package memories

import (
	"net/http"

	"be/internal/auth"

	"github.com/gin-gonic/gin"
)

type memoryRequest struct {
	Content         string    `json:"content" binding:"required"`
	MemoryType      string    `json:"memory_type"`
	ImportanceScore int       `json:"importance_score"`
	Embedding       []float32 `json:"embedding"`
}

type searchMemoryRequest struct {
	Embedding []float32 `json:"embedding" binding:"required"`
	Limit     int       `json:"limit"`
}

func CreateMemoryHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req memoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memory := &Memory{
		UserID:          userID,
		Content:         req.Content,
		MemoryType:      req.MemoryType,
		ImportanceScore: req.ImportanceScore,
		Embedding:       ToEmbedding(req.Embedding),
	}

	if err := CreateMemoryService(memory); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, memory)
}

func GetMemoriesHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	list, err := GetMemoriesByUserIDService(userID, c.Query("memory_type"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func SearchMemoriesHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req searchMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := SearchMemoriesService(userID, req.Embedding, req.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func GetMemoryByIDHandler(c *gin.Context) {
	memory, err := GetMemoryByIDService(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memory)
}

func UpdateMemoryHandler(c *gin.Context) {
	var req memoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memory := &Memory{
		Content:         req.Content,
		MemoryType:      req.MemoryType,
		ImportanceScore: req.ImportanceScore,
		Embedding:       ToEmbedding(req.Embedding),
	}

	if err := UpdateMemoryService(c.Param("id"), memory); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	updated, _ := GetMemoryByIDService(c.Param("id"))
	c.JSON(http.StatusOK, updated)
}

func DeleteMemoryHandler(c *gin.Context) {
	if err := DeleteMemoryService(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "memory deleted"})
}
