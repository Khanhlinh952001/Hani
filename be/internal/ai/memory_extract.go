package ai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// MemoryFact is one storable line for long-term recall (Korean + optional Vietnamese gloss).
type MemoryFact struct {
	Content       string
	TranslationVi string
	Type          string // life | emotional
}

// ExtractMemoryFacts returns factual and emotional memories worth saving, or nil.
func ExtractMemoryFacts(ctx context.Context, userMsg, assistantMsg string) ([]MemoryFact, error) {
	userMsg = strings.TrimSpace(userMsg)
	if userMsg == "" {
		return nil, nil
	}

	key := APIKey()
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	prompt := fmt.Sprintf(`From this chat exchange, extract up to TWO short memories worth keeping for a Korean AI companion app used by Vietnamese learners.

Types:
- life: job, hobby, plan, preference, daily detail about the user
- emotional: how they made each other feel, comfort, missing each other, staying up late talking, warmth

If nothing worth saving, output exactly: NONE

User: %s
Hani: %s

Output format (one per line, max 2 lines). Each line has THREE parts separated by |:
type|Korean sentence|Vietnamese gloss

Example:
life|주말에 밀크티 마시는 걸 좋아함|Thích uống trà sữa cuối tuần
emotional|밤늦게까지 같이 이야기해서 따뜻했음|Đêm qua nói chuyện khuya cùng nhau, ấm áp

Rules:
- Part 2 MUST be natural Korean (한국어), short (under 40 chars).
- Part 3 MUST be natural Vietnamese for the learner.
- Do NOT use English. No quotes.`, userMsg, strings.TrimSpace(assistantMsg))

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       ChatModel(),
		Temperature: 0.2,
		MaxTokens:   120,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, nil
	}

	var out []MemoryFact
	for _, line := range strings.Split(resp.Choices[0].Message.Content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, "NONE") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}
		typ := strings.TrimSpace(parts[0])
		content := strings.TrimSpace(parts[1])
		vi := ""
		if len(parts) >= 3 {
			vi = strings.TrimSpace(parts[2])
		}
		if content == "" {
			continue
		}
		if typ != "emotional" {
			typ = "life"
		}
		out = append(out, MemoryFact{Content: content, TranslationVi: vi, Type: typ})
	}
	return out, nil
}

// ExtractMemoryFact keeps backward compatibility — returns first life fact.
func ExtractMemoryFact(ctx context.Context, userMsg, assistantMsg string) (string, error) {
	facts, err := ExtractMemoryFacts(ctx, userMsg, assistantMsg)
	if err != nil || len(facts) == 0 {
		return "", err
	}
	for _, f := range facts {
		if f.Type == "life" {
			return f.Content, nil
		}
	}
	return facts[0].Content, nil
}
