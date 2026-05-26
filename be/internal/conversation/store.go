package conversation

import (
	"errors"
	"time"

	"be/internal/db"
	"be/internal/modules/messages"
	"be/internal/modules/sessions"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Turn struct {
	ID            uuid.UUID
	Role          string
	Content       string
	TranslationVi string
}

func UserExists(userID int) bool {
	var count int64
	db.DB.Model(&users.User{}).Where("id = ?", userID).Count(&count)
	return count > 0
}

func CreateSession(userID int) (*sessions.Session, error) {
	s := &sessions.Session{UserID: userID}
	if err := db.DB.Create(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func GetSession(id uuid.UUID) (*sessions.Session, error) {
	var s sessions.Session
	err := db.DB.First(&s, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &s, nil
}

func GetSessionForUser(id uuid.UUID, userID int) (*sessions.Session, error) {
	s, err := GetSession(id)
	if err != nil {
		return nil, err
	}
	if s.UserID != userID {
		return nil, errors.New("session not found")
	}
	return s, nil
}

func EndSession(id uuid.UUID) error {
	now := time.Now()
	result := db.DB.Model(&sessions.Session{}).Where("id = ?", id).Update("ended_at", now)
	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}
	return result.Error
}

func SaveMessage(sessionID uuid.UUID, role, content, translationVi string) (*messages.Message, error) {
	msg := &messages.Message{
		SessionID:     sessionID,
		Role:          role,
		Content:       content,
		TranslationVi: translationVi,
	}
	if err := db.DB.Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

// RecentTurns returns last N messages for prompt context (not full history).
func RecentTurns(sessionID uuid.UUID, limit int) ([]Turn, error) {
	var msgs []messages.Message
	err := db.DB.Where("session_id = ?", sessionID).
		Order("created_at desc").
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	// reverse to chronological
	turns := make([]Turn, 0, len(msgs))
	for i := len(msgs) - 1; i >= 0; i-- {
		role := msgs[i].Role
		if role == "assistant" {
			role = "assistant"
		}
		turns = append(turns, Turn{
			ID:            msgs[i].ID,
			Role:          role,
			Content:       msgs[i].Content,
			TranslationVi: msgs[i].TranslationVi,
		})
	}
	return turns, nil
}

// CountUserMessages counts all messages across a user's sessions (relationship depth).
func CountUserMessages(userID int) (int, error) {
	var count int64
	err := db.DB.Table("messages").
		Joins("JOIN conversation_sessions ON conversation_sessions.id = messages.session_id").
		Where("conversation_sessions.user_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

// LastUserMessageAt returns when the user last sent a message (any session).
func LastUserMessageAt(userID int) (time.Time, error) {
	var createdAt time.Time
	err := db.DB.Table("messages").
		Select("messages.created_at").
		Joins("JOIN conversation_sessions ON conversation_sessions.id = messages.session_id").
		Where("conversation_sessions.user_id = ? AND messages.role = ?", userID, "user").
		Order("messages.created_at desc").
		Limit(1).
		Pluck("messages.created_at", &createdAt).Error
	return createdAt, err
}
