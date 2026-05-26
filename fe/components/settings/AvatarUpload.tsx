"use client";

import { useCallback, useRef, useState } from "react";
import { Camera, Loader2 } from "lucide-react";
import { UserAvatar } from "@/components/brand/UserAvatar";
import { useAuth } from "@/hooks/useAuth";
import type { AuthUser } from "@/lib/auth/types";
import { cn } from "@/lib/utils";

const ACCEPT = "image/jpeg,image/png,image/webp,image/gif";

type Props = {
  onSuccess?: (message: string) => void;
  onError?: (message: string) => void;
};

export function AvatarUpload({ onSuccess, onError }: Props) {
  const { user, uploadAvatar } = useAuth();
  const inputRef = useRef<HTMLInputElement>(null);
  const [busy, setBusy] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [cacheBust, setCacheBust] = useState(0);

  const displayUser: AuthUser | null = user
    ? {
        ...user,
        avatar: preview ?? user.avatar,
      }
    : null;

  const handlePick = useCallback(() => {
    if (!busy) inputRef.current?.click();
  }, [busy]);

  const handleFile = useCallback(
    async (file: File | undefined) => {
      if (!file || !user) return;

      if (!file.type.startsWith("image/")) {
        onError?.("Chỉ chọn file ảnh (JPEG, PNG, WebP, GIF).");
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        onError?.("Ảnh tối đa 5MB.");
        return;
      }

      const objectUrl = URL.createObjectURL(file);
      setPreview(objectUrl);
      setBusy(true);

      try {
        await uploadAvatar(file);
        setCacheBust((n) => n + 1);
        onSuccess?.("Đã cập nhật ảnh đại diện.");
      } catch (e) {
        setPreview(null);
        onError?.(e instanceof Error ? e.message : "Tải ảnh thất bại");
      } finally {
        URL.revokeObjectURL(objectUrl);
        setPreview(null);
        setBusy(false);
        if (inputRef.current) inputRef.current.value = "";
      }
    },
    [onError, onSuccess, uploadAvatar, user]
  );

  return (
    <div className="flex flex-col items-center gap-3">
      <input
        ref={inputRef}
        type="file"
        accept={ACCEPT}
        className="sr-only"
        aria-hidden
        onChange={(e) => void handleFile(e.target.files?.[0])}
      />
      <button
        type="button"
        disabled={busy || !user}
        onClick={handlePick}
        className={cn(
          "group relative rounded-full outline-none transition",
          "focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
          busy && "pointer-events-none opacity-70"
        )}
        aria-label="Chọn ảnh đại diện"
      >
        <UserAvatar
          key={cacheBust}
          user={displayUser}
          size="lg"
          className="size-24 ring-4 ring-primary/20"
        />
        <span
          className={cn(
            "absolute inset-0 flex flex-col items-center justify-center gap-1 rounded-full",
            "bg-foreground/45 text-primary-foreground opacity-0 transition group-hover:opacity-100 max-sm:opacity-90",
            busy && "opacity-100"
          )}
        >
          {busy ? (
            <Loader2 className="size-7 animate-spin" aria-hidden />
          ) : (
            <>
              <Camera className="size-7" aria-hidden />
              <span className="text-xs font-medium">Đổi ảnh</span>
            </>
          )}
        </span>
      </button>
      <p className="text-center text-sm text-muted-foreground">
        Nhấn vào ảnh để chọn từ máy · tối đa 5MB
      </p>
    </div>
  );
}
