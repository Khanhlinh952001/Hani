package push

import (
	"net/http"
	"strings"

	"be/internal/auth"

	"github.com/gin-gonic/gin"
)

type registerDeviceBody struct {
	FCMToken   string `json:"fcm_token" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"`
}

type heartbeatBody struct {
	FCMToken string `json:"fcm_token" binding:"required"`
}

func RegisterDeviceHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body registerDeviceBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deviceType := strings.ToLower(strings.TrimSpace(body.DeviceType))
	if deviceType != "android" && deviceType != "ios" && deviceType != "web" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_type must be android, ios, or web"})
		return
	}

	d, err := RegisterDevice(userID, RegisterDeviceInput{
		FCMToken:   strings.TrimSpace(body.FCMToken),
		DeviceType: deviceType,
		UserAgent:  c.GetHeader("User-Agent"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

func HeartbeatHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body heartbeatBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := Heartbeat(userID, strings.TrimSpace(body.FCMToken)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func TestPushHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	title := "Hani ❤️"
	body := "Test thông báo — em đây nè!"
	if err := SendTestToUser(c.Request.Context(), userID, title, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "title": title, "body": body})
}

func RevokeDeviceHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}

	if err := RevokeDevice(userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
