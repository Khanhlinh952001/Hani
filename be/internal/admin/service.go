package admin

import (
	"errors"
	"strconv"

	"be/internal/db"
	"be/internal/modules/memories"
	"be/internal/modules/messages"
	"be/internal/modules/sessions"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stats struct {
	Users    int64 `json:"users"`
	Sessions int64 `json:"sessions"`
	Messages int64 `json:"messages"`
	Memories int64 `json:"memories"`
}

func GetStatsService() (*Stats, error) {
	var s Stats
	if err := db.DB.Model(&users.User{}).Count(&s.Users).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&sessions.Session{}).Count(&s.Sessions).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&messages.Message{}).Count(&s.Messages).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&memories.Memory{}).Count(&s.Memories).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func ListUsersService() ([]users.User, error) {
	return users.GetAllUsersService()
}

type patchUserInput struct {
	Name             string `json:"name"`
	Status           *int   `json:"status"`
	Role             *int   `json:"role"`
	SubscriptionPlan string `json:"subscription_plan"`
	IsActive         *bool  `json:"is_active"`
}

func PatchUserService(id string, in patchUserInput) (*users.User, error) {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if _, err := users.GetUserByIDService(id); err != nil {
		return nil, err
	}

	patch := &users.User{}
	if in.Name != "" {
		patch.Name = in.Name
	}
	if in.Status != nil {
		patch.Status = *in.Status
	}
	if in.Role != nil {
		patch.Role = *in.Role
	}
	if in.SubscriptionPlan != "" {
		patch.SubscriptionPlan = in.SubscriptionPlan
	}
	if err := users.UpdateUserService(id, patch, ""); err != nil {
		return nil, err
	}
	if in.IsActive != nil {
		if err := users.SetUserActiveService(id, *in.IsActive); err != nil {
			return nil, err
		}
	}
	return users.GetUserByIDService(strconv.Itoa(parsed))
}

func DeleteUserCascadeService(id string) error {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var sessionIDs []uuid.UUID
		if err := tx.Model(&sessions.Session{}).
			Where("user_id = ?", parsed).
			Pluck("id", &sessionIDs).Error; err != nil {
			return err
		}
		if len(sessionIDs) > 0 {
			if err := tx.Where("session_id IN ?", sessionIDs).Delete(&messages.Message{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", sessionIDs).Delete(&sessions.Session{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("user_id = ?", parsed).Delete(&memories.Memory{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&users.User{}, parsed)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found")
		}
		return nil
	})
}

func ListUserSessionsService(userID int) ([]sessions.Session, error) {
	var list []sessions.Session
	err := db.DB.Where("user_id = ?", userID).Order("started_at desc").Find(&list).Error
	return list, err
}

func ListSessionMessagesService(sessionID string) ([]messages.Message, error) {
	parsed, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.New("invalid session id")
	}
	return messages.GetMessagesBySessionIDService(parsed.String())
}

func ListUserMemoriesService(userID int) ([]memories.Memory, error) {
	return memories.GetMemoriesByUserIDService(userID, "")
}

func ClearUserMemoriesService(userID int) error {
	return memories.DeleteMemoriesByUserIDService(userID)
}

func ClearUserConversationService(userID int) error {
	sess, err := sessions.GetOrCreateUserSession(userID)
	if err != nil {
		return err
	}
	if err := sessions.ClearSessionMessagesService(sess.ID); err != nil {
		return err
	}
	return memories.DeleteMemoriesByUserIDService(userID)
}
