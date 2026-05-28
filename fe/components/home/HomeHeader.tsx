"use client";

import Link from "next/link";
import { Settings } from "lucide-react";
import { HaniMark } from "@/components/brand/HaniMark";
import { useCompanionProfile } from "@/hooks/useCompanionProfile";
import { Button } from "@/components/ui/button";

type Props = {
  subtitleKo: string;
  subtitleVi?: string;
  moodEmoji: string;
  moodKo: string;
};

export function HomeHeader({
  subtitleKo,
  subtitleVi,
  moodEmoji,
  moodKo,
}: Props) {
  const { avatarUrl, displayName } = useCompanionProfile();

  return (
    <header className="home-header">
      <div className="home-header-avatar-wrap">
        <HaniMark
          size="xl"
          className="home-header-avatar"
          pulse
          src={avatarUrl}
          alt={displayName}
        />
        <span className="hani-avatar-online" title="Online">
          <span className="sr-only">Online</span>
        </span>
      </div>

      <div className="home-header-copy">
        <p className="home-header-name">{displayName}</p>
        <p className="home-header-status">
          <span className="home-header-mood">
            {moodEmoji} {moodKo}
          </span>
          <span className="home-header-dot" aria-hidden>
            ·
          </span>
          <span className="home-header-sub">{subtitleKo}</span>
        </p>
        {subtitleVi ? (
          <p className="home-header-sub-vi">{subtitleVi}</p>
        ) : null}
      </div>

      <Button
        variant="ghost"
        size="icon-sm"
        className="home-header-settings"
        asChild
      >
        <Link href="/settings" aria-label="Cài đặt">
          <Settings className="size-[1.125rem]" />
        </Link>
      </Button>
    </header>
  );
}
