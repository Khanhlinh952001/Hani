export type TtsLanguage = "auto" | "ko" | "vi" | "en";

export type TtsProvider = "openai" | "soniox";

export type AppSettings = {
  showVietnamese: boolean;
  ttsProvider: TtsProvider;
  ttsVoice: string;
  ttsLanguage: TtsLanguage;
};

export const DEFAULT_SETTINGS: AppSettings = {
  showVietnamese: true,
  ttsProvider: "openai",
  ttsVoice: "nova",
  ttsLanguage: "ko",
};

export const TTS_PROVIDER_OPTIONS: {
  id: TtsProvider;
  label: string;
  desc: string;
}[] = [
  {
    id: "openai",
    label: "OpenAI TTS",
    desc: "Giọng nova, shimmer… — cần OPENAI_API_KEY trên server",
  },
  {
    id: "soniox",
    label: "Soniox TTS",
    desc: "Giọng Kenji, Mina… — cần SONIOX_API_KEY trên server",
  },
];

export const TTS_VOICE_OPTIONS = [
  { id: "nova", label: "Nova", desc: "Nữ, ấm — khuyên dùng" },
  { id: "shimmer", label: "Shimmer", desc: "Nữ, nhẹ nhàng" },
  { id: "alloy", label: "Alloy", desc: "Trung tính" },
  { id: "echo", label: "Echo", desc: "Nam" },
  { id: "fable", label: "Fable", desc: "Kể chuyện" },
  { id: "onyx", label: "Onyx", desc: "Nam, trầm" },
] as const;

export const SONIOX_VOICE_OPTIONS = [
  { id: "Kenji", label: "Kenji", desc: "Nam — tiếng Hàn (mặc định)" },
  { id: "Mina", label: "Mina", desc: "Nữ — tiếng Hàn / Việt" },
  { id: "Emma", label: "Emma", desc: "Nữ — đa ngôn ngữ" },
] as const;

export const TTS_LANGUAGE_OPTIONS: { id: TtsLanguage; label: string; desc: string }[] = [
  { id: "ko", label: "Tiếng Hàn", desc: "Luôn đọc như tiếng Hàn" },
  { id: "auto", label: "Tự nhận", desc: "Theo nội dung câu" },
  { id: "vi", label: "Tiếng Việt", desc: "Đọc như tiếng Việt" },
  { id: "en", label: "English", desc: "Đọc như tiếng Anh" },
];
