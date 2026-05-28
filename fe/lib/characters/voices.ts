/** @deprecated Use sonioxVoiceForPreset from @/lib/tts/voice-map */
import { sonioxVoiceForPreset } from "@/lib/tts/voice-map";

export function voiceForCharacter(slug: string | undefined): string | undefined {
  if (!slug) return undefined;
  return sonioxVoiceForPreset(slug);
}
