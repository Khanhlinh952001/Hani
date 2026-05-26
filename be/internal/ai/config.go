package ai

import (
	"os"

	"be/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

func APIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func ChatModel() string {
	return config.GetEnv("OPENAI_MODEL", string(openai.GPT4oMini))
}

func EmbeddingModel() string {
	return config.GetEnv("OPENAI_EMBEDDING_MODEL", "text-embedding-3-small")
}

func EmbeddingDimensions() int {
	return 1536
}
