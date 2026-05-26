package messages

import (
	"errors"
	"be/internal/db"
	"be/internal/modules/sessions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func repoCreateMessage(msg *Message) error {
	return db.DB.Create(msg).Error
}

func repoGetMessagesBySessionID(sessionID uuid.UUID) ([]Message, error) {
	var list []Message
	err := db.DB.Where("session_id = ?", sessionID).Order("created_at asc").Find(&list).Error
	return list, err
}

func repoGetMessageByID(id uuid.UUID) (*Message, error) {
	var msg Message
	result := db.DB.First(&msg, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, result.Error
	}
	return &msg, nil
}

func repoDeleteMessage(id uuid.UUID) error {
	var msg Message
	result := db.DB.First(&msg, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("message not found")
		}
		return result.Error
	}
	return db.DB.Delete(&msg).Error
}

func repoSessionExists(sessionID uuid.UUID) bool {
	var count int64
	db.DB.Model(&sessions.Session{}).Where("id = ?", sessionID).Count(&count)
	return count > 0
}
