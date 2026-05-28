import {
  AppSettings,
  DEFAULT_SETTINGS,
  SONIOX_VOICE_OPTIONS,
  type TtsLanguage,
  type TtsProvider,
} from "./types";

const SONIOX_VOICE_IDS = new Set(SONIOX_VOICE_OPTIONS.map((v) => v.id));

const KEY = "hani_settings";

function parseLanguage(v: unknown): TtsLanguage {
  if (v === "ko" || v === "vi" || v === "en" || v === "auto") return v;
  return DEFAULT_SETTINGS.ttsLanguage;
}

function parseProvider(_v: unknown): TtsProvider {
  return "soniox";
}

function parseVoice(voice: unknown, provider: TtsProvider): string {
  if (typeof voice !== "string") {
    return SONIOX_VOICE_OPTIONS[0].id;
  }
  if (SONIOX_VOICE_IDS.has(voice as (typeof SONIOX_VOICE_OPTIONS)[number]["id"])) {
    return voice;
  }
  // Migrate old OpenAI voice ids saved in localStorage.
  const openaiToSoniox: Record<string, string> = {
    nova: "Mina",
    shimmer: "Mina",
    alloy: "Emma",
    echo: "Kenji",
    fable: "Kenji",
    onyx: "Kenji",
  };
  if (openaiToSoniox[voice]) {
    return openaiToSoniox[voice];
  }
  return SONIOX_VOICE_OPTIONS[0].id;
}

export function loadSettings(): AppSettings {
  if (typeof window === "undefined") return DEFAULT_SETTINGS;
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return DEFAULT_SETTINGS;
    const parsed = JSON.parse(raw) as Partial<AppSettings>;
    const ttsProvider = parseProvider(parsed.ttsProvider);
    return {
      showVietnamese:
        typeof parsed.showVietnamese === "boolean"
          ? parsed.showVietnamese
          : DEFAULT_SETTINGS.showVietnamese,
      ttsProvider,
      ttsVoice: parseVoice(parsed.ttsVoice, ttsProvider),
      ttsLanguage: parseLanguage(parsed.ttsLanguage),
    };
  } catch {
    return DEFAULT_SETTINGS;
  }
}

export function saveSettings(settings: AppSettings) {
  localStorage.setItem(KEY, JSON.stringify(settings));
}
