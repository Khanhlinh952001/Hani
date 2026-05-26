package websocket

import (
	"context"
	"time"

	"be/internal/ai"
	"be/internal/conversation"
	"be/internal/memory"
)

func toAITurns(recent []conversation.Turn) []ai.Turn {
	out := make([]ai.Turn, 0, len(recent))
	for _, t := range recent {
		out = append(out, ai.Turn{Role: t.Role, Content: t.Content})
	}
	return out
}

func (s *RealtimeSession) buildReplyInput(
	recent []ai.Turn,
	userMsg string,
	proactive ai.ProactiveKind,
	retrieved memory.RetrievedMemories,
) ai.ReplyInput {
	now := time.Now()
	lastUserAt, _ := conversation.LastUserMessageAt(s.userID)
	hoursSince := 999.0
	if !lastUserAt.IsZero() {
		hoursSince = now.Sub(lastUserAt).Hours()
	}
	if userMsg != "" {
		hoursSince = 0
	}

	life := s.life
	if life.CurrentActivity == "" {
		life = ai.BootstrapLife(now, s.emotion, hoursSince, s.userID)
	}
	mood := s.mood
	if mood.Affection == 0 && mood.Loneliness == 0 && mood.Stress == 0 {
		mood = ai.DeriveMood(s.emotion, life, hoursSince)
	}

	in := ai.ReplyInput{
		UserName:          s.userName,
		RelationshipStage: s.relationship,
		EmotionState:      s.emotion,
		Mood:              mood,
		Life:              life,
		FactualMemories:   retrieved.Factual,
		EmotionalMemories: retrieved.Emotional,
		RecentTurns:       recent,
		UserMessage:       userMsg,
		TimeContext:       ai.FormatTimeContext(now),
		HoursSinceUser:    hoursSince,
		IncludeVietnamese: s.showVietnamese,
		ProactiveKind:     proactive,
	}
	in.InnerThought = ai.GenerateInnerThought(in)
	return in
}

// sendOpening greets only on a brand-new session (no prior messages).
func (s *RealtimeSession) sendOpening(ctx context.Context) error {
	recent, _ := conversation.RecentTurns(s.sessionID, 1)
	if len(recent) > 0 {
		return nil
	}

	_ = s.write(ServerMessage{Type: EventTypingStart, SessionID: s.sessionID.String()})

	memCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	retrieved, _ := memory.Retrieve(memCtx, s.userID, "greeting conversation start", memorySearchLimit)
	cancel()

	now := time.Now()
	openingEmotion := ai.BootstrapEmotion(nil, 0, now)
	openingEmotion.MissingUser = true
	openingEmotion.Mood = "soft"
	openingEmotion.CurrentFeeling = "연락 올까 설레"
	life := ai.BootstrapLife(now, openingEmotion, 999, s.userID)
	mood := ai.DeriveMood(openingEmotion, life, 999)
	inner := ai.GenerateInnerThought(ai.ReplyInput{
		RelationshipStage: s.relationship,
		EmotionState:      openingEmotion,
		Mood:              mood,
		Life:              life,
	})

	aiCtx, aiCancel := context.WithTimeout(ctx, 45*time.Second)
	defer aiCancel()
	_, err := s.streamAssistantTurn(aiCtx, func(onDelta func(string) error) (ai.BilingualReply, error) {
		return ai.StreamOpening(aiCtx, ai.OpeningInput{
			UserName:          s.userName,
			RelationshipStage: s.relationship,
			EmotionState:      openingEmotion,
			Mood:              mood,
			Life:              life,
			InnerThought:      inner,
			FactualMemories:   retrieved.Factual,
			EmotionalMemories: retrieved.Emotional,
			TimeContext:       ai.FormatTimeContext(now),
			IncludeVietnamese: s.showVietnamese,
		}, onDelta)
	}, true)
	return err
}

