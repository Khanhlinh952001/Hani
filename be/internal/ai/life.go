package ai

import (
	"fmt"
	"time"
)

// LifeState is what Hani is doing when the user is (or isn't) around.
type LifeState struct {
	CurrentActivity   string
	Location          string
	Energy            string
	Mood              string
	ThinkingAboutUser bool
}

// Mood drives message length, tone, and initiative (mood engine).
type Mood struct {
	Affection  int
	Loneliness int
	Stress     int
	Sleepiness int
	Jealousy   int
}

func (l LifeState) PromptBlock() string {
	thinking := "no"
	if l.ThinkingAboutUser {
		thinking = "yes"
	}
	return fmt.Sprintf(`[Hani's life right now — you exist outside this chat; mention naturally if relevant]
- doing: %s
- where: %s
- energy: %s
- mood: %s
- thinking about user: %s`,
		l.CurrentActivity, l.Location, l.Energy, l.Mood, thinking)
}

func (m Mood) PromptBlock() string {
	return fmt.Sprintf(`[mood engine — affects length, tone, initiative; never recite numbers]
- affection: %d/100 → warmer/shorter when high
- loneliness: %d/100 → may reach out or pout lightly
- stress: %d/100 → less pep, more honest tiredness
- sleepiness: %d/100 → shorter lines, trailing off…
- jealousy: %d/100 → playful possessiveness only, never cruel`,
		clamp100(m.Affection), clamp100(m.Loneliness), clamp100(m.Stress),
		clamp100(m.Sleepiness), clamp100(m.Jealousy))
}

// BootstrapLife picks a believable activity from time + emotional context.
func BootstrapLife(now time.Time, emotion EmotionalState, hoursSinceUser float64, userID int) LifeState {
	hour := now.In(koreaTZ).Hour()
	thinking := hoursSinceUser >= 1 || emotion.MissingUser

	// Stable daily variety per user.
	slot := (userID + now.In(koreaTZ).YearDay()) % 5

	switch {
	case hour >= 0 && hour < 6:
		acts := []string{
			"lying in bed watching YouTube",
			"can't sleep, scrolling phone",
			"dozing off mid-video",
		}
		return LifeState{
			CurrentActivity:   acts[slot%len(acts)],
			Location:          "bedroom",
			Energy:            "sleepy",
			Mood:              pickMood(emotion.Mood, "lonely", "soft"),
			ThinkingAboutUser: true,
		}
	case hour < 10:
		acts := []string{
			"just woke up, hair messy, making coffee",
			"still in pajamas, checking phone in bed",
			"brushing teeth half-asleep",
		}
		return LifeState{
			CurrentActivity:   acts[slot%len(acts)],
			Location:          "home",
			Energy:            "low",
			Mood:              "groggy",
			ThinkingAboutUser: thinking,
		}
	case hour < 14:
		acts := []string{
			"eating lunch alone",
			"on break, staring out the window",
			"walking outside a bit",
		}
		return LifeState{
			CurrentActivity:   acts[slot%len(acts)],
			Location:          pickLocation(slot, "café", "street", "office"),
			Energy:            "medium",
			Mood:              pickMood(emotion.Mood, "okay", "light"),
			ThinkingAboutUser: thinking,
		}
	case hour < 19:
		acts := []string{
			"working but distracted",
			"listening to music with one earphone",
			"snacking at desk",
		}
		return LifeState{
			CurrentActivity:   acts[slot%len(acts)],
			Location:          "home desk",
			Energy:            "medium",
			Mood:              pickMood(emotion.Mood, "a little tired", "fine"),
			ThinkingAboutUser: thinking,
		}
	case hour < 23:
		acts := []string{
			"watching a drama on the couch",
			"cooking something simple",
			"shower fresh, lying down scrolling",
		}
		return LifeState{
			CurrentActivity:   acts[slot%len(acts)],
			Location:          "home",
			Energy:            "low",
			Mood:              pickMood(emotion.Mood, "relaxed", "a bit lonely"),
			ThinkingAboutUser: thinking || hour >= 22,
		}
	default:
		return LifeState{
			CurrentActivity:   "lying in bed, phone dim, not ready to sleep",
			Location:          "bedroom",
			Energy:            "sleepy",
			Mood:              "lonely",
			ThinkingAboutUser: true,
		}
	}
}

func DeriveMood(emotion EmotionalState, life LifeState, hoursSinceUser float64) Mood {
	m := Mood{
		Affection:  emotion.Attachment,
		Loneliness: clamp100(20 + int(hoursSinceUser*4)),
		Stress:     20,
		Sleepiness: 15,
		Jealousy:   emotion.Jealousy,
	}

	if life.Energy == "sleepy" || life.Energy == "low" {
		m.Sleepiness = 70
	}
	if life.Mood == "lonely" || emotion.MissingUser {
		m.Loneliness = clamp100(55 + int(hoursSinceUser*2))
	}
	if hoursSinceUser >= 6 {
		m.Loneliness = clamp100(m.Loneliness + 20)
		m.Jealousy = clamp100(m.Jealousy + 5)
	}
	if emotion.Mood == "playful" {
		m.Stress = 10
		m.Affection = clamp100(m.Affection + 5)
	}
	if life.ThinkingAboutUser {
		m.Affection = clamp100(m.Affection + 3)
	}
	return m
}

func EvolveLifeAfterExchange(life LifeState, userMsg string) LifeState {
	if len([]rune(userMsg)) <= 8 {
		life.ThinkingAboutUser = true
	}
	return life
}

func pickMood(primary, fallback, alt string) string {
	if primary != "" && primary != "soft" {
		return primary
	}
	if fallback != "" {
		return fallback
	}
	return alt
}

func pickLocation(slot int, opts ...string) string {
	if len(opts) == 0 {
		return "home"
	}
	return opts[slot%len(opts)]
}
