package websocket

import (
	"context"
	"log"

	"be/internal/ai"
	"be/internal/conversation"
	"be/internal/tts"
)

// streamAssistantTurn streams TTS; text (KO+VI) is sent together once the LLM finishes.
func (s *RealtimeSession) streamAssistantTurn(
	ctx context.Context,
	generate func(onDelta func(string) error) (ai.BilingualReply, error),
	persist bool,
) (ai.BilingualReply, error) {
	if !s.voiceEnabled {
		return s.streamAssistantTurnTextOnly(ctx, generate, persist)
	}

	splitter := &ai.SentenceBuffer{}

	sentences := make(chan string, 24)
	ttsErrCh := make(chan error, 1)
	audioStarted := false

	go func() {
		var ttsErr error
		for sent := range sentences {
			if !tts.WorthSpeaking(sent) {
				continue
			}
			err := tts.StreamSpeechFor(ctx, s.ttsProvider, sent, s.ttsOptions(), func(_ int, b64 string) error {
				if !audioStarted {
					audioStarted = true
					_ = s.write(ServerMessage{
						Type:      EventTypingEnd,
						SessionID: s.sessionID.String(),
					})
				}
				return s.write(ServerMessage{
					Type:      EventAIAudioChunk,
					Audio:     b64,
					Format:    tts.AudioFormatFor(s.ttsProvider),
					SessionID: s.sessionID.String(),
				})
			})
			if err != nil {
				ttsErr = err
				break
			}
			_ = s.write(ServerMessage{
				Type:      EventAIAudioSegmentEnd,
				SessionID: s.sessionID.String(),
			})
		}
		ttsErrCh <- ttsErr
	}()

	// Voice mode: no character streaming to UI — subtitle carries full KO+VI together.
	onDelta := func(delta string) error {
		for _, sent := range splitter.Feed(delta) {
			sentences <- sent
		}
		return nil
	}

	reply, err := generate(onDelta)
	if err != nil {
		close(sentences)
		<-ttsErrCh
		_ = s.write(ServerMessage{Type: EventTypingEnd})
		return ai.BilingualReply{}, err
	}

	// Show Korean + Vietnamese together as soon as text is ready (before TTS finishes).
	reply = s.publishAssistantText(ctx, reply, persist)

	if tail := splitter.Flush(); tail != "" && tts.WorthSpeaking(tail) {
		sentences <- tail
	}
	close(sentences)

	if ttsErr := <-ttsErrCh; ttsErr != nil {
		log.Printf("[ws] tts: %v", ttsErr)
		_ = s.write(ServerMessage{
			Type:      EventError,
			Message:   "voice failed: " + ttsErr.Error(),
			SessionID: s.sessionID.String(),
		})
	}

	if !audioStarted {
		_ = s.write(ServerMessage{Type: EventTypingEnd, SessionID: s.sessionID.String()})
	}

	_ = s.write(ServerMessage{
		Type:      EventAIAudioEnd,
		SessionID: s.sessionID.String(),
		Finished:  true,
	})
	return reply, nil
}

func (s *RealtimeSession) streamAssistantTurnTextOnly(
	ctx context.Context,
	generate func(onDelta func(string) error) (ai.BilingualReply, error),
	persist bool,
) (ai.BilingualReply, error) {
	onDelta := func(string) error { return nil }

	reply, err := generate(onDelta)
	if err != nil {
		_ = s.write(ServerMessage{Type: EventTypingEnd, SessionID: s.sessionID.String()})
		return ai.BilingualReply{}, err
	}

	_ = s.write(ServerMessage{Type: EventTypingEnd, SessionID: s.sessionID.String()})
	return s.finishAssistantTurn(ctx, reply, persist)
}

func (s *RealtimeSession) publishAssistantText(_ context.Context, reply ai.BilingualReply, persist bool) ai.BilingualReply {
	if s.showVietnamese && reply.Vietnamese == "" {
		log.Printf("[ws] missing vietnamese from LLM for: %q", reply.Korean)
	}

	_ = s.write(ServerMessage{
		Type:        EventSubtitle,
		Text:        reply.Korean,
		Translation: reply.Vietnamese,
		FullText:    reply.Korean,
		SessionID:   s.sessionID.String(),
	})

	if persist {
		if _, err := conversation.SaveMessage(s.sessionID, "assistant", reply.Korean, reply.Vietnamese); err != nil {
			log.Printf("[ws] save assistant message: %v", err)
		}
	}
	return reply
}

func (s *RealtimeSession) finishAssistantTurn(ctx context.Context, reply ai.BilingualReply, persist bool) (ai.BilingualReply, error) {
	reply = s.publishAssistantText(ctx, reply, persist)
	_ = s.write(ServerMessage{
		Type:      EventAIAudioEnd,
		SessionID: s.sessionID.String(),
		Finished:  true,
	})
	return reply, nil
}
