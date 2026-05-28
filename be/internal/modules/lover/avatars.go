package lover

// Companion avatar paths (fe/public).
var (
	FemaleAvatars = []string{
		"/gird.jpg",
		"/girld1.jpg",
		"/girld2.jpg",
		"/girld3.jpg",
		"/girld4.jpg",
	}
	MaleAvatars = []string{
		"/boy.jpg",
		"/boy1.jpg",
		"/boy2.jpg",
		"/boy3.jpg",
		"/boy4.jpg",
	}
)

// personalityAvatarSlot picks a stable portrait per personality archetype.
var personalityAvatarSlot = map[string]int{
	"cute_soft":      0,
	"mature_calm":    1,
	"playful_funny":  2,
	"clingy_caring":  3,
	"cold_sweet":     4,
	"energetic":      2,
	"protective":     1,
	"ceo_vibe":       0,
	"romantic":       3,
}

var presetAvatars = map[string]string{
	"hani": "/gird.jpg",
	"mina": "/girld2.jpg",
	"joon": "/boy.jpg",
}

// AvatarForCompanion returns avatar path for gender + personality (or first of pool).
func AvatarForCompanion(companionGender, personalityID string) string {
	pool := FemaleAvatars
	if companionGender == "male" {
		pool = MaleAvatars
	}
	if len(pool) == 0 {
		return "/gird.jpg"
	}
	slot := 0
	if personalityID != "" {
		if s, ok := personalityAvatarSlot[personalityID]; ok {
			slot = s
		}
	}
	if slot >= len(pool) {
		slot = slot % len(pool)
	}
	return pool[slot]
}

func DefaultAvatarForGender(companionGender string) string {
	return AvatarForCompanion(companionGender, "")
}
