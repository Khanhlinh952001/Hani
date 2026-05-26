package ai

import "strings"

// SentenceBuffer extracts complete sentences from streaming Korean/English text.
type SentenceBuffer struct {
	carry strings.Builder
}

func (b *SentenceBuffer) Feed(delta string) []string {
	if delta == "" {
		return nil
	}
	b.carry.WriteString(delta)
	raw := b.carry.String()

	var out []string
	start := 0
	for start < len(raw) {
		end := findSentenceEnd(raw, start)
		if end < 0 {
			break
		}
		sent := strings.TrimSpace(raw[start : end+1])
		if runeLen(sent) >= 2 {
			out = append(out, sent)
		}
		start = end + 1
	}

	b.carry.Reset()
	if start < len(raw) {
		b.carry.WriteString(raw[start:])
	}

	// Flush long fragments without waiting for punctuation (Korean often streams without `.`).
	if runeLen(b.carry.String()) >= 48 {
		if tail := strings.TrimSpace(b.carry.String()); tail != "" {
			b.carry.Reset()
			out = append(out, tail)
		}
	}
	return out
}

func (b *SentenceBuffer) Flush() string {
	s := strings.TrimSpace(b.carry.String())
	b.carry.Reset()
	return s
}

func findSentenceEnd(s string, from int) int {
	best := -1
	for _, sep := range []string{".", "!", "?", "。", "\n", "…"} {
		pos := from
		for {
			i := strings.Index(s[pos:], sep)
			if i < 0 {
				break
			}
			idx := pos + i + len(sep) - 1
			if idx > best {
				best = idx
			}
			pos = idx + 1
		}
	}
	return best
}

func runeLen(s string) int {
	return len([]rune(s))
}
