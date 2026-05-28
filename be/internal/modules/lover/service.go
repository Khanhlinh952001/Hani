package lover

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"be/internal/modules/characters"
	"be/internal/modules/users"
	"be/internal/tts"
)

type CreateProfileInput struct {
	CompanionGender       string
	PersonalityTemplateID string
	SpeakingStyleTags     []string
	VoiceProfileID        string
	DisplayName           string
	AvatarURL             string
}

var presetVoiceMap = map[string]string{
	"hani": "soft_female_01",
	"mina": "cute_female_02",
	"joon": "deep_male_01",
}

var presetPersonalityMap = map[string]string{
	"hani": "cute_soft",
	"mina": "playful_funny",
	"joon": "mature_calm",
}

func ListPersonalitiesService() ([]PersonalityTemplate, error) {
	return repoListPersonalities()
}

func ListVoicesService(companionGender string) ([]VoiceProfile, error) {
	return repoListVoices(companionGender)
}

func SpeakingStylesService() []SpeakingStyleOption {
	return SpeakingStyleOptions()
}

func NameSuggestions(gender string) []string {
	female := []string{"Hani", "Mina", "Jiyoon", "Yuna", "Haru", "Sora", "Eunbi"}
	male := []string{"Joon", "Minho", "Haru", "Jun", "Taeyang", "Hyun", "Sung"}
	switch gender {
	case "male":
		return male
	case "female":
		return female
	default:
		out := make([]string, 0, len(female)+len(male))
		out = append(out, female...)
		out = append(out, male...)
		return out
	}
}

func GetProfileForUserService(userID int) (*PublicProfile, error) {
	p, err := repoGetProfileByUserID(userID)
	if err != nil || p == nil {
		return nil, err
	}
	return enrichPublic(p)
}

func GetProfilePromptForUser(userID int) (displayName, prompt string, voice VoiceProfile, ok bool) {
	p, err := repoGetProfileByUserID(userID)
	if err != nil || p == nil {
		return "", "", VoiceProfile{}, false
	}
	v, _ := repoGetVoice(p.VoiceProfileID)
	if v == nil {
		v = &VoiceProfile{Provider: "soniox", VoiceID: "Mina", Language: "ko"}
	}
	if p.TtsVoice != "" {
		v.VoiceID = p.TtsVoice
	}
	return p.DisplayName, p.ComposedPrompt, *v, true
}

func resolveVoice(voiceProfileID string) (*VoiceProfile, error) {
	voiceProfileID = strings.TrimSpace(voiceProfileID)
	v, err := repoGetVoice(voiceProfileID)
	if err != nil {
		return nil, err
	}
	if v.VoiceID == "" {
		v.VoiceID = "Mina"
	}
	if v.Provider == "" {
		v.Provider = "soniox"
	}
	return v, nil
}

func applyCreateInput(p *AIProfile, in CreateProfileInput) error {
	personalityID := strings.TrimSpace(in.PersonalityTemplateID)
	tpl, err := repoGetPersonality(personalityID)
	if err != nil {
		return err
	}

	voiceID := strings.TrimSpace(in.VoiceProfileID)
	voice, err := resolveVoice(voiceID)
	if err != nil {
		return err
	}

	name := strings.TrimSpace(in.DisplayName)
	if name == "" {
		return errors.New("display_name is required")
	}

	gender := strings.TrimSpace(in.CompanionGender)
	if gender != "male" && gender != "female" {
		gender = "female"
	}

	styles := in.SpeakingStyleTags
	if len(styles) > 2 {
		styles = styles[:2]
	}

	avatar := strings.TrimSpace(in.AvatarURL)
	if avatar == "" {
		avatar = AvatarForCompanion(gender, personalityID)
	}

	p.DisplayName = name
	p.CompanionGender = gender
	p.PersonalityTemplateID = personalityID
	p.SpeakingStyleTags = styles
	p.VoiceProfileID = voiceID
	p.TtsVoice = voice.VoiceID
	p.AvatarURL = avatar
	p.ComposedPrompt = ComposePrompt(tpl, styles, gender, name)
	p.PresetSlug = ""
	return nil
}

func syncUserProfileLink(userID int, p *AIProfile) error {
	if err := users.SetAiProfileService(userID, p.ID.String(), p.CompanionGender); err != nil {
		return err
	}
	if p.PresetSlug != "" {
		_ = users.SetSelectedCharacterService(userID, p.PresetSlug)
	}
	return nil
}

// SyncUserProfileLinkIfNeeded fixes users missing ai_profile_id when a profile row exists.
func SyncUserProfileLinkIfNeeded(userID int) (*PublicProfile, error) {
	p, err := repoGetProfileByUserID(userID)
	if err != nil || p == nil {
		return nil, err
	}
	u, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		return nil, err
	}
	needsLink := u.AiProfileID == nil || strings.TrimSpace(*u.AiProfileID) == ""
	if needsLink {
		if err := syncUserProfileLink(userID, p); err != nil {
			return nil, err
		}
	}
	return enrichPublic(p)
}

