"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import { Volume2, Heart, Loader2 } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import { fetchCharacters } from "@/lib/characters/api";
import { previewCharacterVoice } from "@/lib/characters/api";
import type { PublicCharacter } from "@/lib/characters/types";
import { createQuickPreset, playBase64Audio } from "@/lib/lover/api";
import { syncCompanionVoiceFromProfile } from "@/lib/lover/sync-voice";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function QuickPickTab() {
  const router = useRouter();
  const { applyUser } = useAuth();
  const { setTtsVoice } = useSettings();
  const [characters, setCharacters] = useState<PublicCharacter[]>([]);
  const [activeIndex, setActiveIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [previewing, setPreviewing] = useState(false);
  const [selecting, setSelecting] = useState(false);
  const scrollerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchCharacters()
      .then(setCharacters)
      .catch((e) =>
        setError(e instanceof Error ? e.message : "Không tải được danh sách")
      )
      .finally(() => setLoading(false));
  }, []);

  const active = characters[activeIndex];

  const onScroll = useCallback(() => {
    const el = scrollerRef.current;
    if (!el || characters.length === 0) return;
    const card = el.querySelector<HTMLElement>(".choose-card");
    if (!card) return;
    const idx = Math.round(el.scrollLeft / (card.offsetWidth + 16));
    setActiveIndex(Math.min(Math.max(idx, 0), characters.length - 1));
  }, [characters.length]);

  const confirm = useCallback(async () => {
    if (!active || selecting) return;
    setSelecting(true);
    setError(null);
    try {
      const { user, profile } = await createQuickPreset(active.id);
      applyUser(user);
      syncCompanionVoiceFromProfile(profile);
      if (profile.tts_voice) setTtsVoice(profile.tts_voice);
      router.replace("/chat");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không lưu được");
      setSelecting(false);
    }
  }, [active, selecting, router, applyUser, setTtsVoice]);

  if (loading) {
    return (
      <div className="flex justify-center py-16">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="quick-pick-tab">
      <p className="mb-3 text-center text-sm text-muted-foreground">
        Gặp ngay Hani, Mina hoặc Joon — personality & giọng có sẵn
      </p>
      <div
        ref={scrollerRef}
        className="choose-hani-scroller"
        onScroll={onScroll}
      >
        {characters.map((c, i) => (
          <article
            key={c.id}
            className={cn(
              "choose-card",
              i === activeIndex && "choose-card-active"
            )}
          >
            <div className="choose-card-avatar-wrap">
              <div className="choose-card-avatar-inner">
                <Image
                  src={c.avatar_url}
                  alt={c.name}
                  width={200}
                  height={200}
                  className="choose-card-avatar"
                />
              </div>
              <span className="choose-card-online" />
            </div>
            <div className="choose-card-body">
              <p className="choose-card-name">
                {c.name}
                <span className="choose-card-name-ko">{c.display_name}</span>
              </p>
              <p className="choose-card-intro">&ldquo;{c.intro_message_ko}&rdquo;</p>
            </div>
          </article>
        ))}
      </div>
      {error ? (
        <p className="mt-2 text-center text-sm text-destructive">{error}</p>
      ) : null}
      <div className="mt-4 flex flex-col gap-2 px-4">
        <Button
          type="button"
          variant="outline"
          disabled={!active || previewing}
          onClick={async () => {
            if (!active) return;
            setPreviewing(true);
            try {
              const { audio, format } = await previewCharacterVoice(active.id);
              playBase64Audio(audio, format);
            } finally {
              setPreviewing(false);
            }
          }}
        >
          {previewing ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Volume2 className="size-4" />
          )}
          Nghe giọng
        </Button>
        <Button
          type="button"
          disabled={!active || selecting}
          onClick={() => void confirm()}
        >
          {selecting ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Heart className="size-4 fill-current" />
          )}
          {active ? `Gặp ${active.name}` : "Chọn"}
        </Button>
      </div>
    </div>
  );
}
