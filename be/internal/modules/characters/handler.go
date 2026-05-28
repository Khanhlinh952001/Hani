package characters

import (
	"net/http"
	"strconv"

	"be/internal/auth"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

type selectRequest struct {
	CharacterID string `json:"character_id" binding:"required"`
}

func ListHandler(c *gin.Context) {
	userGender := ""
	if userID, ok := auth.UserID(c); ok {
		if u, err := users.GetUserByIDService(strconv.Itoa(userID)); err == nil {
			userGender = u.Gender
		}
	} else if g := c.Query("gender"); g != "" {
		userGender = g
	}

	list, err := ListForUserService(userGender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func SelectHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req selectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := SelectForUserService(userID, req.CharacterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := users.SetSelectedCharacterService(userID, req.CharacterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func PreviewVoiceHandler(c *gin.Context) {
	slug := c.Param("id")
	b64, format, err := PreviewVoiceKO(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"audio": b64, "format": format})
}
