package ai

import "testing"

func TestParseBilingual_stripsNumberedFormatLeak(t *testing.T) {
	raw := `1) 여보~ 저녁 메뉴는 생각해봤어?
2)
---VI---
1) Em ơi, tối nay ăn gì chưa?
2)`

	got := ParseBilingual(raw)
	wantKo := "여보~ 저녁 메뉴는 생각해봤어?"
	if got.Korean != wantKo {
		t.Fatalf("Korean = %q, want %q", got.Korean, wantKo)
	}
	if got.Vietnamese != "Em ơi, tối nay ăn gì chưa?" {
		t.Fatalf("Vietnamese = %q", got.Vietnamese)
	}
}
