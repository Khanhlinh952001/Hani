package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"be/internal/ai"
	memmod "be/internal/modules/memories"

	"github.com/pgvector/pgvector-go"
)

const defaultLimit = 8

// Retrieve relevant memories for the user utterance (pgvector), split by type.
func Retrieve(ctx context.Context, userID int, userText string, limit int) (RetrievedMemories, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	embedding, err := ai.EmbedText(ctx, userText)
	if err != nil {
		return RetrievedMemories{}, err
	}

	list, err := memmod.SearchByVector(userID, pgvector.NewVector(embedding), limit)
	if err != nil {
		return RetrievedMemories{}, err
	}

	var out RetrievedMemories
	now := time.Now()
	for _, m := range list {
		line := formatMemoryLine(m, now)
		if line == "" {
			continue
		}
		switch m.MemoryType {
		case "emotional":
			out.Emotional = append(out.Emotional, line)
		default:
			out.Factual = append(out.Factual, line)
		}
	}
	return out, nil
}

func formatMemoryLine(m memmod.Memory, now time.Time) string {
	line := strings.TrimSpace(m.Content)
	if line == "" {
		return ""
	}

	age := now.Sub(m.CreatedAt)
	switch m.MemoryType {
	case "emotional":
		if age > 90*24*time.Hour {
			return fmt.Sprintf("(vague feeling — long ago) %s — you remember how it felt more than exact words", line)
		}
		if age > 30*24*time.Hour {
			return fmt.Sprintf("(fuzzy) %s — not every detail, but the feeling stayed", line)
		}
		return fmt.Sprintf("[emotional] %s", line)
	default:
		if age > 60*24*time.Hour {
			return fmt.Sprintf("(might be fuzzy) %s", line)
		}
		if age > 14*24*time.Hour {
			return fmt.Sprintf("(older detail) %s", line)
		}
		return line
	}
}

// RetrieveLines is a convenience wrapper returning a flat list (legacy callers).
func RetrieveLines(ctx context.Context, userID int, userText string, limit int) ([]string, error) {
	got, err := Retrieve(ctx, userID, userText, limit)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(got.Factual)+len(got.Emotional))
	out = append(out, got.Factual...)
	out = append(out, got.Emotional...)
	return out, nil
}
