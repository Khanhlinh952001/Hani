package websocket

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"be/internal/auth"
	"be/internal/config"
	"be/internal/modules/characters"
	"be/internal/modules/lover"
	"be/internal/modules/sessions"
	"be/internal/modules/users"

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

	userGender := ""
	characterID := "hani"
	characterName := "Hani"
	personalityPrompt := ""
	defaultVoice := ""
	defaultLang := "ko"
	if u, err := users.GetUserByIDService(strconv.Itoa(userID)); err == nil {
		userGender = u.Gender
		if u.SelectedCharacterID != "" {
			characterID = u.SelectedCharacterID
		}
	}
	defaultProvider := ""
	if name, prompt, voice, ok := lover.GetProfilePromptForUser(userID); ok {
		characterName = name
		personalityPrompt = prompt
		defaultVoice = voice.VoiceID
		if voice.Language != "" {
			defaultLang = voice.Language
		}
		defaultProvider = voice.Provider
		characterID = "custom"
	} else if ch, err := characters.GetByIDService(characterID); err == nil {
		characterName = ch.Name
		personalityPrompt = ch.PersonalityPrompt
		defaultVoice = ch.VoiceID
		if ch.TTSLanguage != "" {
			defaultLang = ch.TTSLanguage
		}
		defaultProvider = ch.VoiceProvider
	}

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
		ttsProvider = defaultProvider
	}
	if ttsProvider == "" {
		ttsProvider = config.GetEnv("TTS_PROVIDER", "soniox")
	}
	// App uses Soniox TTS only (OpenAI is for chat/embeddings).
	ttsProvider = "soniox"
	ttsVoice := c.Query("tts_voice")
	if ttsVoice == "" {
		ttsVoice = defaultVoice
	}
	ttsLanguage := c.Query("tts_language")
	if ttsLanguage == "" {
		ttsLanguage = defaultLang
	}
	rs := NewRealtimeSession(
		userID,
		claims.Name,
		userGender,
		sessionID,
		conn,
		DefaultHub,
		characterID,
		characterName,
		personalityPrompt,
		ttsProvider,
		ttsVoice,
		ttsLanguage,
		showVi,
		voiceEnabled,
	)
	DefaultHub.Register(rs)

	// "ready" with session_id is sent from runOneTurn (STT runs in the browser).

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rs.Run(ctx)
}
