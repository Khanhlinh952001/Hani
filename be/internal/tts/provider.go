package tts

import (
	"context"
)

// NormalizeProvider returns openai or soniox.
func NormalizeProvider(provider string) string {
	if provider == "soniox" {
		return "soniox"
	}
	return "openai"
}

// DefaultProvider reads TTS_PROVIDER from env (fallback soniox).
func DefaultProvider() string {
	return NormalizeProvider(getenv("TTS_PROVIDER", "soniox"))
}

// StreamSpeech uses env TTS_PROVIDER (legacy).
func StreamSpeech(ctx context.Context, text string, opts *Options, onChunk func(index int, b64 string) error) error {
	return StreamSpeechFor(ctx, DefaultProvider(), text, opts, onChunk)
}

// StreamSpeechFor synthesizes speech with the given provider.
func StreamSpeechFor(ctx context.Context, provider, text string, opts *Options, onChunk func(index int, b64 string) error) error {
	switch NormalizeProvider(provider) {
	case "soniox":
		return sonioxStreamSpeech(ctx, text, opts, onChunk)
	default:
		return openAIStreamSpeech(ctx, text, opts, onChunk)
	}
}

// AudioFormat returns format for env default provider.
func AudioFormat() string {
	return AudioFormatFor(DefaultProvider())
}

// AudioFormatFor returns the audio format label for a provider.
func AudioFormatFor(provider string) string {
	if NormalizeProvider(provider) == "soniox" {
		return sonioxAudioFormat()
	}
	return "mp3"
}
