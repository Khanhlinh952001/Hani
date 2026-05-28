package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"be/internal/conversation"
	"be/internal/ai"
	"be/internal/tts"

	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

const (
	maxRecentTurns     = 12
	maxHistoryOnReady  = 3 // max messages sent to client on connect/reload
	memorySearchLimit  = 5
	turnTimeout        = 2 * time.Minute
)

var errSessionEnded = errors.New("session ended")

var upgrader = gorillaws.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// RealtimeSession handles one websocket connection.
type RealtimeSession struct {
	connID    string
	userID     int
	userName   string
	userGender string
	sessionID  uuid.UUID

	conn      *Conn
	hub       *Hub
	readySent bool

	ttsProvider     string // openai | soniox (from client settings)
	ttsVoice        string
	ttsLanguage     string
	showVietnamese  bool
	voiceEnabled    bool // false = text-only practice (no TTS)

	emotion         ai.EmotionalState
	relationship    ai.RelationshipStage
	life            ai.LifeState
	mood            ai.Mood
}

func (s *RealtimeSession) write(msg ServerMessage) error {
	return s.conn.WriteJSON(msg)
}

func (s *RealtimeSession) Run(ctx context.Context) {
	defer func() {
		s.hub.Unregister(s.connID)
		_ = s.conn.Close()
	}()

	recent, _ := conversation.RecentTurns(s.sessionID, maxRecentTurns)
	msgCount, _ := conversation.CountUserMessages(s.userID)
	lastUserAt, _ := conversation.LastUserMessageAt(s.userID)
	hoursSince := 999.0
	if !lastUserAt.IsZero() {
		hoursSince = time.Since(lastUserAt).Hours()
	}

	now := time.Now()
	s.relationship = ai.RelationshipStageFromMessageCount(msgCount)
	s.emotion = ai.BootstrapEmotion(toAITurns(recent), msgCount, now)
	s.life = ai.BootstrapLife(now, s.emotion, hoursSince, s.userID)
	s.mood = ai.DeriveMood(s.emotion, s.life, hoursSince)

	if err := s.sendReadyWithHistory(recent); err != nil {
		if !errors.Is(err, errConnClosed) {
			log.Printf("[ws %s] ready: %v", s.connID, err)
		}
		return
	}
	s.readySent = true

	go func() {
		if len(recent) == 0 {
			if err := s.sendOpening(ctx); err != nil && !errors.Is(err, errConnClosed) {
				log.Printf("[ws %s] opening: %v", s.connID, err)
				_ = s.write(ServerMessage{Type: EventError, Message: err.Error()})
			}
			return
		}
		kind := ai.DecideProactive(hoursSince, time.Now(), true)
		if kind != ai.ProactiveNone {
			if err := s.sendProactiveReachOut(ctx, recent, kind, hoursSince); err != nil {
				log.Printf("[ws %s] proactive: %v", s.connID, err)
			}
		}
	}()

	for {
		if err := s.runOneTurn(ctx); err != nil {
			if errors.Is(err, errSessionEnded) ||
				errors.Is(err, errConnClosed) ||
				ctx.Err() != nil {
				return
			}
			log.Printf("[ws %s] turn: %v", s.connID, err)
			if werr := s.write(ServerMessage{Type: EventError, Message: err.Error()}); werr != nil && !errors.Is(werr, errConnClosed) {
				log.Printf("[ws %s] error reply: %v", s.connID, werr)
			}
			return
		}
	}
}

