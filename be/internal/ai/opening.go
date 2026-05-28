package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpeningInput struct {
	UserName          string
	UserGender        string
	RelationshipStage RelationshipStage
	EmotionState      EmotionalState
	Mood              Mood
	Life              LifeState
	InnerThought      string
	FactualMemories   []string
	EmotionalMemories []string
	TimeContext       string
	IncludeVietnamese bool
	PersonalityPrompt string
}

// StreamOpening generates Hani's proactive greeting when user opens chat.
func StreamOpening(
	ctx context.Context,
	in OpeningInput,
	onDelta func(delta string) error,
) (BilingualReply, error) {
	key := APIKey()
	if key == "" {
		return BilingualReply{}, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	var prompt strings.Builder
	prompt.WriteString("Your partner just opened the chat.\n")
	prompt.WriteString("Start first — one short natural line. Not a survey.\n")
	prompt.WriteString("You were doing something before they opened chat — let that color your mood.\n")
	prompt.WriteString("Show your mood from [Hani emotional state] and [Hani's life right now]. Do NOT sound like a chatbot hello.\n")

	turnIn := ReplyInput{
		UserName:          in.UserName,
		UserGender:        in.UserGender,
		RelationshipStage: in.RelationshipStage,
		EmotionState:      in.EmotionState,
		Mood:              in.Mood,
		Life:              in.Life,
		InnerThought:      in.InnerThought,
		FactualMemories:   in.FactualMemories,
		EmotionalMemories: in.EmotionalMemories,
		TimeContext:       in.TimeContext,
		IncludeVietnamese: in.IncludeVietnamese,
		PersonalityPrompt: in.PersonalityPrompt,
	}
	if turnIn.InnerThought == "" {
		turnIn.InnerThought = GenerateInnerThought(turnIn)
	}

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: BuildSystemContent(in.RelationshipStage, in.PersonalityPrompt) + "\n\n" + BuildTurnContext(turnIn)},
		{Role: openai.ChatMessageRoleUser, Content: prompt.String()},
	}
	if in.IncludeVietnamese {
		appendBilingualFormatReminder(&messages)
	}

	if in.IncludeVietnamese {
		return completeBilingualReply(ctx, messages, 0.9, onDelta, false)
	}

	client := openai.NewClient(key)
	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       ChatModel(),
		Messages:    messages,
		Temperature: 0.9,
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
		if err := streamKoreanOnly(onDelta, delta, &acc, &koSent, false); err != nil {
			return BilingualReply{Korean: cleanDisplayText(acc.String())}, err
		}
	}

	out := BilingualReply{Korean: cleanDisplayText(strings.TrimSpace(acc.String()))}
	if out.Korean == "" {
		return BilingualReply{}, fmt.Errorf("empty opening response")
	}
	return out, nil
}
