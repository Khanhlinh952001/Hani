package ai

import (
	"testing"
	"time"
)

func TestBootstrapLifeHasActivity(t *testing.T) {
	now := time.Date(2026, 5, 23, 1, 0, 0, 0, koreaTZ)
	life := BootstrapLife(now, DefaultEmotionalState(), 3, 42)
	if life.CurrentActivity == "" {
		t.Fatal("expected activity")
	}
	if life.Location == "" {
		t.Fatal("expected location")
	}
	if !life.ThinkingAboutUser {
		t.Fatal("expected thinking about user after 3h gap")
	}
}

func TestDeriveMoodLonelinessAfterGap(t *testing.T) {
	life := LifeState{Energy: "medium", Mood: "okay"}
	mood := DeriveMood(DefaultEmotionalState(), life, 8)
	if mood.Loneliness < 55 {
		t.Fatalf("expected high loneliness after 8h, got %d", mood.Loneliness)
	}
}

func TestDecideProactive(t *testing.T) {
	now := time.Date(2026, 5, 23, 23, 30, 0, 0, koreaTZ)
	if DecideProactive(7, now, true) != ProactiveMissedYou {
		t.Fatal("expected missed-you proactive")
	}
	if DecideProactive(3, now, true) != ProactiveLateNight {
		t.Fatal("expected late-night proactive")
	}
	if DecideProactive(1, now, true) != ProactiveNone {
		t.Fatal("expected no proactive")
	}
}

func TestGenerateInnerThought(t *testing.T) {
	thought := GenerateInnerThought(ReplyInput{
		RelationshipStage: StageDeeplyAttached,
		Life:              LifeState{ThinkingAboutUser: true, CurrentActivity: "watching YouTube"},
		Mood:              Mood{Loneliness: 70},
		HoursSinceUser:    7,
	})
	if thought == "" {
		t.Fatal("expected inner thought")
	}
}
