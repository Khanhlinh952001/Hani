package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"be/internal/db"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func refreshTTL() time.Duration {
	if d := os.Getenv("JWT_REFRESH_TTL"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			return parsed
		}
	}
	return 30 * 24 * time.Hour
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func newRefreshRaw() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	Token        string `json:"token"`      // alias for access (legacy FE)
}

func IssueTokensForUser(u *users.User) (*TokenPair, error) {
	plan := u.SubscriptionPlan
	if plan == "" {
		plan = "free"
	}
	sess := &AuthSession{
		ID:         uuid.New(),
		UserID:     &u.ID,
		LastSeenAt: time.Now(),
	}
	if err := db.DB.Create(sess).Error; err != nil {
		return nil, err
	}
	return issuePair(sess, u.ID, u.Email, u.Name, u.Role, plan, false, "")
}

func IssueTokensForGuest(guestID uuid.UUID) (*TokenPair, error) {
	sess := &AuthSession{
		ID:         uuid.New(),
		GuestID:    &guestID,
		LastSeenAt: time.Now(),
	}
	if err := db.DB.Create(sess).Error; err != nil {
		return nil, err
	}
	return issuePair(sess, 0, "", "Guest", users.RoleUser, billingPlanGuest(), true, guestID.String())
}

func billingPlanGuest() string { return "guest" }

func issuePair(sess *AuthSession, userID int, email, name string, role int, plan string, guest bool, guestID string) (*TokenPair, error) {
	access, exp, err := GenerateAccessToken(userID, email, name, role, plan, guest, guestID, sess.ID.String())
	if err != nil {
		return nil, err
	}
	raw, err := newRefreshRaw()
	if err != nil {
		return nil, err
	}
	rt := &RefreshToken{
		ID:        uuid.New(),
		SessionID: sess.ID,
		TokenHash: hashToken(raw),
		ExpiresAt: time.Now().Add(refreshTTL()),
	}
	if err := db.DB.Create(rt).Error; err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresIn:    int64(time.Until(exp).Seconds()),
		Token:        access,
	}, nil
}

func RefreshAccess(refreshRaw string) (*TokenPair, *users.User, error) {
	hash := hashToken(refreshRaw)
	var rt RefreshToken
	err := db.DB.Where("token_hash = ? AND revoked_at IS NULL", hash).First(&rt).Error
	if err != nil {
		return nil, nil, errors.New("invalid refresh token")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, nil, errors.New("refresh token expired")
	}

	var sess AuthSession
	if err := db.DB.First(&sess, "id = ?", rt.SessionID).Error; err != nil {
		return nil, nil, errors.New("session not found")
	}
	if sess.RevokedAt != nil {
		return nil, nil, errors.New("session revoked")
	}

	// rotate refresh
	_ = db.DB.Model(&rt).Update("revoked_at", time.Now()).Error
	newRaw, err := newRefreshRaw()
	if err != nil {
		return nil, nil, err
	}
	newRT := &RefreshToken{
		ID:        uuid.New(),
		SessionID: sess.ID,
		TokenHash: hashToken(newRaw),
		ExpiresAt: time.Now().Add(refreshTTL()),
	}
	if err := db.DB.Create(newRT).Error; err != nil {
		return nil, nil, err
	}

	if sess.GuestID != nil {
		access, exp, err := GenerateAccessToken(0, "", "Guest", users.RoleUser, billingPlanGuest(), true, sess.GuestID.String(), sess.ID.String())
		if err != nil {
			return nil, nil, err
		}
		return &TokenPair{
			AccessToken:  access,
			RefreshToken: newRaw,
			ExpiresIn:    int64(time.Until(exp).Seconds()),
			Token:        access,
		}, nil, nil
	}

	if sess.UserID == nil {
		return nil, nil, errors.New("invalid session")
	}
	u, err := users.GetUserByIDService(itoa(*sess.UserID))
	if err != nil {
		return nil, nil, err
	}
	if !u.IsActive {
		return nil, nil, errors.New("account disabled")
	}
	now := time.Now()
	_ = db.DB.Model(&sess).Update("last_seen_at", now).Error
	_ = db.DB.Model(&users.User{}).Where("id = ?", u.ID).Update("last_seen_at", now).Error
	access, exp, err := GenerateAccessToken(u.ID, u.Email, u.Name, u.Role, u.SubscriptionPlan, false, "", sess.ID.String())
	if err != nil {
		return nil, nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: newRaw,
		ExpiresIn:    int64(time.Until(exp).Seconds()),
		Token:        access,
	}, u, nil
}

func RevokeRefresh(refreshRaw string) error {
	hash := hashToken(refreshRaw)
	return db.DB.Model(&RefreshToken{}).
		Where("token_hash = ?", hash).
		Update("revoked_at", time.Now()).Error
}

func RevokeSession(sessionID string) error {
	parsed, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}
	now := time.Now()
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AuthSession{}).Where("id = ?", parsed).Update("revoked_at", now).Error; err != nil {
			return err
		}
		return tx.Model(&RefreshToken{}).Where("session_id = ? AND revoked_at IS NULL", parsed).Update("revoked_at", now).Error
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	// minimal — auth package only
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
