"use client";

import {
  FormEvent,
  KeyboardEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { ArrowUp, Loader2, Sparkles, WifiOff } from "lucide-react";
import { ConnectionStatus } from "@/lib/ws/events";
import { cn } from "@/lib/utils";

const QUICK_EMOJI = ["💕", "😊", "🥺", "✨", "🌸"] as const;

const MAX_LINES = 5;

type Props = {
  status: ConnectionStatus;
  isConnected: boolean;
  showVietnamese: boolean;
  onSend: (text: string) => void;
  chatOnly?: boolean;
};

function statusHint(
  status: ConnectionStatus,
  isConnected: boolean,
  showVi: boolean
): { ko: string; vi?: string } | null {
  if (!isConnected) {
    return { ko: "연결 중이에요…", vi: showVi ? "Đang kết nối…" : undefined };
  }
  if (status === "thinking") {
    return { ko: "하니가 답을 쓰고 있어요", vi: showVi ? "Hani đang trả lời…" : undefined };
  }
  if (status === "speaking" || status === "listening") {
    return { ko: "잠시만 기다려 주세요", vi: showVi ? "Đợi Hani một chút nhé" : undefined };
  }
  if (status === "connecting") {
    return { ko: "채팅방 준비 중…", vi: showVi ? "Đang mở cuộc trò chuyện…" : undefined };
  }
  if (status === "error") {
    return { ko: "다시 시도해 주세요", vi: showVi ? "Thử kết nối lại nhé" : undefined };
  }
  return null;
}

export function TextComposer({
  status,
  isConnected,
  showVietnamese,
  onSend,
}: Props) {
  const [text, setText] = useState("");
  const [emojiOpen, setEmojiOpen] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const emojiRef = useRef<HTMLDivElement>(null);

  const ready = isConnected && status === "ready";
  const trimmed = text.trim();
  const canSend = ready && trimmed.length > 0;
  const hint = statusHint(status, isConnected, showVietnamese);

  const resize = useCallback(() => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = "auto";
    const lineHeight = 22;
    const maxHeight = lineHeight * MAX_LINES + 12;
    el.style.height = `${Math.min(el.scrollHeight, maxHeight)}px`;
  }, []);

  useEffect(() => {
    resize();
  }, [text, resize]);

  useEffect(() => {
    if (!emojiOpen) return;
    const onPointerDown = (e: MouseEvent) => {
      if (emojiRef.current?.contains(e.target as Node)) return;
      setEmojiOpen(false);
    };
    document.addEventListener("pointerdown", onPointerDown);
    return () => document.removeEventListener("pointerdown", onPointerDown);
  }, [emojiOpen]);

  const insertAtCursor = useCallback((chunk: string) => {
    const el = textareaRef.current;
    if (!el) {
      setText((prev) => prev + chunk);
      return;
    }
    const start = el.selectionStart ?? text.length;
    const end = el.selectionEnd ?? text.length;
    const next = text.slice(0, start) + chunk + text.slice(end);
    setText(next);
    requestAnimationFrame(() => {
      el.focus();
      const pos = start + chunk.length;
      el.setSelectionRange(pos, pos);
      resize();
    });
  }, [text, resize]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!canSend) return;
    onSend(trimmed);
    setText("");
    setEmojiOpen(false);
    requestAnimationFrame(resize);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (canSend) {
        onSend(trimmed);
        setText("");
        setEmojiOpen(false);
        requestAnimationFrame(resize);
      }
    }
  };

  return (
    <div className="composer-root">
      {hint ? (
        <div
          className="composer-status"
          role="status"
          aria-live="polite"
        >
          {!isConnected || status === "connecting" ? (
            <Loader2 className="size-3.5 shrink-0 animate-spin text-primary/70" />
          ) : status === "error" ? (
            <WifiOff className="size-3.5 shrink-0 text-destructive/80" />
          ) : (
            <span className="composer-status-dot" aria-hidden />
          )}
          <span className="composer-status-ko">{hint.ko}</span>
          {showVietnamese && hint.vi ? (
            <span className="composer-status-vi">{hint.vi}</span>
          ) : null}
        </div>
      ) : null}

      <form className="composer-dock" onSubmit={handleSubmit}>
        <div className="composer-dock-accent" aria-hidden />

        <div className="composer-field-wrap">
          <div className="composer-leading">
            <button
              type="button"
              className={cn(
                "composer-icon-btn",
                emojiOpen && "composer-icon-btn-active"
              )}
              disabled={!ready}
              aria-label="Biểu cảm"
              aria-expanded={emojiOpen}
              onClick={() => setEmojiOpen((o) => !o)}
            >
              <Sparkles className="size-4.5" strokeWidth={2} />
            </button>

            {emojiOpen ? (
              <div ref={emojiRef} className="composer-emoji-pop">
                {QUICK_EMOJI.map((e) => (
                  <button
                    key={e}
                    type="button"
                    className="composer-emoji-btn"
                    onClick={() => {
                      insertAtCursor(e);
                      setEmojiOpen(false);
                    }}
                  >
                    {e}
                  </button>
                ))}
              </div>
            ) : null}
          </div>

          <textarea
            ref={textareaRef}
            className="composer-textarea"
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={
              ready
                ? "메시지를 입력해요…"
                : "잠시만 기다려 주세요…"
            }
            disabled={!ready}
            autoComplete="off"
            rows={1}
            enterKeyHint="send"
            aria-label="Tin nhắn"
          />

          <button
            type="submit"
            className={cn(
              "composer-send",
              canSend && "composer-send-ready"
            )}
            disabled={!canSend}
            aria-label="Gửi tin nhắn"
          >
            <ArrowUp className="size-4.5" strokeWidth={2.75} />
          </button>
        </div>
      </form>
    </div>
  );
}
