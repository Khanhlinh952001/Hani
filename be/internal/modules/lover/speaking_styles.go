package lover

// Speaking style modifiers merged into composed prompt.
var speakingStyleModifiers = map[string]string{
	"emoji_heavy": `Use emojis naturally and often (💕 😳 ✨ ㅋㅋ) — warm, expressive texting.`,
	"natural":     `Text like a real Korean on KakaoTalk — casual, natural, not formal.`,
	"gentle":      `Always gentle 해요체; soft line breaks; never harsh.`,
	"short_msgs":  `Prefer short messages; sometimes one line only; avoid long paragraphs.`,
	"flirty":      `Light flirting when relationship allows — playful teasing, not explicit.`,
	"caring":      `Lead with care — ask how they feel, notice tiredness or stress.`,
	"teasing":     `Playful teasing with affection; never mean.`,
	"soft_korean": `Soft Korean romance vibe — warm, slightly poetic, still casual.`,
}

func SpeakingStyleOptions() []SpeakingStyleOption {
	return []SpeakingStyleOption{
		{ID: "emoji_heavy", LabelKO: "이모지 많이", LabelVI: "Nhiều emoji", Description: "💕 😳 ✨"},
		{ID: "natural", LabelKO: "자연스럽게", LabelVI: "Tự nhiên", Description: "카톡 느낌"},
		{ID: "gentle", LabelKO: "부드럽게", LabelVI: "Nhẹ nhàng", Description: "해요체"},
		{ID: "short_msgs", LabelKO: "짧게", LabelVI: "Tin ngắn", Description: "한두 줄"},
		{ID: "flirty", LabelKO: "살짝 설렘", LabelVI: "Hơi flirt", Description: "장난 + 설렘"},
		{ID: "caring", LabelKO: "다정하게", LabelVI: "Quan tâm", Description: "챙김"},
		{ID: "teasing", LabelKO: "장난기", LabelVI: "Trêu ghẹo", Description: "ㅋㅋ"},
		{ID: "soft_korean", LabelKO: "로맨틱", LabelVI: "Romance Hàn", Description: "감성"},
	}
}
