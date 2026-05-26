package tts

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"be/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

const chunkSize = 4096

func APIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func voiceFromName(name string) openai.SpeechVoice {
	switch name {
	case "alloy":
		return openai.VoiceAlloy
	case "echo":
		return openai.VoiceEcho
	case "fable":
		return openai.VoiceFable
	case "onyx":
		return openai.VoiceOnyx
	case "shimmer":
		return openai.VoiceShimmer
	case "nova":
		return openai.VoiceNova
	default:
		return openai.VoiceNova
	}
}

func Voice() openai.SpeechVoice {
	return voiceFromName(config.GetEnv("OPENAI_TTS_VOICE", "nova"))
}

func voiceForRequest(opts *Options) openai.SpeechVoice {
	if opts != nil && opts.Voice != "" {
		return voiceFromName(opts.Voice)
	}
	return Voice()
}

func Model() openai.SpeechModel {
	m := config.GetEnv("OPENAI_TTS_MODEL", "tts-1")
	if m == "tts-1-hd" {
		return openai.TTSModel1HD
	}
	return openai.TTSModel1
}

func openAIStreamSpeech(ctx context.Context, text string, opts *Options, onChunk func(index int, b64 string) error) error {
	key := APIKey()
	if key == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set")
	}
	text = trimText(text)
	if text == "" {
		return fmt.Errorf("empty tts text")
	}

	client := openai.NewClient(key)
	resp, err := client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          Model(),
		Input:          text,
		Voice:          voiceForRequest(opts),
		ResponseFormat: openai.SpeechResponseFormatMp3,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	buf := make([]byte, chunkSize)
	idx := 0
	for {
		n, err := resp.Read(buf)
		if n > 0 {
			b64 := base64.StdEncoding.EncodeToString(buf[:n])
			if err := onChunk(idx, b64); err != nil {
				return err
			}
			idx++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

