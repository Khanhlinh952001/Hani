"use client";

import { useEffect, useRef } from "react";
import type { PracticeMode } from "@/lib/practice/mode";
import type { AuthUser } from "@/lib/auth/types";
import { ChatMessage } from "@/lib/ws/events";
import { HaniMark } from "@/components/brand/HaniMark";
import { UserAvatar } from "@/components/brand/UserAvatar";
import { BilingualText } from "./BilingualText";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type Props = {
  user: AuthUser | null;
  messages: ChatMessage[];
  partial: string;
  partialVi?: string;
  status: string;
  holding?: boolean;
  showVietnamese: boolean;
  practiceMode: PracticeMode;
};

export function MessageList({
  user,
  messages,
  partial,
  partialVi,
  status,
  holding = false,
  showVietnamese,
  practiceMode,
}: Props) {
  const isSpeakMode = practiceMode === "speak";
  const isListeningLive =
    isSpeakMode && (holding || status === "listening");
  const showPartial =
    isListeningLive && (!!partial || holding);
  const endRef = useRef<HTMLDivElement>(null);

  const visibleMessages = isSpeakMode
    ? messages
    : messages.filter((msg) => !(msg.role === "assistant" && msg.streaming));

  const lastMsg = visibleMessages[visibleMessages.length - 1];
  const awaitingAssistant =
    status === "thinking" &&
    lastMsg?.role === "assistant" &&
    !!lastMsg.content;

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, partial, status]);

  return (
    <div
      className={cn(
        "flex min-h-full flex-col gap-3.5 pb-4",
        !isSpeakMode && "message-list--chat"
      )}
    >
      {visibleMessages.length === 0 &&
        !showPartial &&
        status !== "thinking" &&
        status !== "speaking" && (
          <div className="hani-empty">
            <HaniMark size="xl" className="hani-avatar-glow" pulse />
            <BilingualText
              className="max-w-[16rem]"
              ko={
                status === "connecting"
                  ? "연결 중이에요, 잠시만요…"
                  : status === "error"
                    ? "연결에 실패했어요"
                    : isSpeakMode
                      ? "곧 연락할게 — 버튼을 길게 눌러 말해요"
                      : "곧 연락할게 — 아래에 메시지를 입력해요"
              }
              vi={
                showVietnamese
                  ? status === "connecting"
                    ? "Đang kết nối Hani…"
                    : status === "error"
                      ? "Kết nối thất bại"
                      : isSpeakMode
                        ? "Giữ nút để trả lời Hani"
                        : "Gõ tin nhắn bên dưới"
                  : undefined
              }
            />
          </div>
        )}

      {visibleMessages.map((msg, i) => {
        const isUser = msg.role === "user";
        const enterClass = isSpeakMode
          ? isUser
            ? "bubble-enter-user"
            : "bubble-enter-hani"
          : msg.justArrived
            ? isUser
              ? "bubble-msg-pop-user"
              : "bubble-msg-pop-hani"
            : "";

        return (
          <div
            key={msg.id}
            className={cn(
              "flex items-end gap-2",
              isUser ? "flex-row-reverse" : "flex-row"
            )}
            style={{ "--bubble-i": i } as React.CSSProperties}
          >
            {isUser ? (
              <UserAvatar user={user} size="sm" className="mb-0.5 shrink-0" />
            ) : (
              <HaniMark size="sm" className="mb-0.5 shrink-0" />
            )}
            <article
              className={cn(
                "bubble min-w-0",
                isUser ? "bubble-user" : "bubble-hani",
                enterClass,
                isSpeakMode && msg.streaming && "bubble-streaming"
              )}
            >
             
              <BilingualText
                ko={msg.content}
                vi={showVietnamese ? msg.translationVi : undefined}
              />
            </article>
          </div>
        );
      })}

      {showPartial && (
        <div className="flex flex-row-reverse items-end gap-2">
          <UserAvatar user={user} size="sm" className="mb-0.5 shrink-0" />
          <article className="bubble bubble-user bubble-partial bubble-enter-user min-w-0">
            <Badge variant="secondary" className="mb-1.5 h-5 px-1.5 text-[10px]">
              나 · 듣는 중
            </Badge>
            <BilingualText
              ko={partial || (holding ? "…" : "")}
              vi={showVietnamese ? partialVi : undefined}
            />
            {holding || status === "listening" ? (
              <span className="cursor-blink">▍</span>
            ) : null}
          </article>
        </div>
      )}

      {status === "thinking" &&
        !showPartial &&
        !awaitingAssistant &&
        !visibleMessages.some((m) => m.streaming) && (
          <div className="flex items-end gap-2">
            <HaniMark size="sm" pulse className="mb-0.5 shrink-0" />
            <article className="bubble bubble-hani bubble-thinking bubble-enter-hani">
              <Badge variant="outline" className="mb-1.5 h-5 px-1.5 text-[10px]">
                하니
              </Badge>
              <div className="typing-dots">
                <span />
                <span />
                <span />
              </div>
            </article>
          </div>
        )}

      <div ref={endRef} className="h-px shrink-0" aria-hidden />
    </div>
  );
}
