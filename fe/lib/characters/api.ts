import { API_URL } from "../config";
import { authHeaders } from "../auth/api";
import type { AuthUser } from "../auth/types";
import type { PublicCharacter } from "./types";

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

export async function fetchCharacters(): Promise<PublicCharacter[]> {
  const res = await fetch(`${API_URL}/api/characters`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function selectCharacter(characterId: string): Promise<AuthUser> {
  const res = await fetch(`${API_URL}/api/characters/select`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify({ character_id: characterId }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function previewCharacterVoice(
  characterId: string
): Promise<{ audio: string; format: string }> {
  const res = await fetch(
    `${API_URL}/api/characters/${encodeURIComponent(characterId)}/preview-voice`,
    { headers: authHeaders() }
  );
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}
