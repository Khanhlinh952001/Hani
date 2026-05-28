"use client";

import Link from "next/link";
import { Mic } from "lucide-react";
import { WaveformBars } from "./WaveformBars";

type Props = {
  voiceMinutes: number;
  onNavigate?: () => void;
};

export function VoicePracticeCard({ voiceMinutes, onNavigate }: Props) {
  const stat =
    voiceMinutes > 0
      ? `Hôm nay đã nói ${voiceMinutes} phút`
      : "Hôm nay chưa nói — thử nhé";

  return (
    <Link href="/speak" className="home-action-card home-action-voice" onClick={onNavigate}>
      <div className="home-action-top">
        <div className="home-action-icon home-action-icon-voice">
          <Mic className="size-5" strokeWidth={2.25} />
        </div>
        <div className="home-action-live">
          <span className="home-live-dot" />
          Live
        </div>
      </div>

      <div className="home-action-body">
        <p className="home-action-ko">말하기 연습</p>
        <p className="home-action-title">Voice Practice</p>
        <p className="home-action-desc">Giữ để nói chuyện với Hani</p>
      </div>

      <div className="home-action-footer">
        <WaveformBars />
        <span className="home-action-stat">{stat}</span>
      </div>
    </Link>
  );
}
