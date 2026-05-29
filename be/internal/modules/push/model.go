package push

import (
	"time"

	"github.com/google/uuid"
)

type UserDevice struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	UserID     int        `json:"user_id" gorm:"not null;index"`
	FCMToken   string     `json:"fcm_token" gorm:"uniqueIndex;not null"`
	DeviceType string     `json:"device_type" gorm:"not null"` // android | ios | web
	UserAgent  string     `json:"user_agent"`
	LastSeenAt time.Time  `json:"last_seen_at" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

func (UserDevice) TableName() string { return "user_devices" }

type Notification struct {
	ID     uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	UserID int        `json:"user_id" gorm:"not null;index"`
	Kind   string     `json:"kind" gorm:"not null;index"` // miss_1d | miss_3d | miss_7d
	Title  string     `json:"title" gorm:"not null"`
	Body   string     `json:"body" gorm:"not null"`
	SentAt time.Time  `json:"sent_at" gorm:"not null;index"`
	ReadAt *time.Time `json:"read_at,omitempty"`
}

func (Notification) TableName() string { return "notifications" }

const (
	KindMiss1Day = "miss_1d"
	KindMiss3Day = "miss_3d"
	KindMiss7Day = "miss_7d"
)
