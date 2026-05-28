"use client";

import Link from "next/link";
import { Flame, Sparkles, Brain, MessageCircleHeart } from "lucide-react";

type Props = {
  streak: number;
  dailyLabel: string;
  dailyPhrase: string;
  dailyVi: string;
  proactiveKo: string;
  proactiveVi?: string;
};

export function EmotionalStrip({
  streak,
  dailyLabel,
  dailyPhrase,
  dailyVi,
  proactiveKo,
  proactiveVi,
}: Props) {
  return (
    <section className="home-emotional" aria-label="Companion insights">
      <div className="home-emotional-grid">
        <div className="home-chip home-chip-streak">
          <Flame className="size-4 text-primary" />
          <div>
            <p className="home-chip-label">Streak</p>
            <p className="home-chip-value">{streak} ngày</p>
          </div>
        </div>

        <Link href="/memory" className="home-chip home-chip-memory">
          <Brain className="size-4 text-primary" />
          <div>
            <p className="home-chip-label">Memory</p>
            <p className="home-chip-value">Hani nhớ bạn</p>
          </div>
        </Link>
      </div>

      <div className="home-daily-phrase">
        <Sparkles className="size-4 shrink-0 text-primary" />
        <div className="min-w-0 flex-1">
          <p className="home-chip-label">{dailyLabel}</p>
          <p className="home-daily-ko">{dailyPhrase}</p>
          <p className="home-daily-vi">{dailyVi}</p>
        </div>
      </div>

      <div className="home-proactive">
        <MessageCircleHeart className="size-4 shrink-0 text-primary" />
        <div className="min-w-0 flex-1">
          <p className="home-proactive-ko">{proactiveKo}</p>
          {proactiveVi ? (
            <p className="home-proactive-vi">{proactiveVi}</p>
          ) : null}
          <div className="home-typing-indicator" aria-hidden>
            <span />
            <span />
            <span />
          </div>
        </div>
      </div>
    </section>
  );
}
