package billing

import "errors"

var (
	ErrQuotaExceeded = errors.New("quota_exceeded")
	ErrPlanRequired  = errors.New("plan_required")
	ErrVoiceDisabled = errors.New("voice_not_allowed")
)
