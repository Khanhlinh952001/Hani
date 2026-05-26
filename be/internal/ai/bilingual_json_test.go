package ai

import "testing"

func TestParseJSONBilingual(t *testing.T) {
	raw := `{"korean":"뭐해 지금?","vietnamese":"Bạn đang làm gì vậy?"}`
	got, err := ParseJSONBilingual(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.Korean != "뭐해 지금?" {
		t.Fatalf("korean = %q", got.Korean)
	}
	if got.Vietnamese != "Bạn đang làm gì vậy?" {
		t.Fatalf("vietnamese = %q", got.Vietnamese)
	}
}
