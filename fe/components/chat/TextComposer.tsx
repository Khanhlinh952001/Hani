"use client";

import { FormEvent, useState } from "react";
import { Send } from "lucide-react";
import { ConnectionStatus } from "@/lib/ws/events";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type Props = {
  status: ConnectionStatus;
  isConnected: boolean;
  showVietnamese: boolean;
  onSend: (text: string) => void;
  chatOnly?: boolean;
};

export function TextComposer({
  status,
  isConnected,
  onSend,
}: Props) {
  const [text, setText] = useState("");
  const canSend =
    isConnected && status === "ready" && text.trim().length > 0;

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const value = text.trim();
    if (!value || !canSend) return;
    onSend(value);
    setText("");
  };

  return (
    <form className="composer-bar" onSubmit={handleSubmit}>
      <Input
        className="composer-input"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="메시지를 입력해요…"
        disabled={!isConnected || status !== "ready"}
        autoComplete="off"
        enterKeyHint="send"
      />
      <Button
        type="submit"
        size="icon"
        className="composer-send"
        disabled={!canSend}
        aria-label="Gửi"
      >
        <Send className="size-4" />
      </Button>
    </form>
  );
}
