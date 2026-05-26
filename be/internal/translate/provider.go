package translate

import (
	"context"
	"os"

	"be/internal/ai"
)

// ToVietnamese translates Korean chat text to Vietnamese (Soniox by default).
func ToVietnamese(ctx context.Context, korean string) (string, error) {
	if os.Getenv("TRANSLATE_PROVIDER") == "openai" {
		return ai.TranslateToVietnamese(ctx, korean)
	}
	return sonioxKoreanToVietnamese(ctx, korean)
}
