package billing

import (
	"be/internal/db"
	"log"
)

const (
	PlanGuest    = "guest"
	PlanFree     = "free"
	PlanPlus     = "plus"
	PlanPremium  = "premium"
)

func intPtr(n int) *int { return &n }

func SeedPlanLimits() {
	rows := []PlanLimit{
		{Plan: PlanGuest, DailyMessages: intPtr(5), DailyVoiceSeconds: intPtr(0), AllowVoice: false, AllowMemory: false, AllowPremiumVoices: false, MaxMemories: intPtr(0)},
		{Plan: PlanFree, DailyMessages: intPtr(30), DailyVoiceSeconds: intPtr(300), AllowVoice: true, AllowMemory: true, AllowPremiumVoices: false, MaxMemories: intPtr(50)},
		{Plan: PlanPlus, DailyMessages: intPtr(1000), DailyVoiceSeconds: intPtr(3600), AllowVoice: true, AllowMemory: true, AllowPremiumVoices: true, MaxMemories: intPtr(500)},
		{Plan: PlanPremium, DailyMessages: nil, DailyVoiceSeconds: nil, AllowVoice: true, AllowMemory: true, AllowPremiumVoices: true, MaxMemories: nil},
	}
	for _, row := range rows {
		var existing PlanLimit
		if db.DB.First(&existing, "plan = ?", row.Plan).Error == nil {
			_ = db.DB.Model(&existing).Updates(map[string]interface{}{
				"daily_messages":        row.DailyMessages,
				"daily_voice_seconds":   row.DailyVoiceSeconds,
				"allow_voice":           row.AllowVoice,
				"allow_memory":          row.AllowMemory,
				"allow_premium_voices":  row.AllowPremiumVoices,
				"max_memories":          row.MaxMemories,
			}).Error
			continue
		}
		if err := db.DB.Create(&row).Error; err != nil {
			log.Println("seed plan_limits:", row.Plan, err)
		}
	}
}
