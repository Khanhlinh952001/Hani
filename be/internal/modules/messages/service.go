package messages

import (
	"errors"

	"github.com/google/uuid"
)

func CreateMessageService(msg *Message) error {
	if !repoSessionExists(msg.SessionID) {
		return errors.New("session not found")
	}
	if msg.Role != "user" && msg.Role != "assistant" && msg.Role != "system" {
		return errors.New("role must be user, assistant, or system")
	}
	return repoCreateMessage(msg)
}

func GetMessagesBySessionIDService(sessionID string) ([]Message, error) {
	parsed, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.New("invalid session id")
	}
	if !repoSessionExists(parsed) {
		return nil, errors.New("session not found")
	}
	return repoGetMessagesBySessionID(parsed)
}

func GetMessageByIDService(id string) (*Message, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid message id")
	}
	return repoGetMessageByID(parsed)
}

func DeleteMessageService(id string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid message id")
	}
	return repoDeleteMessage(parsed)
}
