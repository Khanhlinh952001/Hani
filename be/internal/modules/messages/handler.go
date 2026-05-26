package messages

import (
	"net/http"

	"be/internal/auth"
	"be/internal/modules/sessions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type messageRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

func CreateMessageHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req messageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}

	if _, err := sessions.GetSessionForUserService(req.SessionID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	msg := &Message{
		SessionID: sessionID,
		Role:      req.Role,
		Content:   req.Content,
	}

	if err := CreateMessageService(msg); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func GetMessagesHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	if _, err := uuid.Parse(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	if _, err := sessions.GetSessionForUserService(sessionID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	list, err := GetMessagesBySessionIDService(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func GetMessageByIDHandler(c *gin.Context) {
	msg, err := GetMessageByIDService(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func DeleteMessageHandler(c *gin.Context) {
	if err := DeleteMessageService(c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}
