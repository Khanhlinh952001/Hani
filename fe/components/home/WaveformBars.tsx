"use client";

import { cn } from "@/lib/utils";

export function WaveformBars({ className }: { className?: string }) {
  return (
    <div className={cn("hani-waveform", className)} aria-hidden>
      {Array.from({ length: 5 }).map((_, i) => (
        <span
          key={i}
          className="hani-waveform-bar"
          style={{ animationDelay: `${i * 0.12}s` }}
        />
      ))}
    </div>
  );
}
