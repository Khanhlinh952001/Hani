package characters

import (
	"errors"
	"time"

	"be/internal/db"

	"gorm.io/gorm"
)

func repoListAll() ([]Character, error) {
	var list []Character
	err := db.DB.Order("sort_order asc, id asc").Find(&list).Error
	return list, err
}

func repoGetByID(id string) (*Character, error) {
	var c Character
	err := db.DB.First(&c, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("character not found")
		}
		return nil, err
	}
	return &c, nil
}

func repoUpsertUserMemory(userID int, characterID string) error {
	var row UserCharacterMemory
	err := db.DB.Where("user_id = ? AND character_id = ?", userID, characterID).First(&row).Error
	now := time.Now()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = UserCharacterMemory{
			UserID:             userID,
			CharacterID:        characterID,
			IntimacyLevel:      1,
			RelationshipStatus: "new",
			LastInteractionAt:  now,
		}
		return db.DB.Create(&row).Error
	}
	if err != nil {
		return err
	}
	row.LastInteractionAt = now
	return db.DB.Save(&row).Error
}

func repoTouchInteraction(userID int, characterID string) {
	_ = db.DB.Model(&UserCharacterMemory{}).
		Where("user_id = ? AND character_id = ?", userID, characterID).
		Update("last_interaction_at", time.Now()).Error
}
