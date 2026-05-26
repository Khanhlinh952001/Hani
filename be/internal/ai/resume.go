package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type ResumeInput struct {
	UserName          string
	RelationshipStage RelationshipStage
	EmotionState      EmotionalState
	Mood              Mood
	Life              LifeState
	InnerThought      string
	RecentTurns       []Turn
	FactualMemories   []string
	EmotionalMemories []string
	TimeContext       string
	HoursSinceUser    float64
	IncludeVietnamese bool
	ProactiveKind     ProactiveKind
}

// StreamResume greets someone returning to an ongoing chat (references last thread).
func StreamResume(
	ctx context.Context,
	in ResumeInput,
	onDelta func(delta string) error,
) (BilingualReply, error) {
	key := APIKey()
	if key == "" {
		return BilingualReply{}, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	var prompt strings.Builder
	prompt.WriteString("They opened chat again — you already have history.\n")
	prompt.WriteString("ONE short proactive line (1–2 sentences). Reference something specific from recent chat.\n")
	prompt.WriteString("Do NOT say generic 안녕. Do NOT start with 여보.\n")

	turnIn := ReplyInput{
		UserName:          in.UserName,
		RelationshipStage: in.RelationshipStage,
		EmotionState:      in.EmotionState,
		Mood:              in.Mood,
		Life:              in.Life,
		InnerThought:      in.InnerThought,
		FactualMemories:   in.FactualMemories,
		EmotionalMemories: in.EmotionalMemories,
		RecentTurns:       in.RecentTurns,
		UserMessage:       prompt.String(),
		TimeContext:       in.TimeContext,
		HoursSinceUser:    in.HoursSinceUser,
		IncludeVietnamese: in.IncludeVietnamese,
		ProactiveKind:     in.ProactiveKind,
	}
	if turnIn.InnerThought == "" {
		turnIn.InnerThought = GenerateInnerThought(turnIn)
	}

	messages := BuildChatMessages(turnIn)

	stripOpener := HasPriorAssistant(in.RecentTurns)
	if in.IncludeVietnamese {
		return completeBilingualReply(ctx, messages, 0.88, onDelta, stripOpener)
	}

	client := openai.NewClient(key)
	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       ChatModel(),
		Messages:    messages,
		Temperature: 0.88,
		Stream:      true,
	})
	if err != nil {
		return BilingualReply{}, err
	}
	defer stream.Close()

	var acc strings.Builder
	var koSent int
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return BilingualReply{Korean: cleanDisplayText(acc.String())}, err
		}
		if len(resp.Choices) == 0 {
			continue
		}
		delta := resp.Choices[0].Delta.Content
		if delta == "" {
			continue
		}
		if err := streamKoreanOnly(onDelta, delta, &acc, &koSent, stripOpener); err != nil {
			return BilingualReply{Korean: cleanDisplayText(acc.String())}, err
		}
	}

	out := BilingualReply{Korean: cleanDisplayText(strings.TrimSpace(acc.String()))}
	if stripOpener {
		out.Korean = StripYeoboOpener(out.Korean)
	}
	if out.Korean == "" {
		return BilingualReply{}, fmt.Errorf("empty resume response")
	}
	return out, nil
}
