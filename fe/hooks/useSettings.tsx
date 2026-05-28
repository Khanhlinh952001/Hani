"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { loadSettings, saveSettings } from "@/lib/settings/storage";
import type { AppSettings, TtsLanguage, TtsProvider } from "@/lib/settings/types";
import { SONIOX_VOICE_OPTIONS } from "@/lib/settings/types";

type SettingsContextValue = AppSettings & {
  setShowVietnamese: (value: boolean) => void;
  setTtsProvider: (provider: TtsProvider) => void;
  setTtsVoice: (voice: string) => void;
  setTtsLanguage: (lang: TtsLanguage) => void;
  ready: boolean;
};

const SettingsContext = createContext<SettingsContextValue | null>(null);

export function SettingsProvider({ children }: { children: React.ReactNode }) {
  const [settings, setSettings] = useState(loadSettings);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    setSettings(loadSettings());
    setReady(true);
  }, []);

  const update = useCallback((patch: Partial<AppSettings>) => {
    setSettings((prev) => {
      const next = { ...prev, ...patch };
      saveSettings(next);
      return next;
    });
  }, []);

  const setShowVietnamese = useCallback(
    (value: boolean) => update({ showVietnamese: value }),
    [update]
  );

  const setTtsProvider = useCallback((_provider: TtsProvider) => {
    setSettings((prev) => {
      const voice = SONIOX_VOICE_OPTIONS.some((o) => o.id === prev.ttsVoice)
        ? prev.ttsVoice
        : "Mina";
      const next = { ...prev, ttsProvider: "soniox" as const, ttsVoice: voice };
      saveSettings(next);
      return next;
    });
  }, []);

  const setTtsVoice = useCallback(
    (voice: string) => update({ ttsVoice: voice }),
    [update]
  );

  const setTtsLanguage = useCallback(
    (lang: TtsLanguage) => update({ ttsLanguage: lang }),
    [update]
  );

  const value = useMemo(
    () => ({
      ...settings,
      setShowVietnamese,
      setTtsProvider,
      setTtsVoice,
      setTtsLanguage,
      ready,
    }),
    [
      settings,
      setShowVietnamese,
      setTtsProvider,
      setTtsVoice,
      setTtsLanguage,
      ready,
    ]
  );

  return (
    <SettingsContext.Provider value={value}>{children}</SettingsContext.Provider>
  );
}

export function useSettings() {
  const ctx = useContext(SettingsContext);
  if (!ctx) {
    throw new Error("useSettings must be used within SettingsProvider");
  }
  return ctx;
}
