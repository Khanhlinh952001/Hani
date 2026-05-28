package billing

import (
	"time"

	"github.com/google/uuid"
)

type PlanLimit struct {
	Plan                 string `json:"plan" gorm:"primaryKey"`
	DailyMessages        *int   `json:"daily_messages" gorm:"column:daily_messages"`
	DailyVoiceSeconds    *int   `json:"daily_voice_seconds" gorm:"column:daily_voice_seconds"`
	AllowVoice           bool   `json:"allow_voice" gorm:"column:allow_voice"`
	AllowMemory          bool   `json:"allow_memory" gorm:"column:allow_memory"`
	AllowPremiumVoices   bool   `json:"allow_premium_voices" gorm:"column:allow_premium_voices"`
	MaxMemories          *int   `json:"max_memories" gorm:"column:max_memories"`
	UpdatedAt            time.Time
}

func (PlanLimit) TableName() string { return "plan_limits" }

type UserUsage struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	UserID            *int       `json:"user_id" gorm:"uniqueIndex:idx_usage_user_day"`
	GuestID           *uuid.UUID `json:"guest_id" gorm:"type:uuid;uniqueIndex:idx_usage_guest_day"`
	PeriodDate        time.Time  `json:"period_date" gorm:"type:date;uniqueIndex:idx_usage_user_day;uniqueIndex:idx_usage_guest_day"`
	DailyMessages     int       `json:"daily_messages"`
	DailyVoiceSeconds int       `json:"daily_voice_seconds"`
	TokensIn          int64     `json:"tokens_in"`
	TokensOut         int64     `json:"tokens_out"`
	UpdatedAt         time.Time
}

func (UserUsage) TableName() string { return "user_usage" }

type UsageLog struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	UserID    *int       `json:"user_id"`
	GuestID   *uuid.UUID `json:"guest_id" gorm:"type:uuid"`
	EventType string     `json:"event_type" gorm:"size:32"`
	Units     int        `json:"units"`
	CreatedAt time.Time  `json:"created_at"`
}

func (UsageLog) TableName() string { return "usage_logs" }

type Guest struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (Guest) TableName() string { return "guests" }
