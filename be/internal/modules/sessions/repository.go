package sessions

import (
	"errors"
	"time"

	"be/internal/db"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func repoCreateSession(session *Session) error {
	return db.DB.Create(session).Error
}

func repoGetSessionsByUserID(userID int) ([]Session, error) {
	var list []Session
	err := db.DB.Where("user_id = ?", userID).Order("started_at desc").Find(&list).Error
	return list, err
}

func repoGetSessionByID(id uuid.UUID) (*Session, error) {
	var session Session
	result := db.DB.First(&session, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, result.Error
	}
	return &session, nil
}

func repoEndSession(id uuid.UUID) error {
	now := time.Now()
	result := db.DB.Model(&Session{}).Where("id = ?", id).Update("ended_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}

func repoDeleteSession(id uuid.UUID) error {
	var session Session
	result := db.DB.First(&session, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("session not found")
		}
		return result.Error
	}
	if err := db.DB.Exec("DELETE FROM messages WHERE session_id = ?", id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&session).Error
}

func repoUserExists(userID int) bool {
	var count int64
	db.DB.Model(&users.User{}).Where("id = ?", userID).Count(&count)
	return count > 0
}
