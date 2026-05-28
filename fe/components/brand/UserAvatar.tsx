import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { resolveAvatarSrc } from "@/lib/avatar/resolve-src";
import type { AuthUser } from "@/lib/auth/types";
import { cn } from "@/lib/utils";
import { AvatarPreview } from "./AvatarPreview";

type Props = {
  user: AuthUser | null | undefined;
  size?: "sm" | "md" | "lg";
  className?: string;
  /** Bust cache for uploaded avatars (e.g. updated_at or increment after upload). */
  cacheBust?: string;
  /** Tap to open full image (default: true when user has an avatar). */
  previewable?: boolean;
};

const sizeClass = {
  sm: "size-9",
  md: "size-10",
  lg: "size-14",
} as const;

function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

export function UserAvatar({
  user,
  size = "md",
  className,
  cacheBust,
  previewable = true,
}: Props) {
  const name = user?.name?.trim() || "나";
  const bust = cacheBust ?? user?.updated_at;
  const src = resolveAvatarSrc(user?.avatar, bust);

  const avatar = (
    <Avatar
      className={cn(
        sizeClass[size],
        "hani-avatar-shell",
        src ? "hani-avatar-ring-user" : "ring-2 ring-primary/15",
        previewable && src ? "pointer-events-none" : undefined,
        !previewable || !src ? className : undefined
      )}
    >
      {src ? <AvatarImage src={src} alt={name} /> : null}
      <AvatarFallback className="bg-primary/15 text-sm font-semibold text-primary">
        {initials(name)}
      </AvatarFallback>
    </Avatar>
  );

  if (!previewable || !src) {
    return <div className={cn("shrink-0", className)}>{avatar}</div>;
  }

  return (
    <AvatarPreview src={src} alt={name} className={className}>
      {avatar}
    </AvatarPreview>
  );
}
