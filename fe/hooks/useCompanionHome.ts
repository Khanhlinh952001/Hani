"use client";

import { useEffect, useState } from "react";
import {
  clearUnread,
  getDailyPhrase,
  getDynamicSubtitle,
  getHaniMood,
  getProactiveLine,
  readLastPreview,
  readStreak,
  readUnread,
  readVoiceMinutesToday,
  touchStreak,
  type HomePreview,
  DEFAULT_CHAT_PREVIEW,
} from "@/lib/home/companion";

export function useCompanionHome() {
  const [subtitle, setSubtitle] = useState(getDynamicSubtitle);
  const [streak, setStreak] = useState(0);
  const [voiceMinutes, setVoiceMinutes] = useState(0);
  const [preview, setPreview] = useState<HomePreview | null>(null);
  const [unread, setUnread] = useState(false);
  const [daily] = useState(getDailyPhrase);
  const [mood] = useState(getHaniMood);
  const [proactive] = useState(getProactiveLine);

  useEffect(() => {
    setSubtitle(getDynamicSubtitle());
    setStreak(touchStreak());
    setVoiceMinutes(readVoiceMinutesToday());
    setPreview(readLastPreview());
    setUnread(readUnread());
  }, []);

  const chatPreview = preview?.text ?? DEFAULT_CHAT_PREVIEW;

  return {
    subtitle,
    streak,
    voiceMinutes,
    chatPreview,
    unread,
    daily,
    mood,
    proactive,
    clearUnread: () => {
      clearUnread();
      setUnread(false);
    },
  };
}
