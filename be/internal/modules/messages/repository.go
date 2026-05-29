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

func repoGetMessagesPage(sessionID uuid.UUID, limit int, beforeID *uuid.UUID) ([]Message, bool, error) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}

	query := db.DB.Where("session_id = ?", sessionID)
	if beforeID != nil {
		var anchor Message
		if err := db.DB.First(&anchor, "id = ?", *beforeID).Error; err != nil {
			return nil, false, errors.New("message not found")
		}
		query = query.Where(
			"created_at < ? OR (created_at = ? AND id < ?)",
			anchor.CreatedAt, anchor.CreatedAt, *beforeID,
		)
	}

	var msgs []Message
	err := query.Order("created_at desc, id desc").Limit(limit + 1).Find(&msgs).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(msgs) > limit
	if hasMore {
		msgs = msgs[:limit]
	}
	// reverse to chronological
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, hasMore, nil
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