// sendProactiveReachOut — Hani texts first when user returns after silence or late night.
func (s *RealtimeSession) sendProactiveReachOut(
	ctx context.Context,
	recent []conversation.Turn,
	kind ai.ProactiveKind,
	hoursSince float64,
) error {
	_ = s.write(ServerMessage{Type: EventTypingStart, SessionID: s.sessionID.String()})

	query := "continue conversation"
	if len(recent) > 0 {
		query = recent[len(recent)-1].Content
	}
	memCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	retrieved, _ := memory.Retrieve(memCtx, s.userID, query, memorySearchLimit)
	cancel()

	aiTurns := toAITurns(recent)
	turnIn := s.buildReplyInput(aiTurns, "", kind, retrieved)
	turnIn.HoursSinceUser = hoursSince

	aiCtx, aiCancel := context.WithTimeout(ctx, 40*time.Second)
	defer aiCancel()
	_, err := s.streamAssistantTurn(aiCtx, func(onDelta func(string) error) (ai.BilingualReply, error) {
		return ai.StreamResume(aiCtx, ai.ResumeInput{
			UserName:          s.userName,
			RelationshipStage: s.relationship,
			EmotionState:      turnIn.EmotionState,
			Mood:              turnIn.Mood,
			Life:              turnIn.Life,
			InnerThought:      turnIn.InnerThought,
			RecentTurns:       aiTurns,
			FactualMemories:   retrieved.Factual,
			EmotionalMemories: retrieved.Emotional,
			TimeContext:       turnIn.TimeContext,
			HoursSinceUser:    hoursSince,
			IncludeVietnamese: s.showVietnamese,
			ProactiveKind:     kind,
		}, onDelta)
	}, true)
	return err
}

// sendResume — short proactive line when user returns to existing chat.
func (s *RealtimeSession) sendResume(ctx context.Context, recent []conversation.Turn) error {
	_ = s.write(ServerMessage{Type: EventTypingStart, SessionID: s.sessionID.String()})

	query := "continue conversation"
	if len(recent) > 0 {
		query = recent[len(recent)-1].Content
	}
	memCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	retrieved, _ := memory.Retrieve(memCtx, s.userID, query, memorySearchLimit)
	cancel()

	aiTurns := toAITurns(recent)
	turnIn := s.buildReplyInput(aiTurns, "", ai.ProactiveNone, retrieved)

	aiCtx, aiCancel := context.WithTimeout(ctx, 40*time.Second)
	defer aiCancel()
	_, err := s.streamAssistantTurn(aiCtx, func(onDelta func(string) error) (ai.BilingualReply, error) {
		return ai.StreamResume(aiCtx, ai.ResumeInput{
			UserName:          s.userName,
			RelationshipStage: s.relationship,
			EmotionState:      turnIn.EmotionState,
			Mood:              turnIn.Mood,
			Life:              turnIn.Life,
			InnerThought:      turnIn.InnerThought,
			RecentTurns:       aiTurns,
			FactualMemories:   retrieved.Factual,
			EmotionalMemories: retrieved.Emotional,
			TimeContext:       turnIn.TimeContext,
			HoursSinceUser:    turnIn.HoursSinceUser,
			IncludeVietnamese: s.showVietnamese,
		}, onDelta)
	}, true)
	return err
}

func (s *RealtimeSession) replyToUser(ctx context.Context, userText string) error {
	_ = s.write(ServerMessage{Type: EventTypingStart, SessionID: s.sessionID.String()})

	memCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	retrieved, _ := memory.Retrieve(memCtx, s.userID, userText, memorySearchLimit)
	cancel()

	recent, _ := conversation.RecentTurns(s.sessionID, maxRecentTurns)
	aiTurns := toAITurns(recent)
	turnIn := s.buildReplyInput(aiTurns, userText, ai.ProactiveNone, retrieved)

	aiCtx, aiCancel := context.WithTimeout(ctx, 60*time.Second)
	defer aiCancel()

	reply, err := s.streamAssistantTurn(aiCtx, func(onDelta func(string) error) (ai.BilingualReply, error) {
		return ai.StreamReply(aiCtx, turnIn, onDelta)
	}, true)
	if err == nil && reply.Korean != "" {
		s.emotion = ai.EvolveEmotionAfterExchange(s.emotion, userText, reply.Korean)
		s.life = ai.EvolveLifeAfterExchange(s.life, userText)
		s.mood = ai.DeriveMood(s.emotion, s.life, 0)
		memory.SaveFromExchange(s.userID, userText, reply.Korean)
	}
	return err
}
