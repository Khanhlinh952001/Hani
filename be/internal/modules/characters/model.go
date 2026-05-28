package characters

import "time"

// Character is a selectable AI companion persona.
type Character struct {
	ID                string    `json:"id" gorm:"primaryKey"` // slug: hani, mina, joon
	Name              string    `json:"name" gorm:"not null"`
	DisplayName       string    `json:"display_name"`
	Gender            string    `json:"gender"` // female | male
	AvatarURL         string    `json:"avatar_url"`
	PersonalityPrompt string    `json:"-" gorm:"type:text"`
	VoiceProvider     string    `json:"voice_provider"` // openai | soniox
	VoiceID           string    `json:"voice_id"`
	SonioxVoice       string    `json:"soniox_voice,omitempty"`
	TTSLanguage       string    `json:"tts_language,omitempty"`
	IntroMessageKO    string    `json:"intro_message_ko"`
	IntroMessageVI    string    `json:"intro_message_vi"`
	SpeakingStyle     string    `json:"speaking_style"`
	EmotionStyle      string    `json:"emotion_style"`
	EmojiStyle        string    `json:"emoji_style,omitempty"`
	TypingStyle       string    `json:"typing_style,omitempty"`
	SortOrder         int       `json:"sort_order"`
	CreatedAt         time.Time `json:"created_at"`
}

func (Character) TableName() string {
	return "ai_characters"
}

// UserCharacterMemory tracks per-user relationship with a character.
type UserCharacterMemory struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	UserID             int       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_character"`
	CharacterID        string    `json:"character_id" gorm:"not null;uniqueIndex:idx_user_character"`
	IntimacyLevel      int       `json:"intimacy_level" gorm:"default:1"`
	RelationshipStatus string    `json:"relationship_status" gorm:"default:new"`
	LastInteractionAt  time.Time `json:"last_interaction_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (UserCharacterMemory) TableName() string {
	return "user_character_memory"
}

// PublicCharacter is API-safe (no personality prompt).
type PublicCharacter struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Gender         string `json:"gender"`
	AvatarURL      string `json:"avatar_url"`
	IntroMessageKO string `json:"intro_message_ko"`
	IntroMessageVI string `json:"intro_message_vi"`
	SpeakingStyle  string `json:"speaking_style"`
	EmotionStyle   string `json:"emotion_style"`
	EmojiStyle     string `json:"emoji_style,omitempty"`
	TypingStyle    string `json:"typing_style,omitempty"`
	VoiceProvider  string `json:"voice_provider"`
	VoiceID        string `json:"voice_id"`
}

func (c *Character) ToPublic() PublicCharacter {
	return PublicCharacter{
		ID:             c.ID,
		Name:           c.Name,
		DisplayName:    c.DisplayName,
		Gender:         c.Gender,
		AvatarURL:      c.AvatarURL,
		IntroMessageKO: c.IntroMessageKO,
		IntroMessageVI: c.IntroMessageVI,
		SpeakingStyle:  c.SpeakingStyle,
		EmotionStyle:   c.EmotionStyle,
		EmojiStyle:     c.EmojiStyle,
		TypingStyle:    c.TypingStyle,
		VoiceProvider:  c.VoiceProvider,
		VoiceID:        c.VoiceID,
	}
}
