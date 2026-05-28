"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { HaniMark } from "@/components/brand/HaniMark";
import { useAuth } from "@/hooks/useAuth";
import { fetchMe } from "@/lib/auth/api";
import { fetchMyLoverProfile } from "@/lib/lover/api";

const SKIP_PATHS = ["/choose-hani", "/create-lover"];

export function RequireCharacter({ children }: { children: React.ReactNode }) {
  const { user, loading, applyUser } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [checkingProfile, setCheckingProfile] = useState(true);

  const onSkipPath = SKIP_PATHS.some((p) => pathname?.startsWith(p));

  useEffect(() => {
    if (loading || !user || onSkipPath) {
      setCheckingProfile(false);
      return;
    }
    if (user.ai_profile_id || user.selected_character_id) {
      setCheckingProfile(false);
      return;
    }

    let cancelled = false;
    setCheckingProfile(true);
    fetchMyLoverProfile()
      .then(async (profile) => {
        if (cancelled) return;
        if (profile) {
          const fresh = await fetchMe();
          if (!cancelled) applyUser(fresh);
          return;
        }
        router.replace("/create-lover");
      })
      .catch(() => {
        if (!cancelled) router.replace("/create-lover");
      })
      .finally(() => {
        if (!cancelled) setCheckingProfile(false);
      });

    return () => {
      cancelled = true;
    };
  }, [loading, user, onSkipPath, router, applyUser]);

  const needsPick =
    !loading &&
    !checkingProfile &&
    user &&
    !user.ai_profile_id &&
    !user.selected_character_id &&
    !onSkipPath;

  useEffect(() => {
    if (needsPick) {
      router.replace("/create-lover");
    }
  }, [needsPick, router]);

  if (loading || checkingProfile || needsPick) {
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center gap-3">
        <HaniMark size="lg" pulse />
        <p className="text-sm text-muted-foreground">Đang tải…</p>
      </div>
    );
  }

  return <>{children}</>;
}
