"use client";

import { useEffect, useState } from "react";
import { HaniMark } from "@/components/brand/HaniMark";

type Props = {
  showVietnamese?: boolean;
};

const GREETING_KO = "오늘도 만나서 반가워요 💕";
const GREETING_VI = "Hôm nay gặp lại bạn, vui quá";

export function GreetingBubble({ showVietnamese = true }: Props) {
  const [typed, setTyped] = useState("");
  const [done, setDone] = useState(false);

  useEffect(() => {
    let i = 0;
    const id = window.setInterval(() => {
      i += 1;
      setTyped(GREETING_KO.slice(0, i));
      if (i >= GREETING_KO.length) {
        window.clearInterval(id);
        setDone(true);
      }
    }, 42);
    return () => window.clearInterval(id);
  }, []);

  return (
    <div className="home-greeting-row">
      <HaniMark size="sm" className="shrink-0" />
      <div className="home-greeting-bubble">
        <p className="home-greeting-text">
          {typed}
          {!done ? (
            <span className="home-typing-cursor" aria-hidden>
              |
            </span>
          ) : null}
        </p>
        {showVietnamese && done ? (
          <p className="home-greeting-vi">{GREETING_VI}</p>
        ) : null}
        {!done ? (
          <div className="home-typing-dots" aria-hidden>
            <span />
            <span />
            <span />
          </div>
        ) : null}
      </div>
    </div>
  );
}
