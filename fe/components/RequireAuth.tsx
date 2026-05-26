"use client";

import { AuthForm } from "@/components/auth/AuthForm";
import { useAuth } from "@/hooks/useAuth";
import { HaniMark } from "@/components/brand/HaniMark";

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center gap-3">
        <HaniMark size="lg" pulse className="hani-avatar-glow" />
        <p className="text-sm text-muted-foreground">Đang tải…</p>
      </div>
    );
  }

  if (!user) return <AuthForm />;
  return <>{children}</>;
}
