package memories

import (
	"errors"

	"be/internal/db"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

func repoCreateMemory(memory *Memory) error {
	return db.DB.Create(memory).Error
}

func repoGetMemoriesByUserID(userID int, memoryType string) ([]Memory, error) {
	var list []Memory
	query := db.DB.Where("user_id = ?", userID)
	if memoryType != "" {
		query = query.Where("memory_type = ?", memoryType)
	}
	err := query.Order("importance_score desc, created_at desc").Find(&list).Error
	return list, err
}

func repoGetMemoryByID(id uuid.UUID) (*Memory, error) {
	var memory Memory
	result := db.DB.First(&memory, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("memory not found")
		}
		return nil, result.Error
	}
	return &memory, nil
}

func repoSearchMemories(userID int, embedding pgvector.Vector, limit int) ([]Memory, error) {
	var list []Memory
	err := db.DB.Raw(`
		SELECT * FROM memories
		WHERE user_id = ? AND embedding IS NOT NULL
		ORDER BY embedding <=> ?
		LIMIT ?
	`, userID, embedding, limit).Scan(&list).Error
	return list, err
}

func repoUpdateMemory(id uuid.UUID, data *Memory) error {
	var memory Memory
	result := db.DB.First(&memory, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("memory not found")
		}
		return result.Error
	}

	updates := map[string]interface{}{
		"content":          data.Content,
		"memory_type":      data.MemoryType,
		"importance_score": data.ImportanceScore,
	}
	if data.Embedding != nil {
		updates["embedding"] = data.Embedding
	}

	return db.DB.Model(&memory).Updates(updates).Error
}

func repoDeleteMemoriesByUserID(userID int) error {
	return db.DB.Where("user_id = ?", userID).Delete(&Memory{}).Error
}

func repoDeleteMemory(id uuid.UUID) error {
	var memory Memory
	result := db.DB.First(&memory, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("memory not found")
		}
		return result.Error
	}
	return db.DB.Delete(&memory).Error
}

func repoUserExists(userID int) bool {
	var count int64
	db.DB.Model(&users.User{}).Where("id = ?", userID).Count(&count)
	return count > 0
}
