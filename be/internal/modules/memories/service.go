package memories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

func CreateMemoryService(memory *Memory) error {
	if !repoUserExists(memory.UserID) {
		return errors.New("user not found")
	}
	return repoCreateMemory(memory)
}

func GetMemoriesByUserIDService(userID int, memoryType string) ([]Memory, error) {
	if !repoUserExists(userID) {
		return nil, errors.New("user not found")
	}
	return repoGetMemoriesByUserID(userID, memoryType)
}

func GetMemoryByIDService(id string) (*Memory, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid memory id")
	}
	return repoGetMemoryByID(parsed)
}

func SearchMemoriesService(userID int, embedding []float32, limit int) ([]Memory, error) {
	if !repoUserExists(userID) {
		return nil, errors.New("user not found")
	}
	if len(embedding) == 0 {
		return nil, errors.New("embedding is required")
	}
	if limit <= 0 {
		limit = 5
	}
	return repoSearchMemories(userID, pgvector.NewVector(embedding), limit)
}

func UpdateMemoryService(id string, data *Memory) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid memory id")
	}
	return repoUpdateMemory(parsed, data)
}

func DeleteMemoriesByUserIDService(userID int) error {
	if !repoUserExists(userID) {
		return errors.New("user not found")
	}
	return repoDeleteMemoriesByUserID(userID)
}

func DeleteMemoryService(id string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid memory id")
	}
	return repoDeleteMemory(parsed)
}

func ToEmbedding(values []float32) *pgvector.Vector {
	if len(values) == 0 {
		return nil
	}
	v := pgvector.NewVector(values)
	return &v
}
