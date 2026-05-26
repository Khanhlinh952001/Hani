package ai

import (
	"fmt"
	"strings"
	"time"
)

var koreaTZ = time.FixedZone("KST", 9*3600)

// EmotionalState is Hani's internal mood for this turn (not shown to user verbatim).
type EmotionalState struct {
	Mood           string
	Energy         string
	Attachment     int
	MissingUser    bool
	Jealousy       int
	Shyness        int
	CurrentFeeling string
}

// RelationshipStage affects tone, pet names, and intimacy level.
type RelationshipStage string

const (
	StageAwkward         RelationshipStage = "awkward"
	StageClose           RelationshipStage = "close"
	StageRomantic        RelationshipStage = "romantic"
	StageDeeplyAttached  RelationshipStage = "deeply_attached"
)

func DefaultEmotionalState() EmotionalState {
	return EmotionalState{
		Mood:           "soft",
		Energy:         "medium",
		Attachment:     55,
		MissingUser:    false,
		Jealousy:       8,
		Shyness:        35,
		CurrentFeeling: "그냥 보고 싶어",
	}
}

func RelationshipStageFromMessageCount(n int) RelationshipStage {
	switch {
	case n >= 400:
		return StageDeeplyAttached
	case n >= 120:
		return StageRomantic
	case n >= 30:
		return StageClose
	default:
		return StageAwkward
	}
}

func (s RelationshipStage) PromptBlock() string {
	switch s {
	case StageAwkward:
		return `[relationship — awkward early stage]
- Still getting comfortable; a little shy, not too clingy yet
- Pet names rare; teasing is light`
	case StageClose:
		return `[relationship — close couple]
- Comfortable banter; occasional 자기/여보 okay but not every line
- Inside jokes and callbacks welcome`
	case StageRomantic:
		return `[relationship — romantic]
- Warm, misses them easily; soft jealousy sometimes
- Casual intimacy, teasing, short emotional lines`
	case StageDeeplyAttached:
		return `[relationship — deeply attached]
- Very comfortable; can be blunt, clingy, or pouty naturally
- Speaks like someone who knows them deeply — no performance`
	default:
		return ""
	}
}

func (e EmotionalState) PromptBlock() string {
	missing := "no"
	if e.MissingUser {
		missing = "yes"
	}
	return fmt.Sprintf(`[Hani emotional state — feel this internally; never list these labels in your reply]
- mood: %s
- energy: %s
- attachment: %d/100
- missing user: %s
- jealousy: %d/100
- shyness: %d/100
- feeling now: %s`,
		e.Mood, e.energyLabel(), clamp100(e.Attachment), missing,
		clamp100(e.Jealousy), clamp100(e.Shyness), e.CurrentFeeling)
}

func (e EmotionalState) energyLabel() string {
	if e.Energy != "" {
		return e.Energy
	}
	return "medium"
}

func clamp100(n int) int {
	if n < 0 {
		return 0
	}
	if n > 100 {
		return 100
	}
	return n
}

func FormatTimeContext(now time.Time) string {
	t := now.In(koreaTZ)
	hour := t.Hour()
	var vibe string
	switch {
	case hour >= 23 || hour < 6:
		vibe = "late night — sleepy, softer tone"
	case hour < 10:
		vibe = "morning — slow wake-up energy"
	case hour < 14:
		vibe = "daytime — casual"
	case hour < 19:
		vibe = "afternoon/evening — winding down"
	default:
		vibe = "night — relaxed, more affectionate"
	}
	return fmt.Sprintf("[time] %s (%s KST)", vibe, t.Format("15:04"))
}

// BootstrapEmotion sets initial state when a websocket session starts.
func BootstrapEmotion(recent []Turn, messageCount int, now time.Time) EmotionalState {
	state := DefaultEmotionalState()
	applyTimeOfDay(&state, now)

	if messageCount > 80 {
		state.Attachment = 72
		state.Shyness = 20
	}
	if messageCount > 200 {
		state.Attachment = 84
		state.Shyness = 12
	}

	if len(recent) == 0 {
		state.MissingUser = true
		state.CurrentFeeling = "오랜만에 연락 올까 기대돼"
		return state
	}

	if lastUserTurn(recent) == "" {
		state.MissingUser = true
		state.Attachment = clamp100(state.Attachment + 5)
		state.CurrentFeeling = "왜 안 와… 좀 보고 싶어"
	}

	return state
}

// EvolveEmotionAfterExchange nudges state from the latest exchange (heuristic, no extra LLM call).
func EvolveEmotionAfterExchange(state EmotionalState, userMsg, haniMsg string) EmotionalState {
	u := strings.TrimSpace(userMsg)
	if u == "" {
		return state
	}

	state.MissingUser = false

	if len([]rune(u)) <= 12 {
		state.Energy = "high"
		state.Mood = "playful"
	}

	lower := strings.ToLower(u)
	for _, kw := range []string{"보고", "사랑", "그립", "보고싶", "싶어", "고마", "miss"} {
		if strings.Contains(lower, kw) || strings.Contains(u, kw) {
			state.Attachment = clamp100(state.Attachment + 4)
			state.Mood = "warm"
			state.CurrentFeeling = "마음이 따뜻해"
			state.Shyness = clamp100(state.Shyness - 3)
			break
		}
	}

	for _, kw := range []string{"바빠", "못 ", "못해", "못 와", "나중에", "피곤", "미안"} {
		if strings.Contains(u, kw) {
			state.Mood = "soft"
			state.CurrentFeeling = "아쉽지만 기다릴게"
			break
		}
	}

	if strings.Contains(u, "ㅋ") || strings.Contains(u, "ㅎ") {
		state.Mood = "playful"
		state.Energy = "high"
	}

	if len([]rune(haniMsg)) <= 20 {
		state.Energy = "low"
	}

	return state
}

func applyTimeOfDay(state *EmotionalState, now time.Time) {
	hour := now.In(koreaTZ).Hour()
	switch {
	case hour >= 23 || hour < 6:
		state.Energy = "low"
		state.Mood = "sleepy"
		state.CurrentFeeling = "졸린데 너 생각나"
	case hour < 10:
		state.Energy = "medium"
		state.Mood = "soft"
	case hour >= 22:
		state.Energy = "low"
		state.Mood = "clingy"
		state.Attachment = clamp100(state.Attachment + 3)
	}
}

func lastUserTurn(turns []Turn) string {
	for i := len(turns) - 1; i >= 0; i-- {
		if turns[i].Role == "user" {
			return turns[i].Content
		}
	}
	return ""
}
