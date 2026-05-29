package ai

import (
	"strings"
	"testing"
)

func TestFirstSentence(t *testing.T) {
	cases := []struct{ in, want string }{
		{"뭐해? 오랜만이야", "뭐해?"},
		{"음… 그냥\n두 번째 줄", "음… 그냥"},
		{"한 줄만", "한 줄만"},
	}
	for _, c := range cases {
		if got := firstSentence(c.in); got != c.want {
			t.Fatalf("firstSentence(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseBilingual_keepsFullTranslation(t *testing.T) {
	raw := `뭐해 지금?
---VI---
Bạn đang làm gì vậy? Lâu rồi không thấy.`

	got := ParseBilingual(raw)
	if got.Korean != "뭐해 지금?" {
		t.Fatalf("Korean = %q", got.Korean)
	}
	wantVi := "Bạn đang làm gì vậy? Lâu rồi không thấy."
	if got.Vietnamese != wantVi {
		t.Fatalf("Vietnamese = %q, want %q", got.Vietnamese, wantVi)
	}
}

func TestParseBilingual_twoKoreanLines(t *testing.T) {
	raw := `나도 쉬는 시간이라서 잠깐 멍하니 창밖을 보고 있어...
오빠는 잘 하고 있겠지?
---VI---
Em cũng đang nghỉ nên nhìn ra ngoài một chút...
Anh vẫn ổn chứ?`

	got := ParseBilingual(raw)
	if !strings.Contains(got.Korean, "오빠는") {
		t.Fatalf("Korean should keep both lines: %q", got.Korean)
	}
	if !strings.Contains(got.Vietnamese, "Anh") && !strings.Contains(got.Vietnamese, "ổn") {
		t.Fatalf("Vietnamese should keep both lines: %q", got.Vietnamese)
	}
}
