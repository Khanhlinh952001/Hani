package translate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	Text string `json:"text" binding:"required"`
}

type response struct {
	Translation string `json:"translation"`
}

func Handler(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vi, err := ToVietnamese(c.Request.Context(), req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Translation: vi})
}
