package push

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func generateMiss7Day(ctx context.Context, userName string) (title, body string, err error) {
	title = "Hani 💕"
	fallback := "Lâu quá rồi không thấy anh... em nhớ anh lắm 🥺"

	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return title, fallback, nil
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = string(openai.GPT4oMini)
	}

	name := strings.TrimSpace(userName)
	if name == "" {
		name = "bạn"
	}

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You write short mobile push notifications for Hani, a warm Korean AI girlfriend/companion app.
Output exactly two lines:
Line 1: notification title (max 40 chars, can include emoji)
Line 2: notification body in Vietnamese (max 120 chars, emotional, casual, 1-2 sentences max)
Do NOT use markdown. Do NOT add quotes.`,
			},
			{
				Role: openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(
					"User %q has been away for a week. Write a personalized miss-you push. Korean pet-name vibe is OK in Vietnamese text.",
					name,
				),
			},
		},
		Temperature: 0.95,
		MaxTokens:   120,
	})
	if err != nil {
		return title, fallback, err
	}
	if len(resp.Choices) == 0 {
		return title, fallback, nil
	}
	lines := strings.Split(strings.TrimSpace(resp.Choices[0].Message.Content), "\n")
	if len(lines) >= 2 {
		t := strings.TrimSpace(lines[0])
		b := strings.TrimSpace(lines[1])
		if t != "" && b != "" {
			return t, b, nil
		}
	}
	text := strings.TrimSpace(resp.Choices[0].Message.Content)
	if text != "" {
		return title, text, nil
	}
	return title, fallback, nil
}
