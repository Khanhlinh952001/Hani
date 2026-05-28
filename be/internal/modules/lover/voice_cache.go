package lover

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

const voicePreviewDir = "uploads/voices"

func voicePreviewFilePath(voiceProfileID string) string {
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, voiceProfileID)
	return filepath.Join(voicePreviewDir, safe+".mp3")
}

func readCachedPreview(voiceProfileID string) (b64 string, ok bool) {
	path := voicePreviewFilePath(voiceProfileID)
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return "", false
	}
	return base64.StdEncoding.EncodeToString(data), true
}

func writeCachedPreview(voiceProfileID string, raw []byte) error {
	if err := os.MkdirAll(voicePreviewDir, 0o755); err != nil {
		return err
	}
	path := voicePreviewFilePath(voiceProfileID)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return err
	}
	_ = dbUpdateVoicePreviewPath(voiceProfileID, path)
	return nil
}
