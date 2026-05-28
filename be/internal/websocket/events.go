package websocket

// Client → server (JSON text frames unless noted).
const (
	EventSessionStart    = "session_start"
	EventStartListening  = "start_listening"
	EventStopSpeaking    = "stop_speaking"
	EventSessionEnd   = "session_end"
	EventPing         = "ping"
)

// Server → client.
const (
	EventReady             = "ready"
	EventListening         = "listening"
	EventPartialTranscript = "partial_transcript"
	EventFinalTranscript   = "final_transcript"
	EventTypingStart       = "typing_start"
	EventTypingEnd         = "typing_end"
	EventAIResponse        = "ai_response"
	EventSubtitle          = "subtitle"
	EventTranslationDelta  = "translation_delta"
	EventAIAudioChunk       = "ai_audio_chunk"
	EventAIAudioSegmentEnd  = "ai_audio_segment_end"
	EventAIAudioEnd         = "ai_audio_end"
	EventSessionEnded      = "session_ended"
	EventPong              = "pong"
	EventError             = "error"
	EventQuotaExceeded     = "quota_exceeded"
	EventQuotaWarning      = "quota_warning"
)

// ClientMessage control payloads.
type ClientMessage struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"` // final transcript from client STT
	Translation string `json:"translation,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// HistoryMessage is a prior turn sent on connect.
type HistoryMessage struct {
	ID          string `json:"id,omitempty"`
	Role        string `json:"role"`
	Content     string `json:"content"`
	Translation string `json:"translation,omitempty"`
}

// ServerMessage outbound events.
type ServerMessage struct {
	Type       string           `json:"type"`
	Text        string           `json:"text,omitempty"`
	Translation string           `json:"translation,omitempty"`
	SttContext  string           `json:"stt_context,omitempty"`
	Messages   []HistoryMessage `json:"messages,omitempty"`
	Message    string           `json:"message,omitempty"`
	Code       string           `json:"code,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	Delta      string `json:"delta,omitempty"`
	FullText   string `json:"full_text,omitempty"`
	Audio      string `json:"audio,omitempty"` // base64 chunk
	Format     string `json:"format,omitempty"`
	Index      int    `json:"index,omitempty"`
	Finished   bool   `json:"finished,omitempty"`
}
