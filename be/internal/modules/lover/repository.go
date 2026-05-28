package lover

import (
	"errors"

	"be/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func repoListPersonalities() ([]PersonalityTemplate, error) {
	var list []PersonalityTemplate
	err := db.DB.Order("sort_order asc, id asc").Find(&list).Error
	return list, err
}

func repoGetPersonality(id string) (*PersonalityTemplate, error) {
	var t PersonalityTemplate
	err := db.DB.First(&t, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("personality not found")
		}
		return nil, err
	}
	return &t, nil
}

func repoListVoices(gender string) ([]VoiceProfile, error) {
	var list []VoiceProfile
	q := db.DB.Order("sort_order asc, id asc")
	if gender == "female" || gender == "male" {
		q = q.Where("gender = ? OR gender = ''", gender)
	}
	err := q.Find(&list).Error
	return list, err
}

func repoGetVoice(id string) (*VoiceProfile, error) {
	var v VoiceProfile
	err := db.DB.First(&v, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("voice not found")
		}
		return nil, err
	}
	return &v, nil
}

func repoCreateProfile(p *AIProfile) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return db.DB.Create(p).Error
}

func repoSaveProfile(p *AIProfile) error {
	return db.DB.Save(p).Error
}

func repoGetProfileByUserID(userID int) (*AIProfile, error) {
	var p AIProfile
	err := db.DB.Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func repoGetProfileByID(id uuid.UUID) (*AIProfile, error) {
	var p AIProfile
	err := db.DB.First(&p, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &p, nil
}

func dbUpdateVoicePreviewPath(voiceProfileID, path string) error {
	return db.DB.Model(&VoiceProfile{}).Where("id = ?", voiceProfileID).
		Update("preview_audio_path", path).Error
}

func repoBackfillProfileTtsVoices() {
	var profiles []AIProfile
	if err := db.DB.Find(&profiles).Error; err != nil {
		return
	}
	for _, p := range profiles {
		v, err := repoGetVoice(p.VoiceProfileID)
		if err != nil {
			continue
		}
		if v.VoiceID == "" {
			continue
		}
		if p.TtsVoice == v.VoiceID {
			continue
		}
		_ = db.DB.Model(&p).Update("tts_voice", v.VoiceID).Error
	}
}

func repoUpsertRelationshipStats(userID int, profileID uuid.UUID) error {
	var row RelationshipStats
	err := db.DB.Where("user_id = ? AND ai_profile_id = ?", userID, profileID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = RelationshipStats{
			UserID:            userID,
			AIProfileID:       profileID,
			IntimacyLevel:     5,
			TrustLevel:        5,
			RelationshipStage: "stranger",
		}
		return db.DB.Create(&row).Error
	}
	return err
}

