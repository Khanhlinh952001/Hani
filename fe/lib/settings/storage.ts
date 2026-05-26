import {
  AppSettings,
  DEFAULT_SETTINGS,
  SONIOX_VOICE_OPTIONS,
  TTS_VOICE_OPTIONS,
  type TtsLanguage,
  type TtsProvider,
} from "./types";

const OPENAI_VOICE_IDS = new Set(TTS_VOICE_OPTIONS.map((v) => v.id));
const SONIOX_VOICE_IDS = new Set(SONIOX_VOICE_OPTIONS.map((v) => v.id));

const KEY = "hani_settings";

function parseLanguage(v: unknown): TtsLanguage {
  if (v === "ko" || v === "vi" || v === "en" || v === "auto") return v;
  return DEFAULT_SETTINGS.ttsLanguage;
}

function parseProvider(v: unknown): TtsProvider {
  if (v === "soniox" || v === "openai") return v;
  return DEFAULT_SETTINGS.ttsProvider;
}

function parseVoice(voice: unknown, provider: TtsProvider): string {
  if (typeof voice !== "string") {
    return provider === "soniox"
      ? SONIOX_VOICE_OPTIONS[0].id
      : DEFAULT_SETTINGS.ttsVoice;
  }
  if (provider === "soniox") {
    return SONIOX_VOICE_IDS.has(voice as (typeof SONIOX_VOICE_OPTIONS)[number]["id"])
      ? voice
      : SONIOX_VOICE_OPTIONS[0].id;
  }
  return OPENAI_VOICE_IDS.has(voice as (typeof TTS_VOICE_OPTIONS)[number]["id"])
    ? voice
    : DEFAULT_SETTINGS.ttsVoice;
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
