package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Turn struct {
	Role    string // user | assistant
	Content string
}

type jsonBilingualResponse struct {
	Korean     string `json:"korean"`
	Vietnamese string `json:"vietnamese"`
}

var jsonFenceRe = regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")

// completeBilingualReply — one LLM call; Hani writes Korean + Vietnamese herself.
func completeBilingualReply(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	temperature float32,
	onDelta func(string) error,
	stripOpener bool,
) (BilingualReply, error) {
	key := APIKey()
	if key == "" {
		return BilingualReply{}, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       ChatModel(),
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   400,
	})
	if err != nil {
		return BilingualReply{}, err
	}
	if len(resp.Choices) == 0 {
		return BilingualReply{}, fmt.Errorf("empty ai response")
	}

	choice := resp.Choices[0]
	content := strings.TrimSpace(choice.Message.Content)
	if content == "" {
		if r := strings.TrimSpace(choice.Message.Refusal); r != "" {
			return BilingualReply{}, fmt.Errorf("model refused: %s", r)
		}
		return BilingualReply{}, fmt.Errorf("empty model response")
	}

	out, err := ParseBilingualFlexible(content)
	if err != nil {
		return BilingualReply{}, err
	}
	out = normalizeBilingualScripts(out)
	if stripOpener {
		out.Korean = StripYeoboOpener(out.Korean)
	}
	if out.Korean == "" {
		return BilingualReply{}, fmt.Errorf("empty korean in bilingual reply")
	}
	if !hasHangul(out.Korean) {
		return BilingualReply{}, fmt.Errorf("korean line must contain Hangul")
	}
	if out.Vietnamese == "" {
		return BilingualReply{}, fmt.Errorf("missing vietnamese in model response")
	}

	if onDelta != nil {
		if err := onDelta(out.Korean); err != nil {
			return out, err
		}
	}
	return out, nil
}

// ParseBilingualFlexible accepts ---VI--- text or inline JSON from the model.
func ParseBilingualFlexible(raw string) (BilingualReply, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return BilingualReply{}, fmt.Errorf("empty model response")
	}

	if out := ParseBilingual(raw); out.Korean != "" && out.Vietnamese != "" {
		return out, nil
	}

	if jsonStr := extractJSONObject(raw); jsonStr != "" {
		if out, err := ParseJSONBilingual(jsonStr); err == nil {
			return out, nil
		}
	}

	return BilingualReply{}, fmt.Errorf("could not parse korean and vietnamese from model output")
}

func hasHangul(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7A3 {
			return true
		}
	}
	return false
}

// normalizeBilingualScripts swaps fields if the model reversed Korean/Vietnamese.
func normalizeBilingualScripts(out BilingualReply) BilingualReply {
	if hasHangul(out.Korean) {
		return out
	}
	if hasHangul(out.Vietnamese) {
		return BilingualReply{Korean: out.Vietnamese, Vietnamese: out.Korean}
	}
	return out
}

func extractJSONObject(s string) string {
	if m := jsonFenceRe.FindStringSubmatch(s); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return ""
}

// ParseJSONBilingual parses optional JSON-shaped model output.
func ParseJSONBilingual(raw string) (BilingualReply, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return BilingualReply{}, fmt.Errorf("empty json")
	}
	var j jsonBilingualResponse
	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		return BilingualReply{}, err
	}
	ko := cleanDisplayText(strings.TrimSpace(j.Korean))
	vi := cleanDisplayText(strings.TrimSpace(j.Vietnamese))
	if ko == "" || vi == "" {
		return BilingualReply{}, fmt.Errorf("missing korean or vietnamese field")
	}
	return BilingualReply{Korean: ko, Vietnamese: vi}, nil
}

// StreamReply generates Hani's reply. When IncludeVietnamese, LLM returns both languages in one response.
func StreamReply(ctx context.Context, in ReplyInput, onDelta func(delta string) error) (BilingualReply, error) {
	if strings.TrimSpace(in.UserMessage) == "" {
		return BilingualReply{}, fmt.Errorf("empty user message")
	}

	stripOpener := HasPriorAssistant(in.RecentTurns)
	messages := BuildChatMessages(in)

	if in.IncludeVietnamese {
		return completeBilingualReply(ctx, messages, 0.88, onDelta, stripOpener)
	}

	key := APIKey()
	if key == "" {
		return BilingualReply{}, fmt.Errorf("OPENAI_API_KEY is not set")
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
		return BilingualReply{}, fmt.Errorf("empty ai response")
	}
	return out, nil
}
