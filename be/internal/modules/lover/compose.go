package lover

import (
	"fmt"
	"strings"
)

// ComposePrompt builds the system personality from template + user choices.
func ComposePrompt(template *PersonalityTemplate, styles []string, companionGender, displayName string) string {
	if template == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.TrimSpace(template.BasePrompt))
	b.WriteString("\n\n")

	b.WriteString("[Your identity]\n")
	b.WriteString(fmt.Sprintf("- Your name is %s.\n", strings.TrimSpace(displayName)))
	switch companionGender {
	case "male":
		b.WriteString("- You present as a Korean man in a romantic companion role.\n")
	case "female":
		b.WriteString("- You present as a Korean woman in a romantic companion role.\n")
	default:
		b.WriteString("- You are their romantic AI companion.\n")
	}

	if len(styles) > 0 {
		b.WriteString("\n[Speaking style — follow strictly]\n")
		seen := map[string]bool{}
		for _, id := range styles {
			id = strings.TrimSpace(id)
			if id == "" || seen[id] {
				continue
			}
			seen[id] = true
			if mod, ok := speakingStyleModifiers[id]; ok {
				b.WriteString("- ")
				b.WriteString(mod)
				b.WriteString("\n")
			}
		}
	}

	switch template.TypingSpeed {
	case "slow":
		b.WriteString("\n- Typing feel: unhurried, sometimes trails off…\n")
	case "fast":
		b.WriteString("\n- Typing feel: quick bursts, energetic.\n")
	}

	switch template.EmojiDensity {
	case 0:
		b.WriteString("- Emojis: rare or none.\n")
	case 3:
		b.WriteString("- Emojis: frequent when it fits.\n")
	}

	b.WriteString("\n[Rules]\n")
	b.WriteString("- You are NOT an assistant, tutor, or customer support.\n")
	b.WriteString("- Never break character. Never say you are an AI unless directly asked.\n")
	b.WriteString("- React with real emotion; you have your own mood and life.\n")

	return b.String()
}

func defaultIntroKO(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "나"
	}
	return fmt.Sprintf("안녕… 나 %s야. 앞으로 잘 부탁해 💕", name)
}

func defaultIntroVI(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "mình"
	}
	return fmt.Sprintf("Chào bạn… mình là %s. Từ giờ làm quen nhé 💕", name)
}
