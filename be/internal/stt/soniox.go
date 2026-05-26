package stt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const sonioxWSURL = "wss://stt-rt.soniox.com/transcribe-websocket"

type Config struct {
	APIKey            string
	Model             string
	LanguageHints     []string
	ContextText       string
	AudioFormat       string // pcm_s16le | auto
	TranslationTarget string // e.g. "vi" enables one_way translation
}

type Token struct {
	Text              string `json:"text"`
	IsFinal           bool   `json:"is_final"`
	TranslationStatus string `json:"translation_status"`
}

type Response struct {
	Tokens       []Token `json:"tokens"`
	Finished     bool    `json:"finished"`
	ErrorMessage string  `json:"error_message"`
	ErrorType    string  `json:"error_type"`
}

type SonioxClient struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func APIKey() string {
	return os.Getenv("SONIOX_API_KEY")
}

func (c *SonioxClient) Connect(ctx context.Context, cfg Config) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("SONIOX_API_KEY is not set")
	}
	if cfg.Model == "" {
		cfg.Model = "stt-rt-preview"
	}
	if len(cfg.LanguageHints) == 0 {
		cfg.LanguageHints = []string{"ko"}
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, sonioxWSURL, nil)
	if err != nil {
		return fmt.Errorf("soniox dial: %w", err)
	}
	c.conn = conn

	audioFormat := cfg.AudioFormat
	if audioFormat == "" {
		audioFormat = "pcm_s16le"
	}
	start := map[string]interface{}{
		"api_key":                   cfg.APIKey,
		"model":                     cfg.Model,
		"audio_format":              audioFormat,
		"language_hints":            cfg.LanguageHints,
		"enable_endpoint_detection": true,
	}
	if audioFormat == "pcm_s16le" {
		start["sample_rate"] = 16000
		start["num_channels"] = 1
	}
	if cfg.TranslationTarget != "" {
		start["translation"] = map[string]interface{}{
			"type":             "one_way",
			"target_language": cfg.TranslationTarget,
		}
	}
	if cfg.ContextText != "" {
		start["context"] = map[string]interface{}{"text": cfg.ContextText}
	}

	c.mu.Lock()
	err = c.conn.WriteJSON(start)
	c.mu.Unlock()
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("soniox config: %w", err)
	}
	return nil
}

func (c *SonioxClient) SendAudio(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("soniox not connected")
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *SonioxClient) Finish() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("soniox not connected")
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, []byte{})
}

func (c *SonioxClient) Read(ctx context.Context) (*Response, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("soniox not connected")
	}

	type result struct {
		resp *Response
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			ch <- result{err: err}
			return
		}
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			ch <- result{err: err}
			return
		}
		ch <- result{resp: &resp}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		if r.err != nil {
			return nil, r.err
		}
		if r.resp.ErrorMessage != "" {
			return nil, fmt.Errorf("soniox: %s (%s)", r.resp.ErrorMessage, r.resp.ErrorType)
		}
		return r.resp, nil
	}
}

func tokensByStatus(tokens []Token, status string, finalOnly bool) string {
	var b strings.Builder
	for _, t := range tokens {
		if finalOnly && !t.IsFinal {
			continue
		}
		switch status {
		case "translation":
			if t.TranslationStatus != "translation" {
				continue
			}
		case "original":
			if t.TranslationStatus != "original" && t.TranslationStatus != "none" && t.TranslationStatus != "" {
				continue
			}
		default:
			if t.TranslationStatus == "translation" {
				continue
			}
		}
		b.WriteString(t.Text)
	}
	return strings.TrimSpace(b.String())
}

func TranscriptFromTokens(tokens []Token, finalOnly bool) string {
	return tokensByStatus(tokens, "original", finalOnly)
}

func TranslationFromTokens(tokens []Token, finalOnly bool) string {
	out := tokensByStatus(tokens, "translation", finalOnly)
	if out != "" {
		return out
	}
	return tokensByStatus(tokens, "translation", false)
}

// TranslateAudio runs Soniox STT with one-way translation on recorded audio bytes.
func TranslateAudio(ctx context.Context, audio []byte, audioFormat string) (string, error) {
	key := APIKey()
	if key == "" {
		return "", fmt.Errorf("SONIOX_API_KEY is not set")
	}
	if audioFormat == "" {
		audioFormat = "auto"
	}

	client := &SonioxClient{}
	if err := client.Connect(ctx, Config{
		APIKey:            key,
		Model:             "stt-rt-preview",
		LanguageHints:     []string{"ko"},
		AudioFormat:       audioFormat,
		TranslationTarget: "vi",
	}); err != nil {
		return "", err
	}
	defer client.Close()

	const chunk = 32 * 1024
	for i := 0; i < len(audio); i += chunk {
		end := i + chunk
		if end > len(audio) {
			end = len(audio)
		}
		if err := client.SendAudio(audio[i:end]); err != nil {
			return "", err
		}
	}
	if err := client.Finish(); err != nil {
		return "", err
	}

	var tokens []Token
	for {
		resp, err := client.Read(ctx)
		if err != nil {
			return "", err
		}
		if len(resp.Tokens) > 0 {
			tokens = resp.Tokens
		}
		if resp.Finished {
			break
		}
	}

	vi := TranslationFromTokens(tokens, true)
	if vi == "" {
		return "", fmt.Errorf("soniox translate: no vietnamese tokens")
	}
	return vi, nil
}

func (c *SonioxClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}