func CreateProfileService(userID int, in CreateProfileInput) (*PublicProfile, error) {
	existing, err := repoGetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if err := applyCreateInput(existing, in); err != nil {
			return nil, err
		}
		if err := repoSaveProfile(existing); err != nil {
			return nil, err
		}
		_ = repoUpsertRelationshipStats(userID, existing.ID)
		return enrichPublic(existing)
	}

	profile := &AIProfile{UserID: userID}
	if err := applyCreateInput(profile, in); err != nil {
		return nil, err
	}
	profile.IntroMessageKO = defaultIntroKO(profile.DisplayName)
	profile.IntroMessageVI = defaultIntroVI(profile.DisplayName)

	if err := repoCreateProfile(profile); err != nil {
		return nil, err
	}
	_ = repoUpsertRelationshipStats(userID, profile.ID)
	return enrichPublic(profile)
}

func buildPresetProfile(userID int, presetSlug string) (*AIProfile, error) {
	presetSlug = strings.TrimSpace(presetSlug)
	ch, err := characters.GetByIDService(presetSlug)
	if err != nil {
		return nil, err
	}

	voiceID := presetVoiceMap[presetSlug]
	if voiceID == "" {
		voiceID = "soft_female_01"
	}
	personalityID := presetPersonalityMap[presetSlug]
	if personalityID == "" {
		personalityID = "cute_soft"
	}

	avatar := ch.AvatarURL
	if avatar == "" {
		avatar = presetAvatars[presetSlug]
	}

	prompt := ch.PersonalityPrompt
	if prompt == "" {
		if tpl, e := repoGetPersonality(personalityID); e == nil {
			prompt = ComposePrompt(tpl, nil, ch.Gender, ch.Name)
		}
	}

	introKO := ch.IntroMessageKO
	if introKO == "" {
		introKO = defaultIntroKO(ch.Name)
	}
	introVI := ch.IntroMessageVI
	if introVI == "" {
		introVI = defaultIntroVI(ch.Name)
	}

	voice, err := resolveVoice(voiceID)
	if err != nil {
		return nil, err
	}
	if ch.VoiceID != "" {
		voice.VoiceID = ch.VoiceID
	}

	return &AIProfile{
		UserID:                userID,
		DisplayName:           ch.Name,
		CompanionGender:       ch.Gender,
		PersonalityTemplateID: personalityID,
		SpeakingStyleTags:     nil,
		VoiceProfileID:        voiceID,
		TtsVoice:              voice.VoiceID,
		AvatarURL:             avatar,
		ComposedPrompt:        prompt,
		IntroMessageKO:        introKO,
		IntroMessageVI:        introVI,
		PresetSlug:            presetSlug,
	}, nil
}

func applyPresetToProfile(existing *AIProfile, presetSlug string) error {
	built, err := buildPresetProfile(existing.UserID, presetSlug)
	if err != nil {
		return err
	}
	existing.DisplayName = built.DisplayName
	existing.CompanionGender = built.CompanionGender
	existing.PersonalityTemplateID = built.PersonalityTemplateID
	existing.SpeakingStyleTags = built.SpeakingStyleTags
	existing.VoiceProfileID = built.VoiceProfileID
	existing.TtsVoice = built.TtsVoice
	existing.AvatarURL = built.AvatarURL
	existing.ComposedPrompt = built.ComposedPrompt
	existing.IntroMessageKO = built.IntroMessageKO
	existing.IntroMessageVI = built.IntroMessageVI
	existing.PresetSlug = built.PresetSlug
	return nil
}

func CreateFromPresetService(userID int, presetSlug string) (*PublicProfile, error) {
	existing, err := repoGetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if err := applyPresetToProfile(existing, presetSlug); err != nil {
			return nil, err
		}
		if err := repoSaveProfile(existing); err != nil {
			return nil, err
		}
		_ = repoUpsertRelationshipStats(userID, existing.ID)
		return enrichPublic(existing)
	}

	profile, err := buildPresetProfile(userID, presetSlug)
	if err != nil {
		return nil, err
	}
	if err := repoCreateProfile(profile); err != nil {
		return nil, err
	}
	_ = repoUpsertRelationshipStats(userID, profile.ID)
	return enrichPublic(profile)
}

func PreviewVoiceService(ctx context.Context, voiceProfileID, text string) (audioB64, format string, err error) {
	v, err := resolveVoice(voiceProfileID)
	if err != nil {
		return "", "", err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		text = v.PreviewTextKO
	}
	if text == "" {
		return "", "", errors.New("no preview text")
	}

	// Reuse cached intro MP3 for this voice profile (no Soniox call).
	if text == v.PreviewTextKO {
		if b64, ok := readCachedPreview(voiceProfileID); ok {
			return b64, tts.AudioFormatFor(v.Provider), nil
		}
	}

	provider := v.Provider
	opts := &tts.Options{Voice: v.VoiceID, Language: v.Language}

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
	if text == v.PreviewTextKO && len(chunks) > 0 {
		_ = writeCachedPreview(voiceProfileID, chunks)
	}
	return base64.StdEncoding.EncodeToString(chunks), tts.AudioFormatFor(provider), nil
}

func enrichPublic(p *AIProfile) (*PublicProfile, error) {
	tpl, _ := repoGetPersonality(p.PersonalityTemplateID)
	voice, _ := repoGetVoice(p.VoiceProfileID)
	pub := p.ToPublic(tpl, voice, p.SpeakingStyleTags)
	return &pub, nil
}
