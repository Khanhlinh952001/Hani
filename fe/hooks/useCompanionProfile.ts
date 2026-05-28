"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { PRESET_AVATARS, pickCompanionAvatar } from "@/lib/avatars/catalog";
import { HANI_BRAND_LOGO } from "@/lib/brand/assets";
import { fetchMyLoverProfile } from "@/lib/lover/api";
import { syncCompanionVoiceFromProfile } from "@/lib/lover/sync-voice";
import type { LoverProfile } from "@/lib/lover/types";

export function useCompanionProfile() {
  const { user } = useAuth();
  const [profile, setProfile] = useState<LoverProfile | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!user?.ai_profile_id) {
      setProfile(null);
      return;
    }
    setLoading(true);
    fetchMyLoverProfile()
      .then((p) => {
        setProfile(p);
        syncCompanionVoiceFromProfile(p);
      })
      .catch(() => setProfile(null))
      .finally(() => setLoading(false));
  }, [user?.ai_profile_id]);

  const preset = user?.selected_character_id
    ? PRESET_AVATARS[user.selected_character_id]
    : undefined;

  const avatarUrl =
    profile?.avatar_url ||
    preset ||
    (user?.selected_character_id
      ? pickCompanionAvatar(
          user.selected_character_id === "joon" ? "male" : "female",
          user.selected_character_id === "hani"
            ? "cute_soft"
            : user.selected_character_id === "mina"
              ? "playful_funny"
              : "mature_calm"
        )
      : HANI_BRAND_LOGO);

  const displayName = profile?.display_name ?? "Hani";
  const ttsVoice = profile?.tts_voice;

  return { profile, loading, avatarUrl, displayName, ttsVoice };
}
