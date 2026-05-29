package push

import (
	"time"

	"be/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func upsertDevice(d *UserDevice) error {
	var existing UserDevice
	err := db.DB.Where("fcm_token = ?", d.FCMToken).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"user_id":      d.UserID,
			"device_type":  d.DeviceType,
			"user_agent":   d.UserAgent,
			"last_seen_at": d.LastSeenAt,
			"revoked_at":   nil,
		}
		return db.DB.Model(&existing).Updates(updates).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return db.DB.Create(d).Error
}

func touchDevice(token string, userID int, at time.Time) error {
	res := db.DB.Model(&UserDevice{}).
		Where("fcm_token = ? AND user_id = ? AND revoked_at IS NULL", token, userID).
		Update("last_seen_at", at)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	return gorm.ErrRecordNotFound
}

func revokeDevice(token string, userID int) error {
	now := time.Now()
	return db.DB.Model(&UserDevice{}).
		Where("fcm_token = ? AND user_id = ?", token, userID).
		Update("revoked_at", now).Error
}

func activeTokensForUser(userID int) ([]string, error) {
	var tokens []string
	err := db.DB.Model(&UserDevice{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Pluck("fcm_token", &tokens).Error
	return tokens, err
}

func revokeToken(token string) error {
	now := time.Now()
	return db.DB.Model(&UserDevice{}).
		Where("fcm_token = ?", token).
		Update("revoked_at", now).Error
}

func insertNotification(n *Notification) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return db.DB.Create(n).Error
}

func sentKindToday(userID int, kind string, day time.Time) (bool, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	var count int64
	err := db.DB.Model(&Notification{}).
		Where("user_id = ? AND kind = ? AND sent_at >= ? AND sent_at < ?", userID, kind, start, end).
		Count(&count).Error
	return count > 0, err
}

type inactiveUser struct {
	UserID     int
	LastSeenAt time.Time
}

func usersInactiveSince(cutoff time.Time) ([]inactiveUser, error) {
	var rows []inactiveUser
	err := db.DB.Raw(`
		SELECT u.id AS user_id, COALESCE(MAX(d.last_seen_at), u.updated_at) AS last_seen_at
		FROM users u
		INNER JOIN user_devices d ON d.user_id = u.id AND d.revoked_at IS NULL
		WHERE u.is_active = true
		GROUP BY u.id, u.updated_at
		HAVING COALESCE(MAX(d.last_seen_at), u.updated_at) < ?
	`, cutoff).Scan(&rows).Error
	return rows, err
}
