package ai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func EmbedText(ctx context.Context, text string) ([]float32, error) {
	key := APIKey()
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty embed text")
	}

	client := openai.NewClient(key)
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model:      openai.EmbeddingModel(EmbeddingModel()),
		Input:      []string{text},
		Dimensions: EmbeddingDimensions(),
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}
	return resp.Data[0].Embedding, nil
}
