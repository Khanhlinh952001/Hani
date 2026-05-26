package stt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TemporaryKeyHandler(c *gin.Context) {
	key, err := CreateTemporaryTranscribeKey(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": TemporaryKeyErrorMessage(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"apiKey": key})
}
