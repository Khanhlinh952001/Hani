import { saveSettings, loadSettings } from "@/lib/settings/storage";
import type { LoverProfile } from "./types";

/** Apply saved companion voice from DB profile → local settings (once per session). */
export function syncCompanionVoiceFromProfile(profile: LoverProfile | null) {
  if (!profile?.tts_voice) return;
  const current = loadSettings();
  if (
    current.ttsVoice === profile.tts_voice &&
    current.ttsProvider === "soniox"
  ) {
    return;
  }
  saveSettings({
    ...current,
    ttsProvider: "soniox",
    ttsVoice: profile.tts_voice,
  });
}
