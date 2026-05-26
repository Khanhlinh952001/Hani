package translate

import (
	"context"
	"fmt"
	"strings"

	"be/internal/stt"
	"be/internal/tts"
)

// sonioxKoreanToVietnamese synthesizes Korean speech then runs Soniox STT+translation (ko→vi).
func sonioxKoreanToVietnamese(ctx context.Context, korean string) (string, error) {
	korean = strings.TrimSpace(korean)
	if korean == "" {
		return "", nil
	}

	mp3, err := tts.CollectSonioxSpeech(ctx, korean, &tts.Options{Language: "ko"})
	if err != nil {
		return "", fmt.Errorf("soniox tts for translate: %w", err)
	}
	if len(mp3) == 0 {
		return "", fmt.Errorf("soniox translate: empty audio")
	}

	return stt.TranslateAudio(ctx, mp3, "auto")
}
