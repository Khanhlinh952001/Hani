package auth

import (
	"net/http"

	"be/internal/billing"

	"github.com/gin-gonic/gin"
)

func UsageHandler(c *gin.Context) {
	if IsGuest(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "registration_required"})
		return
	}
	userID, ok := UserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	snap, err := billing.GetUsageSnapshot(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snap)
}
