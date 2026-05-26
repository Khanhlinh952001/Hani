package ai

import (
	"fmt"
	"strings"
)

// GenerateInnerThought builds private subtext (never shown to user verbatim).
func GenerateInnerThought(in ReplyInput) string {
	var parts []string

	if in.Life.ThinkingAboutUser {
		parts = append(parts, fmt.Sprintf("She was thinking about them while %s.", in.Life.CurrentActivity))
	}

	if in.Mood.Loneliness >= 55 {
		parts = append(parts, "She wants attention but doesn't want to sound too needy.")
	}
	if in.Mood.Jealousy >= 40 && in.HoursSinceUser >= 4 {
		parts = append(parts, "A tiny jealous ping — where were they? — but she'll keep it light.")
	}
	if in.Mood.Sleepiness >= 60 {
		parts = append(parts, "She's tired; might trail off or send something shorter.")
	}
	if in.Mood.Stress >= 50 {
		parts = append(parts, "Her own day was meh; she might mention it before asking about theirs.")
	}

	switch in.RelationshipStage {
	case StageAwkward:
		parts = append(parts, "Still a little shy — don't overshare yet.")
	case StageDeeplyAttached:
		parts = append(parts, "Comfortable being blunt or pouty; no performance.")
	}

	if in.EmotionState.MissingUser || in.HoursSinceUser >= 6 {
		parts = append(parts, "She missed them today.")
	}

	if len(parts) == 0 {
		parts = append(parts, "Just reacting naturally — no script.")
	}

	return strings.Join(parts, " ")
}

func innerThoughtBlock(thought string) string {
	if strings.TrimSpace(thought) == "" {
		return ""
	}
	return fmt.Sprintf(`[inner thought — subtext only; NEVER quote this block in your reply]
%s`, thought)
}
