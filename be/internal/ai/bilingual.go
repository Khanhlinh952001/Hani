package ai

import (
	"regexp"
	"strings"
)

// ViMarker separates Korean (spoken) from Vietnamese translation.
const ViMarker = "---VI---"

var (
	numberedLineRe   = regexp.MustCompile(`^\d+\)\s*(.*)$`)
	numberedOnlyRe   = regexp.MustCompile(`^\d+\)\s*$`)
	secondSentenceRe = regexp.MustCompile(`[.!?]\s+\S`)
)

// BilingualReply is Hani's Korean line plus optional Vietnamese translation.
type BilingualReply struct {
	Korean     string
	Vietnamese string
}

// ParseBilingual splits model output into Korean and Vietnamese.
func ParseBilingual(full string) BilingualReply {
	full = strings.TrimSpace(full)
	if full == "" {
		return BilingualReply{}
	}

	idx := strings.Index(full, ViMarker)
	if idx < 0 {
		return BilingualReply{Korean: cleanDisplayText(full)}
	}

	ko := cleanDisplayText(strings.TrimSpace(full[:idx]))
	vi := cleanDisplayText(strings.TrimSpace(full[idx+len(ViMarker):]))

	return BilingualReply{
		Korean:     ko,
		Vietnamese: vi,
	}
}

// firstSentence keeps one chat line (first line only; splits a second sentence on ". ? !" + space).
func firstSentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if idx := strings.Index(s, "\n"); idx >= 0 {
		s = strings.TrimSpace(s[:idx])
	}
	if loc := secondSentenceRe.FindStringIndex(s); loc != nil {
		s = strings.TrimSpace(s[:loc[0]+1])
	}
	return s
}

// cleanDisplayText removes format instructions the model sometimes copies into output.
func cleanDisplayText(s string) string {
	s = strings.ReplaceAll(s, ViMarker, "")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var parts []string
	seen := make(map[string]struct{})
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || numberedOnlyRe.MatchString(line) {
			continue
		}
		if m := numberedLineRe.FindStringSubmatch(line); m != nil {
			line = strings.TrimSpace(m[1])
		}
		if line == "" || isFormatMetaLine(line) {
			continue
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		parts = append(parts, line)
	}

	if len(parts) == 0 {
		return stripNumberedPrefix(s)
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts, "\n")
}

// SentenceCount estimates how many sentences a line contains (for translation completeness).
func SentenceCount(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n := 0
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ends := 0
		for _, r := range line {
			switch r {
			case '.', '?', '!', '…':
				ends++
			}
		}
		if ends == 0 {
			n++
		} else {
			n += ends
		}
	}
	if n == 0 {
		return 1
	}
	return n
}

func isFormatMetaLine(line string) bool {
	lower := strings.ToLower(line)
	switch {
	case strings.HasPrefix(line, "---"),
		strings.Contains(lower, "full vietnamese"),
		strings.Contains(lower, "spoken aloud"),
		strings.Contains(lower, "one line exactly"):
		return true
	}
	return false
}

func stripNumberedPrefix(s string) string {
	s = strings.TrimSpace(s)
	for i := 0; i < 5; i++ {
		if m := numberedLineRe.FindStringSubmatch(s); m != nil && m[1] != "" {
			s = strings.TrimSpace(m[1])
			continue
		}
		if idx := strings.Index(s, ")"); idx >= 0 && idx <= 3 {
			rest := strings.TrimSpace(s[idx+1:])
			if len(rest) > 0 && rest[0] != ')' {
				s = rest
				continue
			}
		}
		break
	}
	return strings.TrimSpace(s)
}

func holdbackSuffix(full string) int {
	maxHold := 0
	for hold := len(ViMarker) - 1; hold > 0; hold-- {
		if strings.HasSuffix(full, ViMarker[:hold]) && hold > maxHold {
			maxHold = hold
		}
	}
	return len(full) - maxHold
}

// safeKoreanEnd avoids sending a partial marker suffix to TTS/UI stream.
func safeKoreanEnd(full string, markerAt int) int {
	if markerAt >= 0 {
		return markerAt
	}
	return holdbackSuffix(full)
}

// streamKoreanOnly forwards Korean deltas only until ViMarker appears.
func streamKoreanOnly(onDelta func(string) error, delta string, acc *strings.Builder, koSent *int, stripOpener bool) error {
	if onDelta == nil {
		return nil
	}
	acc.WriteString(delta)
	full := acc.String()
	markerAt := strings.Index(full, ViMarker)
	koPart := full
	if markerAt >= 0 {
		koPart = full[:markerAt]
	}

	if stripOpener && *koSent == 0 && yeoboPrefixHold(koPart) {
		return nil
	}

	koStream := koPart
	if stripOpener {
		koStream = StripYeoboOpener(koPart)
	}
	end := safeKoreanEnd(koStream, -1)
	if *koSent < end {
		chunk := stripStreamNoise(koStream[*koSent:end])
		*koSent = end
		if chunk != "" {
			return onDelta(chunk)
		}
	}
	return nil
}

// stripStreamNoise removes numbered-step prefixes the model sometimes echoes while streaming.
func stripStreamNoise(s string) string {
	s = strings.TrimPrefix(s, "1) ")
	s = strings.TrimPrefix(s, "2) ")
	s = strings.TrimPrefix(s, "3) ")
	s = strings.TrimPrefix(s, "4) ")
	return s
}

func bilingualFormatInstruction(includeVi bool) string {
	if !includeVi {
		return "Reply in Korean only as Hani. No numbering, no format labels."
	}
	return `Write BOTH languages yourself in one reply (no separate translation):

[Korean — 1–2 short sentences, Hani speaking]
---VI---
[Vietnamese — translate ALL Korean sentences above; same meaning, 1–2 lines]

Never skip ---VI---. Vietnamese must cover every Korean sentence (if 2 Korean lines, 2 Vietnamese lines or one block with both ideas).`
}
