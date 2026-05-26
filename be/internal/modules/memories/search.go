package memories

import (
	"be/internal/db"

	"github.com/pgvector/pgvector-go"
)

// SearchByVector finds nearest memories for prompt injection (used by realtime pipeline).
func SearchByVector(userID int, embedding pgvector.Vector, limit int) ([]Memory, error) {
	var list []Memory
	err := db.DB.Raw(`
		SELECT * FROM memories
		WHERE user_id = ? AND embedding IS NOT NULL
		ORDER BY embedding <=> ?
		LIMIT ?
	`, userID, embedding, limit).Scan(&list).Error
	return list, err
}
