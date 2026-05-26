package sessions

import (
	"errors"

	"be/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSessionService(session *Session) error {
	if !repoUserExists(session.UserID) {
		return errors.New("user not found")
	}
	return repoCreateSession(session)
}

func GetSessionsByUserIDService(userID int) ([]Session, error) {
	if !repoUserExists(userID) {
		return nil, errors.New("user not found")
	}
	return repoGetSessionsByUserID(userID)
}

func GetSessionByIDService(id string) (*Session, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid session id")
	}
	return repoGetSessionByID(parsed)
}

func GetSessionForUserService(id string, userID int) (*Session, error) {
	session, err := GetSessionByIDService(id)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func EndSessionService(id string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid session id")
	}
	return repoEndSession(parsed)
}

func DeleteSessionService(id string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid session id")
	}
	return repoDeleteSession(parsed)
}

// GetOrCreateUserSession returns the user's single ongoing conversation
// (most recently active), or creates one if none exists.
func GetOrCreateUserSession(userID int) (*Session, error) {
	if !repoUserExists(userID) {
		return nil, errors.New("user not found")
	}

	var sess Session
	err := db.DB.Table("conversation_sessions").
		Select("conversation_sessions.*").
		Joins("LEFT JOIN messages ON messages.session_id = conversation_sessions.id").
		Where("conversation_sessions.user_id = ?", userID).
		Group("conversation_sessions.id").
		Order("COALESCE(MAX(messages.created_at), conversation_sessions.started_at) DESC").
		Limit(1).
		First(&sess).Error
	if err == nil {
		_ = ReactivateSessionService(sess.ID)
		return &sess, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		created := &Session{UserID: userID}
		if err := repoCreateSession(created); err != nil {
			return nil, err
		}
		return created, nil
	}
	return nil, err
}

func ReactivateSessionService(id uuid.UUID) error {
	return db.DB.Model(&Session{}).Where("id = ?", id).Update("ended_at", nil).Error
}

func ClearSessionMessagesService(sessionID uuid.UUID) error {
	return db.DB.Exec("DELETE FROM messages WHERE session_id = ?", sessionID).Error
}
