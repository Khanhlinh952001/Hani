package stt

import (
	"net/http"

	"be/internal/auth"
	"be/internal/billing"

	"github.com/gin-gonic/gin"
)

func TemporaryKeyHandler(c *gin.Context) {
	if auth.IsGuest(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "voice_not_allowed", "code": "plan_required"})
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	plan := auth.UserPlan(c)
	if plan == "" {
		plan = billing.PlanFree
	}
	if !billing.AllowsVoice(plan) {
		c.JSON(http.StatusForbidden, gin.H{"error": "voice_not_allowed", "code": "plan_required"})
		return
	}
	if err := billing.CheckVoiceAllowed(userID, plan); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error(), "code": "quota_exceeded"})
		return
	}

	key, err := CreateTemporaryTranscribeKey(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": TemporaryKeyErrorMessage(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"apiKey": key})
}
