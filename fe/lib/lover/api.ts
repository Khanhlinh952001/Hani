import { API_URL } from "../config";
import { authHeaders } from "../auth/api";
import type { AuthUser } from "../auth/types";
import type {
  CreateProfilePayload,
  LoverProfile,
  PersonalityTemplate,
  SpeakingStyleOption,
  VoiceProfile,
} from "./types";

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

export async function fetchPersonalities(): Promise<PersonalityTemplate[]> {
  const res = await fetch(`${API_URL}/api/lover/personalities`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function fetchSpeakingStyles(): Promise<SpeakingStyleOption[]> {
  const res = await fetch(`${API_URL}/api/lover/speaking-styles`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function fetchVoices(companionGender: string): Promise<VoiceProfile[]> {
  const q = companionGender
    ? `?companion_gender=${encodeURIComponent(companionGender)}`
    : "";
  const res = await fetch(`${API_URL}/api/lover/voices${q}`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function fetchNameSuggestions(
  companionGender: string
): Promise<string[]> {
  const res = await fetch(
    `${API_URL}/api/lover/name-suggestions?companion_gender=${encodeURIComponent(companionGender)}`,
    { headers: authHeaders() }
  );
  if (!res.ok) throw new Error(await parseError(res));
  const data = (await res.json()) as { names: string[] };
  return data.names;
}

export async function previewLoverVoice(
  voiceProfileId: string,
  text?: string
): Promise<{ audio: string; format: string }> {
  const res = await fetch(`${API_URL}/api/lover/preview-voice`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify({ voice_profile_id: voiceProfileId, text }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function createLoverProfile(
  payload: CreateProfilePayload
): Promise<{ profile: LoverProfile; user: AuthUser }> {
  const res = await fetch(`${API_URL}/api/lover/profile`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function createQuickPreset(
  presetSlug: string
): Promise<{ profile: LoverProfile; user: AuthUser }> {
  const res = await fetch(`${API_URL}/api/lover/profile/quick`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify({ preset_slug: presetSlug }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function fetchMyLoverProfile(): Promise<LoverProfile | null> {
  const res = await fetch(`${API_URL}/api/lover/profile/me`, {
    headers: authHeaders(),
  });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export function playBase64Audio(b64: string, format: string) {
  const mime = format === "mp3" ? "audio/mpeg" : `audio/${format}`;
  const audio = new Audio(`data:${mime};base64,${b64}`);
  void audio.play();
}
