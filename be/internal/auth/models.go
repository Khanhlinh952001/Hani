package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     *int       `gorm:"index"`
	GuestID    *uuid.UUID `gorm:"type:uuid;index"`
	RevokedAt  *time.Time
	LastSeenAt time.Time
	CreatedAt  time.Time
}

func (AuthSession) TableName() string { return "auth_sessions" }

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SessionID uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }
