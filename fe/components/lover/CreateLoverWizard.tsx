"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import {
  ChevronLeft,
  ChevronRight,
  Loader2,
  Shuffle,
  Volume2,
  Heart,
} from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import {
  createLoverProfile,
  fetchNameSuggestions,
  fetchPersonalities,
  fetchSpeakingStyles,
  fetchVoices,
  previewLoverVoice,
  playBase64Audio,
} from "@/lib/lover/api";
import type {
  PersonalityTemplate,
  SpeakingStyleOption,
  VoiceProfile,
} from "@/lib/lover/types";
import { QuickPickTab } from "./QuickPickTab";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { pickCompanionAvatar } from "@/lib/avatars/catalog";
import { syncCompanionVoiceFromProfile } from "@/lib/lover/sync-voice";
import { cn } from "@/lib/utils";

const STEPS = 6;

type Mode = "quick" | "custom";

export function CreateLoverWizard() {
  const router = useRouter();
  const { user, loading: authLoading, applyUser } = useAuth();
  const { setTtsVoice } = useSettings();

  const [mode, setMode] = useState<Mode>("custom");
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [previewing, setPreviewing] = useState(false);

  const [companionGender, setCompanionGender] = useState<"female" | "male">(
    "female"
  );
  const [personalities, setPersonalities] = useState<PersonalityTemplate[]>(
    []
  );
  const [speakingStyles, setSpeakingStyles] = useState<SpeakingStyleOption[]>(
    []
  );
  const [voices, setVoices] = useState<VoiceProfile[]>([]);
  const [nameSuggestions, setNameSuggestions] = useState<string[]>([]);

  const [personalityId, setPersonalityId] = useState("");
  const [styleTags, setStyleTags] = useState<string[]>([]);
  const [voiceId, setVoiceId] = useState("");
  const [displayName, setDisplayName] = useState("");

  useEffect(() => {
    if (authLoading) return;
    if (user?.ai_profile_id || user?.selected_character_id) {
      router.replace("/");
    }
  }, [authLoading, user, router]);

  useEffect(() => {
    if (mode !== "custom") return;
    setLoading(true);
    Promise.all([
      fetchPersonalities(),
      fetchSpeakingStyles(),
      fetchVoices(companionGender),
      fetchNameSuggestions(companionGender),
    ])
      .then(([p, s, v, n]) => {
        setPersonalities(p);
        setSpeakingStyles(s);
        setVoices(v);
        setNameSuggestions(n);
        if (!personalityId && p[0]) setPersonalityId(p[0].id);
        if (!voiceId && v[0]) setVoiceId(v[0].id);
      })
      .catch((e) =>
        setError(e instanceof Error ? e.message : "Không tải được dữ liệu")
      )
      .finally(() => setLoading(false));
  }, [mode, companionGender]);

  const toggleStyle = (id: string) => {
    setStyleTags((prev) => {
      if (prev.includes(id)) return prev.filter((x) => x !== id);
      if (prev.length >= 2) return [prev[1], id];
      return [...prev, id];
    });
  };

  const selectedPersonality = personalities.find((p) => p.id === personalityId);
  const selectedVoice = voices.find((v) => v.id === voiceId);
  const previewAvatar = pickCompanionAvatar(companionGender, personalityId);
  const introKO =
    selectedVoice?.preview_text_ko ??
    `안녕… 나 ${displayName || "…"}야. 앞으로 잘 부탁해 💕`;

  const previewVoice = useCallback(async () => {
    if (!voiceId) return;
    setPreviewing(true);
    try {
      const { audio, format } = await previewLoverVoice(voiceId, introKO);
      playBase64Audio(audio, format);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không phát được giọng");
    } finally {
      setPreviewing(false);
    }
  }, [voiceId, introKO]);

  const finish = useCallback(async () => {
    if (!displayName.trim() || !personalityId || !voiceId) return;
    setSubmitting(true);
    setError(null);
    try {
      const { user: fresh, profile } = await createLoverProfile({
        companion_gender: companionGender,
        personality_template_id: personalityId,
        speaking_style_tags: styleTags,
        voice_profile_id: voiceId,
        display_name: displayName.trim(),
        avatar_url: previewAvatar,
      });
      applyUser(fresh);
      syncCompanionVoiceFromProfile(profile);
      if (profile.tts_voice) setTtsVoice(profile.tts_voice);
      router.replace("/chat");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không tạo được");
      setSubmitting(false);
    }
  }, [
    displayName,
    personalityId,
    voiceId,
    companionGender,
    styleTags,
    applyUser,
    setTtsVoice,
    router,
  ]);

  useEffect(() => {
    if (step === 6 && voiceId) {
      void previewVoice();
    }
  }, [step]);

  if (authLoading) {
    return (
      <div className="create-lover-screen flex items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="create-lover-screen">
      <header className="create-lover-header">
        <p className="choose-hani-eyebrow">Create Your AI Lover</p>
        <h1 className="choose-hani-title">Tạo người yêu AI</h1>
        <div className="create-lover-mode-tabs">
          <button
            type="button"
            className={cn(mode === "quick" && "create-lover-mode-active")}
            onClick={() => setMode("quick")}
          >
            Gặp nhanh
          </button>
          <button
            type="button"
            className={cn(mode === "custom" && "create-lover-mode-active")}
            onClick={() => setMode("custom")}
          >
            Tạo riêng
          </button>
        </div>
      </header>

      {mode === "quick" ? (
        <QuickPickTab />
      ) : (
        <>
          <div className="create-lover-progress">
            {Array.from({ length: STEPS }, (_, i) => (
              <span
                key={i}
                className={cn(
                  "create-lover-dot",
                  i + 1 <= step && "create-lover-dot-active"
                )}
              />
            ))}
          </div>

          {loading && step > 1 ? (
            <div className="flex justify-center py-12">
              <Loader2 className="size-6 animate-spin text-primary" />
            </div>
          ) : (
            <div className="create-lover-step-body">
              {step === 1 && (
                <>
                  <h2 className="create-lover-step-title">
                    Chọn giới tính người đồng hành
                  </h2>
                  <p className="create-lover-step-sub">
                    Tự do chọn — không bị khóa theo giới tính của bạn
                  </p>
                  <div className="create-lover-gender-grid">
                    {(
                      [
                        { id: "female" as const, label: "Nữ", ko: "여성" },
                        { id: "male" as const, label: "Nam", ko: "남성" },
                      ] as const
                    ).map((g) => (
                      <button
                        key={g.id}
                        type="button"
                        className={cn(
                          "create-lover-gender-card",
                          companionGender === g.id &&
                            "create-lover-gender-card-active"
                        )}
                        onClick={() => setCompanionGender(g.id)}
                      >
                        <span className="text-2xl">
                          {g.id === "female" ? "💕" : "🌙"}
                        </span>
                        <span className="font-semibold">{g.label}</span>
                        <span className="text-xs text-muted-foreground">
                          {g.ko}
                        </span>
                      </button>
                    ))}
                  </div>
                </>
              )}

              {step === 2 && (
                <>
                  <h2 className="create-lover-step-title">Chọn personality</h2>
                  <div className="create-lover-personality-grid">
                    {personalities.map((p) => (
                      <button
                        key={p.id}
                        type="button"
                        className={cn(
                          "create-lover-personality-card",
                          personalityId === p.id &&
                            "create-lover-personality-card-active"
                        )}
                        onClick={() => setPersonalityId(p.id)}
                      >
                        <span className="text-xl">{p.icon}</span>
                        <span className="font-medium">{p.name_vi}</span>
                        <span className="text-[0.625rem] text-muted-foreground">
                          {p.description_ko}
                        </span>
                      </button>
                    ))}
                  </div>
                </>
              )}

              {step === 3 && (
                <>
                  <h2 className="create-lover-step-title">Cách nói chuyện</h2>
                  <p className="create-lover-step-sub">Chọn tối đa 2</p>
                  <div className="create-lover-style-grid">
                    {speakingStyles.map((s) => (
                      <button
                        key={s.id}
                        type="button"
                        className={cn(
                          "create-lover-style-chip",
                          styleTags.includes(s.id) &&
                            "create-lover-style-chip-active"
                        )}
                        onClick={() => toggleStyle(s.id)}
                      >
                        <span>{s.label_vi}</span>
                        <span className="text-[0.625rem] opacity-70">
                          {s.description}
                        </span>
                      </button>
                    ))}
                  </div>
                </>
              )}

              {step === 4 && (
                <>
                  <h2 className="create-lover-step-title">Chọn giọng nói</h2>
                  <div className="create-lover-voice-list">
                    {voices.map((v) => (
                      <button
                        key={v.id}
                        type="button"
                        className={cn(
                          "create-lover-voice-row",
                          voiceId === v.id && "create-lover-voice-row-active"
                        )}
                        onClick={() => setVoiceId(v.id)}
                      >
                        <span className="font-medium">{v.name_vi}</span>
                        <span className="text-xs text-muted-foreground">
                          {v.emotion} · {v.speed}
                        </span>
                      </button>
                    ))}
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="mt-3 shadow-md shadow-primary/20 h-10 w-full"
                    disabled={!voiceId || previewing}
                    onClick={() => void previewVoice()}
                  >
                    {previewing ? (
                      <Loader2 className="size-4 animate-spin" />
                    ) : (
                      <Volume2 className="size-4" />
                    )}
                    Nghe thử
                  </Button>
                </>
              )}

              {step === 5 && (
                <>
                  <h2 className="create-lover-step-title">Đặt tên</h2>
                  <Input
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    placeholder="Hani, Mina, Joon…"
                    className="mt-2 h-12 text-center text-lg"
                    maxLength={24}
                  />
                  <div className="mt-3 flex flex-wrap justify-center gap-2">
                    {nameSuggestions.slice(0, 8).map((n) => (
                      <button
                        key={n}
                        type="button"
                        className="create-lover-name-chip"
                        onClick={() => setDisplayName(n)}
                      >
                        {n}
                      </button>
                    ))}
                    <button
                      type="button"
                      className="create-lover-name-chip"
                      onClick={() => {
                        const i = Math.floor(
                          Math.random() * nameSuggestions.length
                        );
                        setDisplayName(nameSuggestions[i] ?? "Hani");
                      }}
                    >
                      <Shuffle className="size-3" />
                    </button>
                  </div>
                </>
              )}

              {step === 6 && (
                <>
                  <h2 className="create-lover-step-title">Gặp nhau nhé</h2>
                  <div className="create-lover-preview">
                    <div className="choose-card-avatar-wrap mx-auto">
                      <div className="choose-card-avatar-inner">
                        <Image
                          src={previewAvatar}
                          alt={displayName || "Companion"}
                          width={160}
                          height={160}
                          className="choose-card-avatar"
                        />
                      </div>
                      <span className="choose-card-online" />
                    </div>
                    <p className="mt-4 font-display text-2xl font-bold">
                      {displayName}
                    </p>
                    {selectedPersonality ? (
                      <p className="text-sm text-primary">
                        {selectedPersonality.icon}{" "}
                        {selectedPersonality.name_vi}
                      </p>
                    ) : null}
                    <p className="mt-3 text-center italic text-foreground">
                      &ldquo;{introKO}&rdquo;
                    </p>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="mt-2"
                      disabled={previewing}
                      onClick={() => void previewVoice()}
                    >
                      <Volume2 className="size-4" />
                      Nghe lại
                    </Button>
                  </div>
                </>
              )}
            </div>
          )}

          {error ? (
            <p className="px-6 text-center text-sm text-destructive">{error}</p>
          ) : null}

          <footer className="create-lover-footer">
            {step > 1 ? (
              <Button
                type="button"
                variant="ghost"
                onClick={() => setStep((s) => s - 1)}
              >
                <ChevronLeft className="size-4" />
                Quay lại
              </Button>
            ) : (
              <span />
            )}
            {step < STEPS ? (
              <Button
                type="button"
                onClick={() => setStep((s) => s + 1)}
                disabled={step === 5 && !displayName.trim()}
              >
                Tiếp
                <ChevronRight className="size-4" />
              </Button>
            ) : (
              <Button
                type="button"
                className="choose-btn-talk"
                disabled={submitting}
                onClick={() => void finish()}
              >
                {submitting ? (
                  <Loader2 className="size-4 animate-spin" />
                ) : (
                  <Heart className="size-4 fill-current" />
                )}
                Trò chuyện với {displayName}
              </Button>
            )}
          </footer>
        </>
      )}
    </div>
  );
}
