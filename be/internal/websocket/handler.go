package websocket

import (
	"context"
	"net/http"
	"time"

	"be/internal/auth"
	"be/internal/config"
	"be/internal/modules/sessions"

	"github.com/gin-gonic/gin"
)

// HandleChat upgrades to websocket (requires JWT via ?token= or Authorization header).
func HandleChat(c *gin.Context) {
	claims, err := auth.ParseRequestToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := claims.UserID

	sess, err := sessions.GetOrCreateUserSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sessionID := sess.ID

	rawConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	conn := newConn(rawConn)
	showVi := c.Query("show_vietnamese") != "0"
	voiceEnabled := c.Query("practice_mode") != "chat"
	ttsProvider := c.Query("tts_provider")
	if ttsProvider == "" {
		ttsProvider = config.GetEnv("TTS_PROVIDER", "openai")
	}
	rs := NewRealtimeSession(
		userID,
		claims.Name,
		sessionID,
		conn,
		DefaultHub,
		ttsProvider,
		c.Query("tts_voice"),
		c.Query("tts_language"),
		showVi,
		voiceEnabled,
	)
	DefaultHub.Register(rs)

	// "ready" with session_id is sent from runOneTurn (STT runs in the browser).

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rs.Run(ctx)
}
