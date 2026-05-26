import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { API_URL } from "@/lib/config";
import type { AuthUser } from "@/lib/auth/types";
import { cn } from "@/lib/utils";

function resolveAvatarSrc(
  avatar?: string,
  cacheBust?: string
): string | undefined {
  const raw = avatar?.trim();
  if (!raw) return undefined;

  let url: string;
  if (
    raw.startsWith("http://") ||
    raw.startsWith("https://") ||
    raw.startsWith("blob:")
  ) {
    url = raw;
  } else if (raw.startsWith("/uploads/")) {
    url = `${API_URL}${raw}`;
  } else if (raw.startsWith("/")) {
    url = raw;
  } else {
    url = raw;
  }

  if (cacheBust && url.includes("/uploads/")) {
    const sep = url.includes("?") ? "&" : "?";
    return `${url}${sep}v=${encodeURIComponent(cacheBust)}`;
  }
  return url;
}

type Props = {
  user: AuthUser | null | undefined;
  size?: "sm" | "md" | "lg";
  className?: string;
  /** Bust cache for uploaded avatars (e.g. updated_at or increment after upload). */
  cacheBust?: string;
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
}: Props) {
  const name = user?.name?.trim() || "나";
  const bust = cacheBust ?? user?.updated_at;
  const src = resolveAvatarSrc(user?.avatar, bust);

  return (
    <Avatar
      className={cn(
        sizeClass[size],
        "ring-2 ring-background shadow-sm",
        className
      )}
    >
      {src ? <AvatarImage src={src} alt={name} /> : null}
      <AvatarFallback className="bg-primary/15 text-sm font-semibold text-primary">
        {initials(name)}
      </AvatarFallback>
    </Avatar>
  );
}
