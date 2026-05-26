package tts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const sonioxTTSURL = "https://tts-rt.soniox.com/tts"

func sonioxAPIKey() string {
	return os.Getenv("SONIOX_API_KEY")
}

func sonioxModel() string {
	if m := os.Getenv("SONIOX_TTS_MODEL"); m != "" {
		return m
	}
	return "tts-rt-v1"
}

func sonioxAudioFormat() string {
	if f := os.Getenv("SONIOX_TTS_AUDIO_FORMAT"); f != "" {
		return f
	}
	return "mp3"
}

func sonioxSampleRate() int {
	if s := os.Getenv("SONIOX_TTS_SAMPLE_RATE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}
	return 24000
}

func sonioxBitrate() int {
	if s := os.Getenv("SONIOX_TTS_BITRATE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}
	return 192000
}

func languageForSpeech(text string, opts *Options) string {
	if opts != nil && opts.Language != "" && opts.Language != "auto" {
		return opts.Language
	}
	if override := os.Getenv("SONIOX_TTS_LANGUAGE"); override != "" && override != "auto" {
		return override
	}
	for _, r := range text {
		if unicode.Is(unicode.Hangul, r) {
			return "ko"
		}
	}
	lower := strings.ToLower(text)
	if strings.ContainsAny(lower, "àáạảãâầấậẩẫăằắặẳẵèéẹẻẽêềếệểễìíịỉĩòóọỏõôồốộổỗơờớợởỡùúụủũưừứựửữỳýỵỷỹđ") {
		return "vi"
	}
	return "ko"
}

var openAITTSVoiceNames = map[string]bool{
	"nova": true, "shimmer": true, "alloy": true, "echo": true, "fable": true, "onyx": true,
}

func voiceForLanguage(lang string, opts *Options) string {
	if opts != nil && opts.Voice != "" && !openAITTSVoiceNames[strings.ToLower(opts.Voice)] {
		return opts.Voice
	}
	if v := os.Getenv("SONIOX_TTS_VOICE"); v != "" {
		return v
	}
	switch lang {
	case "ko":
		return "Kenji"
	case "vi":
		return "Mina"
	default:
		return "Emma"
	}
}

func sonioxStreamSpeech(ctx context.Context, text string, opts *Options, onChunk func(index int, b64 string) error) error {
	key := sonioxAPIKey()
	if key == "" {
		return fmt.Errorf("SONIOX_API_KEY is not set")
	}
	text = trimText(text)
	if text == "" {
		return fmt.Errorf("empty tts text")
	}

	lang := languageForSpeech(text, opts)
	voice := voiceForLanguage(lang, opts)
	format := sonioxAudioFormat()

	body := map[string]interface{}{
		"model":        sonioxModel(),
		"language":     lang,
		"voice":        voice,
		"audio_format": format,
		"text":         text,
		"sample_rate":  sonioxSampleRate(),
	}
	if format == "mp3" || format == "aac" || format == "opus" {
		body["bitrate"] = sonioxBitrate()
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sonioxTTSURL, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			ErrorMessage string `json:"error_message"`
			ErrorType    string `json:"error_type"`
		}
		_ = json.Unmarshal(data, &apiErr)
		msg := apiErr.ErrorMessage
		if msg == "" {
			msg = string(data)
		}
		return fmt.Errorf("soniox tts: %s (%s)", strings.TrimSpace(msg), apiErr.ErrorType)
	}

	if len(data) == 0 {
		return fmt.Errorf("soniox tts: empty audio response")
	}

	idx := 0
	for offset := 0; offset < len(data); offset += chunkSize {
		end := offset + chunkSize
		if end > len(data) {
			end = len(data)
		}
		b64 := base64.StdEncoding.EncodeToString(data[offset:end])
		if err := onChunk(idx, b64); err != nil {
			return err
		}
		idx++
	}
	return nil
}

// CollectSonioxSpeech returns full MP3 audio via Soniox TTS (ignores TTS_PROVIDER).
func CollectSonioxSpeech(ctx context.Context, text string, opts *Options) ([]byte, error) {
	var buf bytes.Buffer
	err := sonioxStreamSpeech(ctx, text, opts, func(_ int, b64 string) error {
		raw, decErr := base64.StdEncoding.DecodeString(b64)
		if decErr != nil {
			return decErr
		}
		buf.Write(raw)
		return nil
	})
	return buf.Bytes(), err
}
