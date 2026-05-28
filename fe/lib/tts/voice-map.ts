/** Soniox TTS voice names (server uses these in voice_profiles.voice_id). */

export const SONIOX_VOICES = {
  femaleSoft: "Mina",
  femalePlayful: "Nina",
  femaleBright: "Claire",
  femaleWarm: "Emma",
  maleCalm: "Kenji",
  maleMature: "Daniel",
  maleEnergetic: "Noah",
} as const;

/** voice_profile_id → Soniox voice passed to WebSocket. */
export const VOICE_PROFILE_TO_SONIOX: Record<string, string> = {
  soft_female_01: SONIOX_VOICES.femaleSoft,
  cute_female_02: SONIOX_VOICES.femalePlayful,
  bright_female_03: SONIOX_VOICES.femaleBright,
  deep_male_01: SONIOX_VOICES.maleCalm,
  calm_male_02: SONIOX_VOICES.maleMature,
  energetic_male_03: SONIOX_VOICES.maleEnergetic,
};

export const PRESET_TO_SONIOX: Record<string, string> = {
  hani: SONIOX_VOICES.femaleSoft,
  mina: SONIOX_VOICES.femalePlayful,
  joon: SONIOX_VOICES.maleCalm,
};

export function sonioxVoiceForProfile(voiceProfileId: string): string {
  return VOICE_PROFILE_TO_SONIOX[voiceProfileId] ?? SONIOX_VOICES.femaleSoft;
}

export function sonioxVoiceForPreset(presetSlug: string): string {
  return PRESET_TO_SONIOX[presetSlug] ?? SONIOX_VOICES.femaleSoft;
}
