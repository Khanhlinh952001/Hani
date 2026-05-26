package ai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// TranslateToVietnamese returns a natural Vietnamese gloss for Korean chat text.
func TranslateToVietnamese(ctx context.Context, korean string) (string, error) {
	korean = strings.TrimSpace(korean)
	if korean == "" {
		return "", nil
	}

	key := APIKey()
	if key == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: ChatModel(),
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You translate Korean to natural Vietnamese for language learners. " +
					"Output only the Vietnamese translation, one or two short lines, no quotes.",
			},
			{Role: openai.ChatMessageRoleUser, Content: korean},
		},
		Temperature: 0.3,
		MaxTokens:   120,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty translation")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
