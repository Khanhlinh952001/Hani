package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// StreamWelcomeBack greets the user when they reconnect to an existing session.
func StreamWelcomeBack(ctx context.Context, userName string, onDelta func(delta string) error) (string, error) {
	key := APIKey()
	if key == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	name := strings.TrimSpace(userName)
	if name == "" {
		name = "친구"
	}

	prompt := strings.Join([]string{
		"The user reconnected to the chat.",
		"User name: " + name,
		"",
		"Give a very short warm Korean greeting (1–2 sentences).",
		"Ask one light follow-up question.",
		"Do not repeat a long introduction.",
		"Reply in Korean only as Hani.",
	}, "\n")

	client := openai.NewClient(key)
	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: ChatModel(),
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: SystemPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0.85,
		Stream:      true,
	})
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var full strings.Builder
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return full.String(), err
		}
		if len(resp.Choices) == 0 {
			continue
		}
		delta := resp.Choices[0].Delta.Content
		if delta == "" {
			continue
		}
		full.WriteString(delta)
		if onDelta != nil {
			if err := onDelta(delta); err != nil {
				return full.String(), err
			}
		}
	}

	out := strings.TrimSpace(full.String())
	if out == "" {
		return "", fmt.Errorf("empty welcome response")
	}
	return out, nil
}
