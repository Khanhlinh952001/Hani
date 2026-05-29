export const ClientEvents = {
  StartListening: "start_listening",
  StopSpeaking: "stop_speaking",
  SessionEnd: "session_end",
  Ping: "ping",
} as const;

export const ServerEvents = {
  Ready: "ready",
  Listening: "listening",
  PartialTranscript: "partial_transcript",
  FinalTranscript: "final_transcript",
  TypingStart: "typing_start",
  TypingEnd: "typing_end",
  AIResponse: "ai_response",
  Subtitle: "subtitle",
  TranslationDelta: "translation_delta",
  AIAudioChunk: "ai_audio_chunk",
  AIAudioSegmentEnd: "ai_audio_segment_end",
  AIAudioEnd: "ai_audio_end",
  SessionEnded: "session_ended",
  Pong: "pong",
  Error: "error",
  QuotaExceeded: "quota_exceeded",
  QuotaWarning: "quota_warning",
} as const;

export type HistoryMessage = {
  id?: string;
  role: "user" | "assistant";
  content: string;
  translation?: string;
};

export type ServerMessage = {
  type: string;
  text?: string;
  translation?: string;
  stt_context?: string;
  messages?: HistoryMessage[];
  history_has_more?: boolean;
  message?: string;
  user_id?: number;
  session_id?: string;
  delta?: string;
  full_text?: string;
  audio?: string;
  format?: string;
  index?: number;
  finished?: boolean;
  code?: string;
};

export type ChatMessage = {
  id: string;
  role: "user" | "assistant";
  content: string;
  translationVi?: string;
  streaming?: boolean;
  /** Chat mode: play one-shot bubble pop when message appears */
  justArrived?: boolean;
};

export type ConnectionStatus =
  | "disconnected"
  | "connecting"
  | "ready"
  | "listening"
  | "thinking"
  | "speaking"
  | "error";
