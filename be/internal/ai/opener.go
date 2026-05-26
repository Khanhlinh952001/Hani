package ai

import (
	"regexp"
	"strings"
)

var (
	yeoboLeadRe = regexp.MustCompile(`(?i)^(?:아|어|음)?[,，]?\s*여보[~\-~]?\s*`)
)

// HasPriorAssistant is true when Hani already spoke in this thread.
func HasPriorAssistant(turns []Turn) bool {
	for _, t := range turns {
		if t.Role == "assistant" {
			return true
		}
	}
	return false
}

func midConversationOpenerBan() string {
	return `MID-CONVERSATION (mandatory — there is chat history above):
- Do NOT start with: "아, 여보", "아, 여보~", "여보~", "여보," or any variant.
- Start by reacting to their LAST message (헐 / 응 / 그치 / ㅋㅋ) or answer directly.
- BAD: "아, 여보~ 잡채도 맛있지!"
- GOOD: "잡채? 좋지~ 내가 해줄게!"`
}

// StripYeoboOpener removes a leading 여보 pet-name greeting from Korean text.
func StripYeoboOpener(ko string) string {
	s := strings.TrimSpace(ko)
	if s == "" {
		return ko
	}
	for i := 0; i < 3; i++ {
		next := yeoboLeadRe.ReplaceAllString(s, "")
		next = strings.TrimSpace(next)
		if next == s {
			break
		}
		s = next
	}
	if s == "" {
		return ko
	}
	return s
}

func yeoboPrefixHold(koPart string) bool {
	koPart = strings.TrimSpace(koPart)
	if koPart == "" {
		return true
	}
	stripped := StripYeoboOpener(koPart)
	if stripped != koPart && stripped != "" {
		return false
	}
	holdPrefixes := []string{
		"아", "아,", "아, ", "아,  ", "아, 여", "아, 여보", "아, 여보~",
		"여보", "여보~", "여보,",
	}
	for _, p := range holdPrefixes {
		if strings.HasPrefix(koPart, p) && len(koPart) <= len(p)+1 {
			return true
		}
	}
	return false
}

func polishMidConversationKorean(ko string, stripOpener bool) string {
	if !stripOpener {
		return ko
	}
	return StripYeoboOpener(ko)
}