func (s *RealtimeSession) runOneTurn(ctx context.Context) error {
	turnCtx, cancel := context.WithTimeout(ctx, turnTimeout)
	defer cancel()

	if !s.readySent {
		recent, _ := conversation.RecentTurns(s.sessionID, maxRecentTurns)
		if err := s.sendReadyWithHistory(recent); err != nil {
			return err
		}
		s.readySent = true
	}

	for {
		msgType, data, err := s.conn.ReadMessage()
		if err != nil {
			return err
		}

		switch msgType {
		case gorillaws.BinaryMessage:
			// STT runs on the client; ignore legacy audio frames.
			continue
		case gorillaws.TextMessage:
			var cm ClientMessage
			if err := json.Unmarshal(data, &cm); err != nil {
				continue
			}
			switch cm.Type {
			case EventPing:
				_ = s.write(ServerMessage{Type: EventPong})
			case EventStartListening:
				recent, _ := conversation.RecentTurns(s.sessionID, maxRecentTurns)
				_ = s.write(ServerMessage{
					Type:       EventListening,
					SessionID:  s.sessionID.String(),
					Text:       "listening",
					SttContext: formatSTTContext(recent),
				})
			case EventSessionEnd:
				_ = conversation.EndSession(s.sessionID)
				_ = s.write(ServerMessage{Type: EventSessionEnded, SessionID: s.sessionID.String()})
				return errSessionEnded
			case EventStopSpeaking:
				return s.processUtterance(turnCtx, strings.TrimSpace(cm.Text), strings.TrimSpace(cm.Translation))
			}
		}
	}
}

func trimHistoryForReady(recent []conversation.Turn) []conversation.Turn {
	if len(recent) <= maxHistoryOnReady {
		return recent
	}
	return recent[len(recent)-maxHistoryOnReady:]
}

func (s *RealtimeSession) sendReadyWithHistory(recent []conversation.Turn) error {
	visible := trimHistoryForReady(recent)
	history := make([]HistoryMessage, 0, len(visible))
	for _, t := range visible {
		history = append(history, HistoryMessage{
			ID:          t.ID.String(),
			Role:        t.Role,
			Content:     t.Content,
			Translation: t.TranslationVi,
		})
	}
	return s.write(ServerMessage{
		Type:       EventReady,
		UserID:     s.userID,
		SessionID:  s.sessionID.String(),
		Text:       readyHint(s.voiceEnabled),
		SttContext: formatSTTContext(recent),
		Messages:   history,
	})
}

func formatSTTContext(turns []conversation.Turn) string {
	if len(turns) == 0 {
		return ""
	}
	var b strings.Builder
	for _, t := range turns {
		if t.Role == "assistant" {
			b.WriteString("Hani: ")
		} else {
			b.WriteString("User: ")
		}
		b.WriteString(t.Content)
		b.WriteString("\n")
	}
	return b.String()
}

func (s *RealtimeSession) processUtterance(ctx context.Context, final, clientVi string) error {
	if final == "" {
		return s.write(ServerMessage{
			Type:    EventError,
			Message: emptyInputHint(s.voiceEnabled),
		})
	}

	vi := clientVi // from Soniox STT when using mic; typed messages may have no vi
	_ = s.write(ServerMessage{
		Type:        EventFinalTranscript,
		Text:        final,
		Translation: vi,
		SessionID:   s.sessionID.String(),
	})

	if _, err := conversation.SaveMessage(s.sessionID, "user", final, vi); err != nil {
		log.Printf("[ws] save user message: %v", err)
	}

	return s.replyToUser(ctx, final)
}

func (s *RealtimeSession) ttsOptions() *tts.Options {
	if s.ttsVoice == "" && s.ttsLanguage == "" {
		return nil
	}
	return &tts.Options{Voice: s.ttsVoice, Language: s.ttsLanguage}
}

func readyHint(voiceEnabled bool) string {
	if voiceEnabled {
		return "mic ready"
	}
	return "text ready"
}

func emptyInputHint(voiceEnabled bool) string {
	if voiceEnabled {
		return "no speech detected — giữ mic lâu hơn (~1 giây) và nói rõ"
	}
	return "nhập tin nhắn trước khi gửi"
}

func NewRealtimeSession(
	userID int,
	userName, userGender string,
	sessionID uuid.UUID,
	conn *Conn,
	hub *Hub,
	ttsProvider, ttsVoice, ttsLanguage string,
	showVietnamese bool,
	voiceEnabled bool,
) *RealtimeSession {
	return &RealtimeSession{
		connID:         newConnID(),
		userID:         userID,
		userName:       userName,
		userGender:     userGender,
		sessionID:      sessionID,
		conn:           conn,
		hub:            hub,
		ttsProvider:    tts.NormalizeProvider(ttsProvider),
		ttsVoice:       ttsVoice,
		ttsLanguage:    ttsLanguage,
		showVietnamese: showVietnamese,
		voiceEnabled:   voiceEnabled,
	}
}
