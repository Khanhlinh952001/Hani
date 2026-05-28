package billing

import (
	"strconv"

	"be/internal/db"
	"be/internal/modules/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsageSnapshot struct {
	Plan              string `json:"plan"`
	DailyMessages     int    `json:"daily_messages"`
	DailyMessagesLimit *int  `json:"daily_messages_limit,omitempty"`
	DailyVoiceSeconds int    `json:"daily_voice_seconds"`
	DailyVoiceLimit   *int   `json:"daily_voice_limit,omitempty"`
	Warning           bool   `json:"warning"`
}

func GetUsageSnapshot(userID int) (*UsageSnapshot, error) {
	u, err := users.GetUserByIDService(itoa(userID))
	if err != nil {
		return nil, err
	}
	return snapshotForPlan(PlanForUser(u), userID, nil)
}

func GetGuestUsageSnapshot(guestID uuid.UUID) (*UsageSnapshot, error) {
	return snapshotForPlan(PlanGuest, 0, &guestID)
}

func snapshotForPlan(plan string, userID int, guestID *uuid.UUID) (*UsageSnapshot, error) {
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return nil, err
	}
	day := todayUTC()
	var row UserUsage
	q := db.DB.Where("period_date = ?", day)
	if guestID != nil {
		q = q.Where("guest_id = ?", *guestID)
	} else {
		q = q.Where("user_id = ?", userID)
	}
	_ = q.First(&row).Error

	snap := &UsageSnapshot{
		Plan:              plan,
		DailyMessages:     row.DailyMessages,
		DailyMessagesLimit: lim.DailyMessages,
		DailyVoiceSeconds: row.DailyVoiceSeconds,
		DailyVoiceLimit:   lim.DailyVoiceSeconds,
	}
	if lim.DailyMessages != nil && *lim.DailyMessages > 0 {
		used := float64(row.DailyMessages) / float64(*lim.DailyMessages)
		snap.Warning = used >= 0.8 && used < 1.0
	}
	return snap, nil
}

// ConsumeMessage increments daily message count; returns ErrQuotaExceeded when over limit.
func ConsumeMessage(userID int, plan string) error {
	if users.IsAdminRole(userID) && AdminBypassQuota() {
		return nil
	}
	return consumeMessage(plan, &userID, nil)
}

func ConsumeGuestMessage(guestID uuid.UUID) error {
	return consumeMessage(PlanGuest, nil, &guestID)
}

func consumeMessage(plan string, userID *int, guestID *uuid.UUID) error {
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return err
	}
	if lim.DailyMessages == nil {
		return recordMessage(userID, guestID, plan)
	}

	day := todayUTC()
	return db.DB.Transaction(func(tx *gorm.DB) error {
		row, err := getOrCreateUsage(tx, userID, guestID, day)
		if err != nil {
			return err
		}
		resetUsageIfNewDay(row, day)
		if row.DailyMessages >= *lim.DailyMessages {
			return ErrQuotaExceeded
		}
		row.DailyMessages++
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		return tx.Create(&UsageLog{
			UserID:    userID,
			GuestID:   guestID,
			EventType: "chat_message",
			Units:     1,
		}).Error
	})
}

func recordMessage(userID *int, guestID *uuid.UUID, plan string) error {
	day := todayUTC()
	return db.DB.Transaction(func(tx *gorm.DB) error {
		row, err := getOrCreateUsage(tx, userID, guestID, day)
		if err != nil {
			return err
		}
		resetUsageIfNewDay(row, day)
		row.DailyMessages++
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		return tx.Create(&UsageLog{
			UserID:    userID,
			GuestID:   guestID,
			EventType: "chat_message",
			Units:     1,
		}).Error
	})
}

// AddVoiceSeconds adds voice usage after a session segment.
func AddVoiceSeconds(userID int, seconds int, plan string) error {
	if seconds <= 0 {
		return nil
	}
	if users.IsAdminRole(userID) && AdminBypassQuota() {
		return nil
	}
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return err
	}
	if !lim.AllowVoice {
		return ErrVoiceDisabled
	}

	day := todayUTC()
	uid := userID
	return db.DB.Transaction(func(tx *gorm.DB) error {
		row, err := getOrCreateUsage(tx, &uid, nil, day)
		if err != nil {
			return err
		}
		resetUsageIfNewDay(row, day)
		if lim.DailyVoiceSeconds != nil && row.DailyVoiceSeconds >= *lim.DailyVoiceSeconds {
			return ErrQuotaExceeded
		}
		row.DailyVoiceSeconds += seconds
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		return tx.Create(&UsageLog{
			UserID:    &uid,
			EventType: "voice_seconds",
			Units:     seconds,
		}).Error
	})
}

// CheckVoiceAllowed verifies voice can start (quota not already exhausted).
func CheckVoiceAllowed(userID int, plan string) error {
	if users.IsAdminRole(userID) && AdminBypassQuota() {
		return nil
	}
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return err
	}
	if !lim.AllowVoice {
		return ErrVoiceDisabled
	}
	if lim.DailyVoiceSeconds == nil {
		return nil
	}
	day := todayUTC()
	var row UserUsage
	if err := db.DB.Where("user_id = ? AND period_date = ?", userID, day).First(&row).Error; err != nil {
		return nil
	}
	if row.DailyVoiceSeconds >= *lim.DailyVoiceSeconds {
		return ErrQuotaExceeded
	}
	return nil
}

func ResetUserUsage(userID int) error {
	day := todayUTC()
	res := db.DB.Model(&UserUsage{}).
		Where("user_id = ? AND period_date = ?", userID, day).
		Updates(map[string]interface{}{
			"daily_messages":      0,
			"daily_voice_seconds": 0,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	uid := userID
	return db.DB.Create(&UserUsage{
		UserID:     &uid,
		PeriodDate: day,
	}).Error
}

func CreateGuest() (*Guest, error) {
	g := &Guest{ID: uuid.New()}
	if err := db.DB.Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func itoa(id int) string {
	return strconv.Itoa(id)
}
