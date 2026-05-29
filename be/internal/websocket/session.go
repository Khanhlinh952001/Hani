package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"be/internal/ai"
	"be/internal/billing"
	"be/internal/conversation"
	"be/internal/tts"

	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

const (
	maxRecentTurns    = 12
	maxHistoryOnReady = 40 // recent messages sent to client on connect (mobile-friendly)
	memorySearchLimit = 5
	turnTimeout       = 2 * time.Minute
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

	characterID        string
	characterName      string
	personalityPrompt  string

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

	isGuest   bool
	guestID   uuid.UUID
	plan      string
	persist   bool
	localTurns []conversation.Turn
}

func (s *RealtimeSession) write(msg ServerMessage) error {
	return s.conn.WriteJSON(msg)
}

func (s *RealtimeSession) Run(ctx context.Context) {
	defer func() {
		s.hub.Unregister(s.connID)
		_ = s.conn.Close()
	}()

	recent, _ := s.recentTurns(maxRecentTurns)
	msgCount := 0
	if !s.isGuest {
		msgCount, _ = conversation.CountUserMessages(s.userID)
	}
	var lastUserAt time.Time
	if !s.isGuest {
		lastUserAt, _ = conversation.LastUserMessageAt(s.userID)
	}
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
		recent, _ := s.recentTurns(maxRecentTurns)
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
				recent, _ := s.recentTurns(maxRecentTurns)
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

func (s *RealtimeSession) sendReadyWithHistory(recent []conversation.Turn) error {
	historyTurns, hasMore, err := s.historyForReady()
	if err != nil {
		return err
	}
	history := make([]HistoryMessage, 0, len(historyTurns))
	for _, t := range historyTurns {
		history = append(history, HistoryMessage{
			ID:          t.ID.String(),
			Role:        t.Role,
			Content:     t.Content,
			Translation: t.TranslationVi,
		})
	}
	return s.write(ServerMessage{
		Type:           EventReady,
		UserID:         s.userID,
		SessionID:      s.sessionID.String(),
		Text:           readyHint(s.voiceEnabled),
		SttContext:     formatSTTContext(recent),
		Messages:       history,
		HistoryHasMore: hasMore,
	})
}

func (s *RealtimeSession) historyForReady() ([]conversation.Turn, bool, error) {
	if !s.persist {
		turns := s.localTurns
		hasMore := len(turns) > maxHistoryOnReady
		if len(turns) > maxHistoryOnReady {
			turns = turns[len(turns)-maxHistoryOnReady:]
		}
		return turns, hasMore, nil
	}
	page, err := conversation.RecentTurnsPage(s.sessionID, maxHistoryOnReady)
	if err != nil {
		return nil, false, err
	}
	return page.Turns, page.HasMore, nil
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

	if err := s.consumeMessageQuota(); err != nil {
		return s.write(ServerMessage{
			Type:    EventQuotaExceeded,
			Message: "Đã hết lượt trò chuyện hôm nay. Nâng cấp Premium để tiếp tục 💕",
			Code:    "quota_exceeded",
		})
	}

	vi := clientVi
	_ = s.write(ServerMessage{
		Type:        EventFinalTranscript,
		Text:        final,
		Translation: vi,
		SessionID:   s.sessionID.String(),
	})

	if s.persist {
		if _, err := conversation.SaveMessage(s.sessionID, "user", final, vi); err != nil {
			log.Printf("[ws] save user message: %v", err)
		}
	} else {
		s.appendLocalTurn("user", final, vi)
	}

	return s.replyToUser(ctx, final)
}

func (s *RealtimeSession) consumeMessageQuota() error {
	if s.isGuest {
		return billing.ConsumeGuestMessage(s.guestID)
	}
	return billing.ConsumeMessage(s.userID, s.plan)
}

func (s *RealtimeSession) recentTurns(limit int) ([]conversation.Turn, error) {
	if !s.persist {
		if len(s.localTurns) <= limit {
			return s.localTurns, nil
		}
		return s.localTurns[len(s.localTurns)-limit:], nil
	}
	return conversation.RecentTurns(s.sessionID, limit)
}

func (s *RealtimeSession) appendLocalTurn(role, content, vi string) {
	s.localTurns = append(s.localTurns, conversation.Turn{
		ID:            uuid.New(),
		Role:          role,
		Content:       content,
		TranslationVi: vi,
	})
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
	characterID, characterName, personalityPrompt string,
	ttsProvider, ttsVoice, ttsLanguage string,
	showVietnamese bool,
	voiceEnabled bool,
	isGuest bool,
	guestID uuid.UUID,
	plan string,
	persist bool,
) *RealtimeSession {
	if characterID == "" {
		characterID = "hani"
	}
	if characterName == "" {
		characterName = "Hani"
	}
	if ttsVoice == "" {
		ttsVoice = "nova"
	}
	return &RealtimeSession{
		connID:            newConnID(),
		userID:            userID,
		userName:          userName,
		userGender:        userGender,
		sessionID:         sessionID,
		characterID:       characterID,
		characterName:     characterName,
		personalityPrompt: personalityPrompt,
		conn:              conn,
		hub:               hub,
		ttsProvider:       tts.NormalizeProvider(ttsProvider),
		ttsVoice:          ttsVoice,
		ttsLanguage:       ttsLanguage,
		showVietnamese:    showVietnamese,
		voiceEnabled:      voiceEnabled,
		isGuest:           isGuest,
		guestID:           guestID,
		plan:              plan,
		persist:           persist,
	}
}
