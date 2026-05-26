package tts

// Options optional per-session overrides (from client settings).
type Options struct {
	Voice    string
	Language string // ko, vi, en, or auto
}
