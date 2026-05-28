import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import { AvatarPreview } from "./AvatarPreview";

export const HANI_AVATAR_SRC = "/avatar.jpg";

type Props = {
  href?: string;
  size?: "sm" | "md" | "lg" | "xl";
  pulse?: boolean;
  className?: string;
  /** Tap to view full avatar (default: true when not wrapped in a link). */
  previewable?: boolean;
};

const sizeClass = {
  sm: "size-9",
  md: "size-10",
  lg: "size-14",
  xl: "size-20",
} as const;

export function HaniMark({
  href,
  size = "md",
  pulse,
  className,
  previewable,
}: Props) {
  const canPreview = previewable ?? !href;

  const avatar = (
    <Avatar
      className={cn(
        sizeClass[size],
        "hani-avatar-shell hani-avatar-ring",
        pulse && "hani-avatar-pulse",
        canPreview && "pointer-events-none",
        !canPreview ? className : undefined
      )}
    >
      <AvatarImage src={HANI_AVATAR_SRC} alt="Hani" />
      <AvatarFallback className="bg-primary font-display text-sm font-bold text-primary-foreground">
        하
      </AvatarFallback>
    </Avatar>
  );

  const wrapped =
    canPreview ? (
      <AvatarPreview src={HANI_AVATAR_SRC} alt="Hani" className={className}>
        {avatar}
      </AvatarPreview>
    ) : (
      <div className={cn("shrink-0", className)}>{avatar}</div>
    );

  if (href) {
    return (
      <Link
        href={href}
        className="shrink-0 rounded-full outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
        aria-label="Trang chủ"
      >
        {wrapped}
      </Link>
    );
  }

  return wrapped;
}
