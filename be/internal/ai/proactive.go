package ai

import "time"

// ProactiveKind is why Hani reaches out first (not a chatbot waiting passively).
type ProactiveKind int

const (
	ProactiveNone ProactiveKind = iota
	ProactiveMissedYou
	ProactiveLateNight
)

func DecideProactive(hoursSinceUser float64, now time.Time, hasHistory bool) ProactiveKind {
	if !hasHistory {
		return ProactiveNone
	}
	hour := now.In(koreaTZ).Hour()
	if hoursSinceUser >= 6 {
		return ProactiveMissedYou
	}
	if (hour >= 23 || hour < 2) && hoursSinceUser >= 2 {
		return ProactiveLateNight
	}
	return ProactiveNone
}

func (k ProactiveKind) PromptHint() string {
	switch k {
	case ProactiveMissedYou:
		return `[proactive reach-out — they came back after a long gap]
You noticed they were gone for hours. ONE short line — curious, a little pouty, not angry.
Examples of vibe (do NOT copy): "뭐해...? 오늘 바빠? ㅠㅠ" / "어디 갔었어 ㅋㅋ"
Do NOT lecture. Do NOT start with 여보.`
	case ProactiveLateNight:
		return `[proactive — they're online late at night]
You caught them up late. ONE short teasing or worried line.
Examples of vibe (do NOT copy): "안 자고 있었네 ㅋㅋ" / "또 늦게까지… 졸리지?"
Do NOT start with 여보.`
	default:
		return ""
	}
}
