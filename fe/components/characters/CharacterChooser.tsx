"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import { Volume2, Heart, Loader2 } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import {
  fetchCharacters,
  previewCharacterVoice,
  selectCharacter,
} from "@/lib/characters/api";
import type { PublicCharacter } from "@/lib/characters/types";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function playBase64Audio(b64: string, format: string) {
  const mime = format === "mp3" ? "audio/mpeg" : `audio/${format}`;
  const audio = new Audio(`data:${mime};base64,${b64}`);
  void audio.play();
}

export function CharacterChooser() {
  const router = useRouter();
  const { user, loading: authLoading, applyUser } = useAuth();
  const { setTtsVoice } = useSettings();

  const [characters, setCharacters] = useState<PublicCharacter[]>([]);
  const [activeIndex, setActiveIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [previewing, setPreviewing] = useState(false);
  const [selecting, setSelecting] = useState(false);
  const scrollerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (authLoading) return;
    if (user?.selected_character_id) {
      router.replace("/");
      return;
    }
    fetchCharacters()
      .then(setCharacters)
      .catch((e) =>
        setError(e instanceof Error ? e.message : "Không tải được danh sách")
      )
      .finally(() => setLoading(false));
  }, [authLoading, user?.selected_character_id, router]);

  const active = characters[activeIndex];

  const onScroll = useCallback(() => {
    const el = scrollerRef.current;
    if (!el || characters.length === 0) return;
    const card = el.querySelector<HTMLElement>(".choose-card");
    if (!card) return;
    const gap = 16;
    const w = card.offsetWidth + gap;
    const idx = Math.round(el.scrollLeft / w);
    setActiveIndex(Math.min(Math.max(idx, 0), characters.length - 1));
  }, [characters.length]);

  const previewVoice = useCallback(async () => {
    if (!active || previewing) return;
    setPreviewing(true);
    try {
      const { audio, format } = await previewCharacterVoice(active.id);
      playBase64Audio(audio, format);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không phát được giọng");
    } finally {
      setPreviewing(false);
    }
  }, [active, previewing]);

  const confirm = useCallback(async () => {
    if (!active || selecting) return;
    setSelecting(true);
    setError(null);
    try {
      const fresh = await selectCharacter(active.id);
      applyUser(fresh);
      if (active.voice_id) setTtsVoice(active.voice_id);
      router.replace("/chat");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không lưu được lựa chọn");
      setSelecting(false);
    }
  }, [active, selecting, router, setTtsVoice, applyUser]);

  if (authLoading || loading) {
    return (
      <div className="choose-hani-screen flex flex-col items-center justify-center gap-3">
        <Loader2 className="size-8 animate-spin text-primary" />
        <p className="text-sm text-muted-foreground">Đang mở thế giới Hani…</p>
      </div>
    );
  }

  if (error && characters.length === 0) {
    return (
      <div className="choose-hani-screen flex flex-col items-center justify-center gap-4 px-6 text-center">
        <p className="text-destructive">{error}</p>
        <Button onClick={() => window.location.reload()}>Thử lại</Button>
      </div>
    );
  }

  return (
    <div className="choose-hani-screen">
      <header className="choose-hani-header">
        <p className="choose-hani-eyebrow">Choose Your Hani</p>
        <h1 className="choose-hani-title">Chọn người đồng hành</h1>
        <p className="choose-hani-sub">
          Mỗi người có cách nói chuyện và giọng riêng — chọn người bạn muốn gặp mỗi ngày
        </p>
      </header>

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
            <div className="choose-card-glow" aria-hidden />
            <div className="choose-card-avatar-wrap">
              <div className="choose-card-avatar-inner">
                <Image
                  src={c.avatar_url}
                  alt={c.name}
                  width={280}
                  height={280}
                  className="choose-card-avatar"
                  priority={i === 0}
                />
              </div>
              <span className="choose-card-online">
                <span className="sr-only">Online</span>
              </span>
            </div>

            <div className="choose-card-body">
              <p className="choose-card-name">
                {c.name}
                <span className="choose-card-name-ko">{c.display_name}</span>
              </p>
              <p className="choose-card-style">{c.speaking_style}</p>
              <p className="choose-card-intro">&ldquo;{c.intro_message_ko}&rdquo;</p>
              {c.intro_message_vi ? (
                <p className="choose-card-intro-vi">{c.intro_message_vi}</p>
              ) : null}
              <div className="choose-card-tags">
                <span>{c.emotion_style}</span>
                {c.emoji_style ? <span>{c.emoji_style}</span> : null}
              </div>
            </div>
          </article>
        ))}
      </div>

      <div className="choose-hani-dots" role="tablist" aria-label="Nhân vật">
        {characters.map((c, i) => (
          <button
            key={c.id}
            type="button"
            role="tab"
            aria-selected={i === activeIndex}
            aria-label={c.name}
            className={cn(
              "choose-dot",
              i === activeIndex && "choose-dot-active"
            )}
            onClick={() => {
              const el = scrollerRef.current;
              const card = el?.querySelector<HTMLElement>(".choose-card");
              if (!el || !card) return;
              const gap = 16;
              el.scrollTo({
                left: i * (card.offsetWidth + gap),
                behavior: "smooth",
              });
              setActiveIndex(i);
            }}
          />
        ))}
      </div>

      {error ? (
        <p className="choose-hani-error px-6 text-center text-sm text-destructive">
          {error}
        </p>
      ) : null}

      <footer className="choose-hani-actions">
        <Button
          type="button"
          variant="outline"
          size="lg"
          className="choose-btn-preview"
          disabled={!active || previewing}
          onClick={() => void previewVoice()}
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
          size="lg"
          className="choose-btn-talk"
          disabled={!active || selecting}
          onClick={() => void confirm()}
        >
          {selecting ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Heart className="size-4 fill-current" />
          )}
          {active ? `Trò chuyện với ${active.name}` : "Chọn"}
        </Button>
      </footer>
    </div>
  );
}
