"use client";

import Link from "next/link";
import { ChevronLeft, Settings } from "lucide-react";
import { ConnectionStatus } from "@/lib/ws/events";
import type { AuthUser } from "@/lib/auth/types";
import type { PracticeMode } from "@/lib/practice/mode";
import { HaniMark } from "@/components/brand/HaniMark";
import { UserAvatar } from "@/components/brand/UserAvatar";
import { useCompanionProfile } from "@/hooks/useCompanionProfile";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function statusDotClass(status: ConnectionStatus, connected: boolean): string {
  if (status === "error") return "bg-destructive";
  if (!connected) return "bg-muted-foreground/40";
  if (status === "thinking" || status === "speaking" || status === "listening") {
    return "bg-primary animate-pulse";
  }
  return "bg-[var(--hani-connected)] shadow-[0_0_8px_var(--hani-connected)]";
}

type Props = {
  practiceMode: PracticeMode;
  status: ConnectionStatus;
  user: AuthUser | null;
  isConnected: boolean;
};

export function ConnectionBar({
  practiceMode,
  status,
  user,
  isConnected,
}: Props) {
  const { avatarUrl, displayName } = useCompanionProfile();
  const pulse =
    practiceMode === "speak" &&
    (status === "speaking" || status === "thinking");

  return (
    <header className="hani-header animate-bar-in flex items-center gap-2 px-2 py-2.5">
      <Button variant="ghost" size="icon-sm" asChild>
        <Link href="/" aria-label="Về trang chủ">
          <ChevronLeft className="size-5" />
        </Link>
      </Button>

      <div className="flex min-w-0 flex-1 items-center gap-2">
        <HaniMark
          href="/"
          size="md"
          pulse={pulse}
          src={avatarUrl}
          alt={displayName}
        />
        <div className="min-w-0">
          <p className="font-display text-base font-bold leading-tight">
            {displayName}
          </p>
          {user ? (
            <p className="truncate text-xs text-muted-foreground">{user.name}</p>
          ) : null}
        </div>
      </div>

      {user ? <UserAvatar user={user} size="sm" /> : null}

      <span
        className={cn(
          "size-2 shrink-0 rounded-full",
          statusDotClass(status, isConnected)
        )}
        title={status}
        aria-hidden
      />

      <Button variant="ghost" size="icon-sm" asChild>
        <Link href="/settings" aria-label="Cài đặt">
          <Settings className="size-4" />
        </Link>
      </Button>
    </header>
  );
}
