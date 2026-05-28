package billing

import (
	"os"

	"be/internal/modules/users"
)

func AdminBypassQuota() bool {
	return os.Getenv("ADMIN_BYPASS_QUOTA") == "true"
}

func PlanForUser(u *users.User) string {
	if u == nil {
		return PlanGuest
	}
	if users.IsAdmin(u.Role) && AdminBypassQuota() {
		return PlanPremium
	}
	if u.SubscriptionPlan != "" {
		return u.SubscriptionPlan
	}
	return PlanFree
}

func GetPlanLimits(plan string) (PlanLimit, error) {
	var lim PlanLimit
	if err := dbGetPlanLimit(plan, &lim); err != nil {
		if plan != PlanFree {
			return GetPlanLimits(PlanFree)
		}
		return lim, err
	}
	return lim, nil
}

func AllowsVoice(plan string) bool {
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return false
	}
	return lim.AllowVoice
}

func AllowsMemory(plan string) bool {
	lim, err := GetPlanLimits(plan)
	if err != nil {
		return false
	}
	return lim.AllowMemory
}
