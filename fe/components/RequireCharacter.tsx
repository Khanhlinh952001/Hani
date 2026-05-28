"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";
import { HaniMark } from "@/components/brand/HaniMark";
import { useAuth } from "@/hooks/useAuth";

const SKIP_PATHS = ["/choose-hani", "/create-lover"];

export function RequireCharacter({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  const needsPick =
    !loading &&
    user &&
    !user.ai_profile_id &&
    !user.selected_character_id &&
    !SKIP_PATHS.some((p) => pathname?.startsWith(p));

  useEffect(() => {
    if (needsPick) {
      router.replace("/create-lover");
    }
  }, [needsPick, router]);

  if (loading || needsPick) {
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center gap-3">
        <HaniMark size="lg" pulse />
        <p className="text-sm text-muted-foreground">Đang tải…</p>
      </div>
    );
  }

  return <>{children}</>;
}
