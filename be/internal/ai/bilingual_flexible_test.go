package ai

import "testing"

func TestParseBilingualFlexible_viMarker(t *testing.T) {
	raw := "뭐해?\n---VI---\nBạn đang làm gì?"
	got, err := ParseBilingualFlexible(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.Korean != "뭐해?" || got.Vietnamese != "Bạn đang làm gì?" {
		t.Fatalf("got %#v", got)
	}
}

func TestParseBilingualFlexible_jsonFallback(t *testing.T) {
	raw := "```json\n{\"korean\":\"안녕\",\"vietnamese\":\"Chào bạn\"}\n```"
	got, err := ParseBilingualFlexible(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.Korean != "안녕" || got.Vietnamese != "Chào bạn" {
		t.Fatalf("got %#v", got)
	}
}
