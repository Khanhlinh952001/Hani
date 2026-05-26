package memory

import (
	"context"
	"log"
	"time"

	"be/internal/ai"
	memmod "be/internal/modules/memories"

	"github.com/pgvector/pgvector-go"
)

// RetrievedMemories splits vector hits into factual vs emotional lines for prompts.
type RetrievedMemories struct {
	Factual   []string
	Emotional []string
}

// SaveFromExchange stores factual + emotional memories in the background (best-effort).
func SaveFromExchange(userID int, userMsg, assistantMsg string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		facts, err := ai.ExtractMemoryFacts(ctx, userMsg, assistantMsg)
		if err != nil || len(facts) == 0 {
			return
		}

		for _, fact := range facts {
			if fact.Content == "" {
				continue
			}
			embedding, err := ai.EmbedText(ctx, fact.Content)
			if err != nil {
				log.Printf("[memory] embed: %v", err)
				continue
			}
			importance := 2
			if fact.Type == "emotional" {
				importance = 3
			}
			m := &memmod.Memory{
				UserID:          userID,
				Content:         fact.Content,
				MemoryType:      fact.Type,
				ImportanceScore: importance,
				Embedding:       ptr(pgvector.NewVector(embedding)),
			}
			if err := memmod.CreateMemoryService(m); err != nil {
				log.Printf("[memory] save: %v", err)
			}
		}
	}()
}

func ptr(v pgvector.Vector) *pgvector.Vector {
	return &v
}
