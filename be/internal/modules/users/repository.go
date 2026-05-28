package users

import (
	"errors"

	"be/internal/db"

	"gorm.io/gorm"
)

func repoCreateUser(user *User) error {
	return db.DB.Create(user).Error
}

func repoGetAllUsers() ([]User, error) {
	var list []User
	err := db.DB.Find(&list).Error
	return list, err
}

func repoGetUserByID(id int) (*User, error) {
	var user User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func repoGetUserByProviderId(provider, providerID string) (*User, error) {
	var user User
	result := db.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func repoGetUserByEmail(email string) (*User, error) {
	var user User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func repoUpdateUser(id int, data *User) error {
	var user User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}

	updates := map[string]interface{}{
		"name":                  data.Name,
		"email":                 data.Email,
		"password":              data.Password,
		"phone_number":          data.PhoneNumber,
		"provider":              data.Provider,
		"provider_id":           data.ProviderId,
		"avatar":                data.Avatar,
		"gender":                data.Gender,
		"selected_character_id": data.SelectedCharacterID,
		"level":                 data.Level,
		"address":               data.Address,
		"status":                data.Status,
		"role":                  data.Role,
		"subscription_plan":     data.SubscriptionPlan,
		"is_active":             data.IsActive,
	}
	if data.AiProfileID != nil {
		updates["ai_profile_id"] = *data.AiProfileID
	} else {
		updates["ai_profile_id"] = nil
	}
	return db.DB.Model(&user).Updates(updates).Error
}

func repoDeleteUser(id int) error {
	var user User
	result := db.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}
	return db.DB.Delete(&user).Error
}
