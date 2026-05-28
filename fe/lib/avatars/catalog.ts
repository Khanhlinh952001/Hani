/** Companion portraits in /public (female: gird + girld1–4, male: boy + boy1–4). */

export const FEMALE_AVATARS = [
  "/gird.jpg",
  "/girld1.jpg",
  "/girld2.jpg",
  "/girld3.jpg",
  "/girld4.jpg",
] as const;

export const MALE_AVATARS = [
  "/boy.jpg",
  "/boy1.jpg",
  "/boy2.jpg",
  "/boy3.jpg",
  "/boy4.jpg",
] as const;

const PERSONALITY_AVATAR_SLOT: Record<string, number> = {
  cute_soft: 0,
  mature_calm: 1,
  playful_funny: 2,
  clingy_caring: 3,
  cold_sweet: 4,
  energetic: 2,
  protective: 1,
  ceo_vibe: 0,
  romantic: 3,
};

export const PRESET_AVATARS: Record<string, string> = {
  hani: "/gird.jpg",
  mina: "/girld2.jpg",
  joon: "/boy.jpg",
};

export function avatarsForGender(gender: "female" | "male" | string) {
  return gender === "male" ? MALE_AVATARS : FEMALE_AVATARS;
}

/** Stable avatar for companion gender + personality template id. */
export function pickCompanionAvatar(
  companionGender: "female" | "male" | string,
  personalityId?: string
): string {
  const pool = avatarsForGender(companionGender);
  let slot = 0;
  if (personalityId && PERSONALITY_AVATAR_SLOT[personalityId] != null) {
    slot = PERSONALITY_AVATAR_SLOT[personalityId];
  }
  return pool[Math.min(slot, pool.length - 1)] ?? pool[0];
}
