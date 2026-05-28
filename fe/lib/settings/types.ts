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
  ttsProvider: "soniox",
  ttsVoice: "Mina",
  ttsLanguage: "ko",
};

/** TTS is Soniox-only; OpenAI is used for chat on the server. */
export const TTS_PROVIDER_OPTIONS: {
  id: TtsProvider;
  label: string;
  desc: string;
}[] = [
  {
    id: "soniox",
    label: "Soniox TTS",
    desc: "Kenji, Mina, Emma — cần SONIOX_API_KEY trên server",
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
  { id: "Mina", label: "Mina", desc: "Nữ — nhẹ, ấm" },
  { id: "Nina", label: "Nina", desc: "Nữ — sáng, năng động" },
  { id: "Claire", label: "Claire", desc: "Nữ — rõ, tự tin" },
  { id: "Emma", label: "Emma", desc: "Nữ — tự nhiên, thân thiện" },
  { id: "Maya", label: "Maya", desc: "Nữ — ổn định, ấm" },
  { id: "Grace", label: "Grace", desc: "Nữ — nhẹ nhàng, êm" },
  { id: "Kenji", label: "Kenji", desc: "Nam — bình tĩnh" },
  { id: "Daniel", label: "Daniel", desc: "Nam — trưởng thành" },
  { id: "Noah", label: "Noah", desc: "Nam — trẻ, sôi nổi" },
] as const;

export const TTS_LANGUAGE_OPTIONS: { id: TtsLanguage; label: string; desc: string }[] = [
  { id: "ko", label: "Tiếng Hàn", desc: "Luôn đọc như tiếng Hàn" },
  { id: "auto", label: "Tự nhận", desc: "Theo nội dung câu" },
  { id: "vi", label: "Tiếng Việt", desc: "Đọc như tiếng Việt" },
  { id: "en", label: "English", desc: "Đọc như tiếng Anh" },
];
