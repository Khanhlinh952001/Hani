"use client";

import { useRef } from "react";
import { Mic } from "lucide-react";
import { ConnectionStatus } from "@/lib/ws/events";
import { BilingualText } from "./BilingualText";
import { cn } from "@/lib/utils";

const hints = {
  active: { ko: "말한 뒤 손을 떼세요", vi: "Nói xong thì thả tay" },
  busy: { ko: "잠시만요…", vi: "Đang xử lý giọng nói…" },
  disabled: { ko: "먼저 연결해 주세요", vi: "Hãy kết nối trước" },
  waitHani: { ko: "하니가 말하는 중", vi: "Đợi Hani trả lời xong" },
  muted: { ko: "마이크가 꺼져 있어요", vi: "Micro đang tắt — bật lại rồi thử" },
  idle: { ko: "누르고 말하기", vi: "Giữ để nói, thả để gửi" },
} as const;

type Props = {
  status: ConnectionStatus;
  isConnected: boolean;
  canPress: boolean;
  holding: boolean;
  busy: boolean;
  showVietnamese: boolean;
  isSourceMuted?: boolean;
  onPressStart: () => void;
  onPressEnd: () => void;
};

export function PushToTalkButton({
  status,
  isConnected,
  canPress,
  holding,
  busy,
  showVietnamese,
  isSourceMuted = false,
  onPressStart,
  onPressEnd,
}: Props) {
  const pressedRef = useRef(false);
  const disabled = !canPress || isSourceMuted;
  const active = holding;

  const hintKey = active
    ? "active"
    : isSourceMuted
      ? "muted"
      : busy
        ? "busy"
        : !isConnected
          ? "disabled"
          : status === "thinking" || status === "speaking"
            ? "waitHani"
            : "idle";

  const hint = hints[hintKey];

  const endPress = (el: HTMLButtonElement, pointerId: number) => {
    if (!pressedRef.current) return;
    pressedRef.current = false;
    if (el.hasPointerCapture(pointerId)) {
      el.releasePointerCapture(pointerId);
    }
    onPressEnd();
  };

  return (
    <div className="flex flex-col items-center gap-4 py-1">
      <BilingualText
        className="text-center text-xs text-muted-foreground"
        ko={hint.ko}
        vi={showVietnamese ? hint.vi : undefined}
      />

      <div className="relative flex items-center justify-center">
        {active && (
          <>
            <span
              className="absolute inset-0 scale-125 animate-ping rounded-full bg-primary/25"
              aria-hidden
            />
            <span
              className="absolute inset-0 scale-110 rounded-full bg-primary/15 blur-md"
              aria-hidden
            />
          </>
        )}
        <button
          type="button"
          disabled={disabled}
          aria-label="Giữ để nói"
          className={cn(
            "ptt-button inline-flex items-center justify-center touch-none select-none disabled:pointer-events-none disabled:opacity-45",
            active && "ptt-button-active",
            isSourceMuted && "ring-2 ring-destructive/40"
          )}
          onPointerDown={(e) => {
            if (disabled || pressedRef.current) return;
            e.preventDefault();
            const el = e.currentTarget;
            el.setPointerCapture(e.pointerId);
            pressedRef.current = true;
            onPressStart();
          }}
          onPointerUp={(e) => {
            endPress(e.currentTarget, e.pointerId);
          }}
          onPointerCancel={(e) => {
            endPress(e.currentTarget, e.pointerId);
          }}
          onLostPointerCapture={(e) => {
            endPress(e.currentTarget, e.pointerId);
          }}
          onContextMenu={(e) => e.preventDefault()}
        >
          <Mic className={cn("size-7", active && "animate-pulse")} />
        </button>
      </div>
    </div>
  );
}
