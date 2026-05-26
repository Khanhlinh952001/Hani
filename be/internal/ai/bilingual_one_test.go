package ai

import "testing"

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

func TestParseBilingual_oneSentenceEach(t *testing.T) {
	raw := `뭐해 지금?
---VI---
Bạn đang làm gì vậy? Lâu rồi không thấy.`

	got := ParseBilingual(raw)
	if got.Korean != "뭐해 지금?" {
		t.Fatalf("Korean = %q", got.Korean)
	}
	if got.Vietnamese != "Bạn đang làm gì vậy?" {
		t.Fatalf("Vietnamese = %q, want first sentence only", got.Vietnamese)
	}
}
