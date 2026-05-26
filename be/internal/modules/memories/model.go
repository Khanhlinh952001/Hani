package memories

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Memory struct {
	ID              uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	UserID          int              `json:"user_id" gorm:"not null;index"`
	Content         string           `json:"content" gorm:"not null;type:text"`
	MemoryType      string           `json:"memory_type"`
	ImportanceScore int              `json:"importance_score" gorm:"default:1"`
	Embedding       *pgvector.Vector `json:"embedding,omitempty" gorm:"type:vector(1536)"`
	CreatedAt       time.Time        `json:"created_at"`
}

func (Memory) TableName() string {
	return "memories"
}

func (m *Memory) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.ImportanceScore == 0 {
		m.ImportanceScore = 1
	}
	return nil
}
