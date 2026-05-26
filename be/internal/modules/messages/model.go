package messages

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	SessionID     uuid.UUID `json:"session_id" gorm:"type:uuid;not null;index"`
	Role          string    `json:"role" gorm:"not null"`
	Content       string    `json:"content" gorm:"not null;type:text"`
	TranslationVi string    `json:"translation_vi,omitempty" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
