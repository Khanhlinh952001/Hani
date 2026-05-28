package characters

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"be/internal/tts"
)

func ListForUserService(userGender string) ([]PublicCharacter, error) {
	list, err := repoListAll()
	if err != nil {
		return nil, err
	}
	sorted := sortByPreference(list, userGender)
	out := make([]PublicCharacter, 0, len(sorted))
	for _, c := range sorted {
		out = append(out, c.ToPublic())
	}
	return out, nil
}

func sortByPreference(list []Character, userGender string) []Character {
	if len(list) <= 1 {
		return list
	}
	preferFemale := userGender == "male"
	preferMale := userGender == "female"

	var first, rest []Character
	for _, c := range list {
		if preferFemale && c.Gender == "female" {
			first = append(first, c)
		} else if preferMale && c.Gender == "male" {
			first = append(first, c)
		} else {
			rest = append(rest, c)
		}
	}
	return append(first, rest...)
}

func GetByIDService(id string) (*Character, error) {
	return repoGetByID(id)
}

func SelectForUserService(userID int, characterID string) (*Character, error) {
	characterID = strings.TrimSpace(characterID)
	if characterID == "" {
		return nil, errors.New("character_id is required")
	}
	c, err := repoGetByID(characterID)
	if err != nil {
		return nil, err
	}
	if err := repoUpsertUserMemory(userID, characterID); err != nil {
		return nil, err
	}
	return c, nil
}

// PreviewVoiceKO synthesizes intro line for character selection.
func PreviewVoiceKO(ctx context.Context, slug string) (audioB64 string, format string, err error) {
	c, err := repoGetByID(slug)
	if err != nil {
		return "", "", err
	}
	text := strings.TrimSpace(c.IntroMessageKO)
	if text == "" {
		return "", "", errors.New("no intro message")
	}
	provider := c.VoiceProvider
	if provider == "" {
		provider = "soniox"
	}
	opts := &tts.Options{Voice: c.VoiceID, Language: c.TTSLanguage}
	if opts.Voice == "" {
		opts.Voice = "Mina"
	}

	var chunks []byte
	err = tts.StreamSpeechFor(ctx, provider, text, opts, func(_ int, b64 string) error {
		raw, decErr := base64.StdEncoding.DecodeString(b64)
		if decErr != nil {
			return decErr
		}
		chunks = append(chunks, raw...)
		return nil
	})
	if err != nil {
		return "", "", err
	}
	format = tts.AudioFormatFor(provider)
	return base64.StdEncoding.EncodeToString(chunks), format, nil
}

func TouchInteraction(userID int, characterID string) {
	if userID == 0 || characterID == "" {
		return
	}
	repoTouchInteraction(userID, characterID)
}
