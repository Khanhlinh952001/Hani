package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const sonioxTempKeyURL = "https://api.soniox.com/v1/auth/temporary-api-key"

type tempKeyRequest struct {
	UsageType        string `json:"usage_type"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type tempKeyResponse struct {
	APIKey string `json:"api_key"`
}

func CreateTemporaryTranscribeKey(ctx context.Context) (string, error) {
	master := APIKey()
	if master == "" {
		return "", fmt.Errorf("SONIOX_API_KEY is not set")
	}

	body, err := json.Marshal(tempKeyRequest{
		UsageType:        "transcribe_websocket",
		ExpiresInSeconds: 300,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sonioxTempKeyURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+master)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("soniox temp key request: %w", err)
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("soniox temp key: HTTP %d: %s", res.StatusCode, string(raw))
	}

	var parsed tempKeyResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("soniox temp key decode: %w", err)
	}
	if parsed.APIKey == "" {
		return "", fmt.Errorf("soniox temp key: empty response")
	}
	return parsed.APIKey, nil
}

func TemporaryKeyErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if os.Getenv("SONIOX_API_KEY") == "" {
		return "Thiếu SONIOX_API_KEY trên server (be/.env)"
	}
	return "Không tạo được khóa Soniox tạm — kiểm tra SONIOX_API_KEY và mạng"
}
