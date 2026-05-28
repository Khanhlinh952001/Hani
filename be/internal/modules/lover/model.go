package lover

import (
	"time"

	"github.com/google/uuid"
)

// PersonalityTemplate is a catalog personality archetype.
type PersonalityTemplate struct {
	ID                   string `json:"id" gorm:"primaryKey"`
	NameKO               string `json:"name_ko"`
	NameVI               string `json:"name_vi"`
	DescriptionKO        string `json:"description_ko"`
	DescriptionVI        string `json:"description_vi"`
	Icon                 string `json:"icon"`
	BasePrompt           string `json:"-" gorm:"type:text"`
	EmojiDensity         int    `json:"emoji_density" gorm:"default:1"`
	TypingSpeed          string `json:"typing_speed" gorm:"default:normal"`
	FlirtingLevel        int    `json:"flirting_level" gorm:"default:1"`
	DefaultSpeakingStyles []string `json:"default_speaking_styles" gorm:"serializer:json"`
	SortOrder            int    `json:"sort_order"`
}

func (PersonalityTemplate) TableName() string { return "ai_personality_templates" }

// VoiceProfile maps TTS settings.
type VoiceProfile struct {
	ID            string `json:"id" gorm:"primaryKey"`
	NameKO        string `json:"name_ko"`
	NameVI        string `json:"name_vi"`
	Gender        string `json:"gender"`
	Provider      string `json:"provider"`
	VoiceID       string `json:"voice_id"`
	Emotion       string `json:"emotion"`
	Speed         string `json:"speed"`
	Language      string `json:"language" gorm:"default:ko"`
	PreviewTextKO   string `json:"preview_text_ko"`
	PreviewAudioPath string `json:"-" gorm:"column:preview_audio_path"` // cached mp3 under uploads/voices/
	SortOrder       int    `json:"sort_order"`
}

func (VoiceProfile) TableName() string { return "voice_profiles" }

// AIProfile is the user's custom AI lover instance.
type AIProfile struct {
	ID                    uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	UserID                int            `json:"user_id" gorm:"not null;uniqueIndex"`
	DisplayName           string         `json:"display_name" gorm:"not null"`
	CompanionGender       string         `json:"companion_gender"` // female | male
	PersonalityTemplateID string         `json:"personality_template_id"`
	SpeakingStyleTags     []string `json:"speaking_style_tags" gorm:"serializer:json"`
	VoiceProfileID        string         `json:"voice_profile_id"`
	TtsVoice              string         `json:"tts_voice" gorm:"column:tts_voice"` // Soniox name: Mina, Kenji, Emma — copied once at create
	AvatarURL             string         `json:"avatar_url"`
	ComposedPrompt        string         `json:"-" gorm:"type:text"`
	IntroMessageKO        string         `json:"intro_message_ko"`
	IntroMessageVI        string         `json:"intro_message_vi"`
	PresetSlug            string         `json:"preset_slug,omitempty"` // hani|mina|joon if quick pick
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

func (AIProfile) TableName() string { return "ai_profiles" }

// RelationshipStats tracks bond progression per user+profile.
type RelationshipStats struct {
	UserID             int       `json:"user_id" gorm:"primaryKey"`
	AIProfileID        uuid.UUID `json:"ai_profile_id" gorm:"type:uuid;primaryKey"`
	IntimacyLevel      int       `json:"intimacy_level" gorm:"default:5"`
	TrustLevel         int       `json:"trust_level" gorm:"default:5"`
	RelationshipStage  string    `json:"relationship_stage" gorm:"default:stranger"`
	DailyStreak        int       `json:"daily_streak" gorm:"default:0"`
	TotalMessages      int       `json:"total_messages" gorm:"default:0"`
	EmotionalBondScore float64   `json:"emotional_bond_score" gorm:"default:0"`
	LastInteractionAt  time.Time `json:"last_interaction_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (RelationshipStats) TableName() string { return "relationship_stats" }

// PublicProfile is API-safe profile view.
type PublicProfile struct {
	ID                    string   `json:"id"`
	DisplayName           string   `json:"display_name"`
	CompanionGender       string   `json:"companion_gender"`
	PersonalityTemplateID string   `json:"personality_template_id"`
	SpeakingStyleTags     []string `json:"speaking_style_tags"`
	VoiceProfileID        string   `json:"voice_profile_id"`
	TtsVoice              string   `json:"tts_voice"` // Soniox voice id saved with profile
	VoiceProvider         string   `json:"voice_provider,omitempty"`
	AvatarURL             string   `json:"avatar_url"`
	IntroMessageKO        string   `json:"intro_message_ko"`
	IntroMessageVI        string   `json:"intro_message_vi"`
	PresetSlug            string   `json:"preset_slug,omitempty"`
	PersonalityNameKO     string   `json:"personality_name_ko,omitempty"`
	PersonalityNameVI     string   `json:"personality_name_vi,omitempty"`
	VoiceNameKO           string   `json:"voice_name_ko,omitempty"`
}

func (p *AIProfile) ToPublic(personality *PersonalityTemplate, voice *VoiceProfile, tags []string) PublicProfile {
	out := PublicProfile{
		ID:                    p.ID.String(),
		DisplayName:           p.DisplayName,
		CompanionGender:       p.CompanionGender,
		PersonalityTemplateID: p.PersonalityTemplateID,
		SpeakingStyleTags:     tags,
		VoiceProfileID:        p.VoiceProfileID,
		TtsVoice:              p.TtsVoice,
		VoiceProvider:         "",
		AvatarURL:             p.AvatarURL,
		IntroMessageKO:        p.IntroMessageKO,
		IntroMessageVI:        p.IntroMessageVI,
		PresetSlug:            p.PresetSlug,
	}
	if personality != nil {
		out.PersonalityNameKO = personality.NameKO
		out.PersonalityNameVI = personality.NameVI
	}
	if voice != nil {
		out.VoiceNameKO = voice.NameKO
		if out.TtsVoice == "" {
			out.TtsVoice = voice.VoiceID
		}
		if out.VoiceProvider == "" {
			out.VoiceProvider = voice.Provider
		}
	}
	return out
}

// SpeakingStyleOption for onboarding UI.
type SpeakingStyleOption struct {
	ID          string `json:"id"`
	LabelKO     string `json:"label_ko"`
	LabelVI     string `json:"label_vi"`
	Description string `json:"description"`
}
