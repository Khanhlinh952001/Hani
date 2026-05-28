package auth

import (
	"net/http"

	"be/internal/billing"

	"github.com/gin-gonic/gin"
)

func GuestHandler(c *gin.Context) {
	g, err := billing.CreateGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pair, err := IssueTokensForGuest(g.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}
	usage, _ := billing.GetGuestUsageSnapshot(g.ID)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_in":    pair.ExpiresIn,
		"token":         pair.Token,
		"guest_id":      g.ID.String(),
		"usage":         usage,
	})
}
