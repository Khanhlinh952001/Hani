package billing

import (
	"time"

	"be/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func dbGetPlanLimit(plan string, dest *PlanLimit) error {
	return db.DB.First(dest, "plan = ?", plan).Error
}

func dbListPlans(dest *[]PlanLimit) error {
	return db.DB.Order("plan").Find(dest).Error
}

func todayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func getOrCreateUsage(tx *gorm.DB, userID *int, guestID *uuid.UUID, day time.Time) (*UserUsage, error) {
	var row UserUsage
	q := tx.Where("period_date = ?", day)
	if userID != nil {
		q = q.Where("user_id = ?", *userID)
	} else if guestID != nil {
		q = q.Where("guest_id = ?", *guestID)
	}
	err := q.First(&row).Error
	if err == nil {
		return &row, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = UserUsage{
		UserID:     userID,
		GuestID:    guestID,
		PeriodDate: day,
	}
	if err := tx.Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func resetUsageIfNewDay(row *UserUsage, day time.Time) {
	if row.PeriodDate.Equal(day) {
		return
	}
	row.PeriodDate = day
	row.DailyMessages = 0
	row.DailyVoiceSeconds = 0
}
