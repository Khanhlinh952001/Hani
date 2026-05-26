package ai

import "testing"

func TestStripYeoboOpener(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"아, 여보~ 잡채도 맛있지!", "잡채도 맛있지!"},
		{"여보~ 뭐해?", "뭐해?"},
		{"잡채? 좋지~", "잡채? 좋지~"},
	}
	for _, c := range cases {
		got := StripYeoboOpener(c.in)
		if got != c.want {
			t.Fatalf("StripYeoboOpener(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRelationshipStageFromMessageCount(t *testing.T) {
	if RelationshipStageFromMessageCount(10) != StageAwkward {
		t.Fatal("expected awkward")
	}
	if RelationshipStageFromMessageCount(50) != StageClose {
		t.Fatal("expected close")
	}
	if RelationshipStageFromMessageCount(150) != StageRomantic {
		t.Fatal("expected romantic")
	}
	if RelationshipStageFromMessageCount(500) != StageDeeplyAttached {
		t.Fatal("expected deeply attached")
	}
}

func TestEvolveEmotionAfterExchange(t *testing.T) {
	state := DefaultEmotionalState()
	next := EvolveEmotionAfterExchange(state, "보고싶어", "나도")
	if next.MissingUser {
		t.Fatal("expected missingUser false after exchange")
	}
	if next.Attachment <= state.Attachment {
		t.Fatalf("attachment should rise, got %d", next.Attachment)
	}
}
