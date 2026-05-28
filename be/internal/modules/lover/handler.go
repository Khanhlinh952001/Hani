package lover

import (
	"net/http"
	"strconv"

	"be/internal/auth"
	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

type createProfileRequest struct {
	CompanionGender       string   `json:"companion_gender" binding:"required"`
	PersonalityTemplateID string   `json:"personality_template_id" binding:"required"`
	SpeakingStyleTags     []string `json:"speaking_style_tags"`
	VoiceProfileID        string   `json:"voice_profile_id" binding:"required"`
	DisplayName           string   `json:"display_name" binding:"required"`
	AvatarURL             string   `json:"avatar_url"`
}

type quickPresetRequest struct {
	PresetSlug string `json:"preset_slug" binding:"required"`
}

type previewVoiceRequest struct {
	VoiceProfileID string `json:"voice_profile_id" binding:"required"`
	Text           string `json:"text"`
}

func ListPersonalitiesHandler(c *gin.Context) {
	list, err := ListPersonalitiesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func ListSpeakingStylesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, SpeakingStylesService())
}

func ListVoicesHandler(c *gin.Context) {
	gender := c.Query("companion_gender")
	list, err := ListVoicesService(gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func NameSuggestionsHandler(c *gin.Context) {
	gender := c.Query("companion_gender")
	c.JSON(http.StatusOK, gin.H{"names": NameSuggestions(gender)})
}

func GetMyProfileHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	profile, err := GetProfileForUserService(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no profile"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func CreateProfileHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := CreateProfileService(userID, CreateProfileInput{
		CompanionGender:       req.CompanionGender,
		PersonalityTemplateID: req.PersonalityTemplateID,
		SpeakingStyleTags:     req.SpeakingStyleTags,
		VoiceProfileID:        req.VoiceProfileID,
		DisplayName:           req.DisplayName,
		AvatarURL:             req.AvatarURL,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := users.SetAiProfileService(userID, profile.ID, req.CompanionGender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile, "user": user})
}

func CreateQuickPresetHandler(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req quickPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := CreateFromPresetService(userID, req.PresetSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := users.SetAiProfileService(userID, profile.ID, profile.CompanionGender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = users.SetSelectedCharacterService(userID, req.PresetSlug)

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile, "user": user})
}

func PreviewVoiceHandler(c *gin.Context) {
	var req previewVoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b64, format, err := PreviewVoiceService(c.Request.Context(), req.VoiceProfileID, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"audio": b64, "format": format})
}
