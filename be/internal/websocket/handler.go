package websocket

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"be/internal/auth"
	"be/internal/billing"
	"be/internal/config"
	"be/internal/modules/characters"
	"be/internal/modules/lover"
	"be/internal/modules/sessions"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleChat upgrades to websocket (requires JWT via ?token= or Authorization header).
func HandleChat(c *gin.Context) {
	claims, err := auth.ParseRequestToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	voiceEnabled := c.Query("practice_mode") != "chat"
	plan := claims.Plan
	if plan == "" {
		plan = billing.PlanFree
	}

	var (
		userID            int
		userName          = claims.Name
		userGender        string
		characterID       = "hani"
		characterName     = "Hani"
		personalityPrompt string
		defaultVoice      string
		defaultLang       = "ko"
		defaultProvider   string
		sessionID       uuid.UUID
		isGuest           = claims.Guest
		guestID           uuid.UUID
		persist           = true
	)

	if isGuest {
		if claims.GuestID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid guest token"})
			return
		}
		guestID, err = uuid.Parse(claims.GuestID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid guest token"})
			return
		}
		plan = billing.PlanGuest
		persist = false
		sessionID = uuid.New()
		if voiceEnabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "voice_not_allowed", "code": "plan_required"})
			return
		}
	} else {
		userID = claims.UserID
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		u, err := users.GetUserByIDService(strconv.Itoa(userID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		if !u.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "account_banned"})
			return
		}
		plan = billing.PlanForUser(u)
		userGender = u.Gender
		if u.SelectedCharacterID != "" {
			characterID = u.SelectedCharacterID
		}

		if voiceEnabled {
			if err := billing.CheckVoiceAllowed(userID, plan); err != nil {
				code := "quota_exceeded"
				if err == billing.ErrVoiceDisabled {
					code = "voice_not_allowed"
				}
				c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error(), "code": code})
				return
			}
		}

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
		sessionID = sess.ID
	}

	rawConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	conn := newConn(rawConn)
	showVi := c.Query("show_vietnamese") != "0"
	ttsProvider := c.Query("tts_provider")
	if ttsProvider == "" {
		ttsProvider = defaultProvider
	}
	if ttsProvider == "" {
		ttsProvider = config.GetEnv("TTS_PROVIDER", "soniox")
	}
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
		userName,
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
		isGuest,
		guestID,
		plan,
		persist,
	)
	DefaultHub.Register(rs)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rs.Run(ctx)
}
