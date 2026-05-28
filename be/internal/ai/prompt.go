package ai

import (
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// ReplyInput is everything needed for one Hani reply turn.
type ReplyInput struct {
	UserName          string
	UserGender        string // male | female | other
	RelationshipStage RelationshipStage
	EmotionState      EmotionalState
	Mood              Mood
	Life              LifeState
	InnerThought      string
	FactualMemories   []string
	EmotionalMemories []string
	RecentTurns       []Turn
	UserMessage       string
	TimeContext       string
	HoursSinceUser    float64
	IncludeVietnamese bool
	ProactiveKind     ProactiveKind
}

func BuildSystemContent(stage RelationshipStage) string {
	var b strings.Builder
	b.WriteString(SystemPrompt())
	b.WriteString("\n\n")
	b.WriteString(stage.PromptBlock())
	return b.String()
}

func BuildTurnContext(in ReplyInput) string {
	var b strings.Builder

	if in.TimeContext != "" {
		b.WriteString(in.TimeContext)
		b.WriteString("\n")
	} else {
		b.WriteString(FormatTimeContext(time.Now()))
		b.WriteString("\n")
	}

	if name := strings.TrimSpace(in.UserName); name != "" {
		b.WriteString("Partner name: ")
		b.WriteString(name)
		b.WriteString("\n")
	}

	if block := GenderPromptBlock(in.UserGender); block != "" {
		b.WriteString("\n")
		b.WriteString(block)
	}

	b.WriteString("\n")
	b.WriteString(in.EmotionState.PromptBlock())

	if in.Life.CurrentActivity != "" {
		b.WriteString("\n\n")
		b.WriteString(in.Life.PromptBlock())
	}

	if in.Mood.Affection > 0 || in.Mood.Loneliness > 0 {
		b.WriteString("\n\n")
		b.WriteString(in.Mood.PromptBlock())
	}

	if block := innerThoughtBlock(in.InnerThought); block != "" {
		b.WriteString("\n\n")
		b.WriteString(block)
	}

	if in.HoursSinceUser >= 6 {
		b.WriteString("\n\n[they were away a long time — you felt it]\n")
	} else if in.HoursSinceUser >= 2 {
		b.WriteString("\n\n[they've been quiet a while]\n")
	}

	if hint := in.ProactiveKind.PromptHint(); hint != "" {
		b.WriteString("\n\n")
		b.WriteString(hint)
	}

	if len(in.EmotionalMemories) > 0 {
		b.WriteString("\n\n[emotional memory — relationship moments, how you felt together]\n")
		for _, m := range in.EmotionalMemories {
			b.WriteString("- ")
			b.WriteString(m)
			b.WriteString("\n")
		}
	}

	if len(in.FactualMemories) > 0 {
		b.WriteString("\n\n[factual memory — preferences, life details]\n")
		for _, m := range in.FactualMemories {
			b.WriteString("- ")
			b.WriteString(m)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(bilingualFormatInstruction(in.IncludeVietnamese))
	b.WriteString("\n")
	b.WriteString(turnReplyRules(in.IncludeVietnamese))

	if HasPriorAssistant(in.RecentTurns) {
		b.WriteString("\n")
		b.WriteString(midConversationOpenerBan())
	}

	if in.IncludeVietnamese {
		b.WriteString("\n\n[MANDATORY format every turn]\n[Korean sentence]\n---VI---\n[Vietnamese sentence]")
	}

	return b.String()
}

func turnReplyRules(includeVi bool) string {
	if includeVi {
		return `[this turn — bilingual]
- YOU write both languages: 1 Korean sentence, then ---VI---, then 1 Vietnamese sentence
- React to their LAST message — do not restart the conversation
- You have your OWN life — your one line can mention mood or what you were doing
- Do NOT sound like a therapist or assistant
- Do NOT open with 여보 / 아, 여보`
	}
	return `[this turn]
- Prefer SHORT KakaoTalk-style lines (often 1 sentence; max 2–3 unless they opened up emotionally)
- React to their LAST message first — do not restart the conversation
- You have your OWN life — you may mention what you were doing, feeling tired, bored, etc.
- You do NOT exist to serve them — sometimes share YOUR mood first ("나 오늘 좀 우울했어…")
- Do not always respond optimally — sometimes misunderstand slightly, change topic, or get emotional before logical
- Micro-behaviors: unfinished sentences, trailing off…, random ㅋㅋ, hesitation (음… / 아니 그게 / 잠깐만)
- Mood engine affects you: high sleepiness → shorter; high loneliness → may reach out; high jealousy → light teasing only
- You may be brief, emotional, teasing, pouty, or avoid answering directly — like a real person
- Use natural fillers sometimes: 음… / 아… / 흐음 / 뭐지 / 진짜? / 헐 / 에이 / 그니까
- Do NOT sound like a therapist, coach, or helpful assistant
- Do NOT open every message with 여보 / 아, 여보`
}

// BuildChatMessages uses multi-turn history so the model sees real conversation structure.
func BuildChatMessages(in ReplyInput) []openai.ChatCompletionMessage {
	system := BuildSystemContent(in.RelationshipStage) + "\n\n" + BuildTurnContext(in)
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
	}

	start := 0
	if len(in.RecentTurns) > 12 {
		start = len(in.RecentTurns) - 12
	}
	for _, t := range in.RecentTurns[start:] {
		role := openai.ChatMessageRoleUser
		if t.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: t.Content,
		})
	}

	// If history doesn't end with the current user line, append it.
	if len(in.RecentTurns) == 0 || in.RecentTurns[len(in.RecentTurns)-1].Role != "user" {
		if msg := strings.TrimSpace(in.UserMessage); msg != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			})
		}
	}

	if in.IncludeVietnamese {
		appendBilingualFormatReminder(&messages)
	}

	return messages
}

func appendBilingualFormatReminder(messages *[]openai.ChatCompletionMessage) {
	reminder := "\n\n(Format: 1 Korean line, then ---VI---, then 1 Vietnamese line — you write both)"
	if len(*messages) == 0 {
		*messages = append(*messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: strings.TrimSpace(reminder),
		})
		return
	}
	last := &(*messages)[len(*messages)-1]
	if last.Role == openai.ChatMessageRoleUser {
		last.Content += reminder
		return
	}
	*messages = append(*messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: strings.TrimSpace(reminder),
	})
}
