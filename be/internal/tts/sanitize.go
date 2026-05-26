package tts

import (
	"strings"
	"unicode"
)

// SanitizeForSpeech removes emoji and symbols (e.g. ♡ 💕) so TTS does not read them aloud.
func SanitizeForSpeech(s string) string {
	s = strings.ReplaceAll(s, "---VI---", "")
	s = strings.ReplaceAll(s, "---WORDS---", "")
	s = strings.ReplaceAll(s, "«VI»", "")
	var b strings.Builder
	for _, r := range s {
		if skipForSpeech(r) {
			continue
		}
		if r == '«' || r == '»' {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// WorthSpeaking returns false for fragments too short or marker debris for TTS.
func WorthSpeaking(s string) bool {
	s = SanitizeForSpeech(s)
	n := len([]rune(s))
	return n >= 4
}

func skipForSpeech(r rune) bool {
	if unicode.IsSymbol(r) {
		return true
	}
	switch {
	case r >= 0x1F000 && r <= 0x1FFFF:
		return true
	case r >= 0x2600 && r <= 0x27BF:
		return true
	case r >= 0xFE00 && r <= 0xFE0F:
		return true
	}
	return false
}

func trimText(s string) string {
	s = SanitizeForSpeech(s)
	const max = 4096
	if len(s) <= max {
		return s
	}
	return s[:max]
}
