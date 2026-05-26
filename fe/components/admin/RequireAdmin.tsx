"use client";

import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { HaniMark } from "@/components/brand/HaniMark";
import { Button } from "@/components/ui/button";

export function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { user, loading, isAdmin } = useAuth();

  if (loading) {
    return (
      <div className="flex min-h-dvh items-center justify-center">
        <HaniMark size="lg" pulse />
      </div>
    );
  }

  if (!user) return null;

  if (!isAdmin) {
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center gap-4 p-6 text-center">
        <p className="text-lg font-medium">Không có quyền truy cập</p>
        <p className="max-w-sm text-sm text-muted-foreground">
          Trang quản trị chỉ dành cho tài khoản admin.
        </p>
        <Button asChild>
          <Link href="/">Về trang chủ</Link>
        </Button>
      </div>
    );
  }

  return <>{children}</>;
}
