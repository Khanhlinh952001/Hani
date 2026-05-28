import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { HANI_BRAND_LOGO } from "@/lib/brand/assets";
import { cn } from "@/lib/utils";
import { AvatarPreview } from "./AvatarPreview";

type Props = {
  href?: string;
  size?: "sm" | "md" | "lg" | "xl" | "2xl";
  pulse?: boolean;
  className?: string;
  /** Companion portrait URL; omit for app brand logo from /logo/. */
  src?: string;
  alt?: string;
  /** Tap to view full image (default: true when not wrapped in a link). */
  previewable?: boolean;
};

const sizeClass = {
  sm: "size-9",
  md: "size-10",
  lg: "size-14",
  xl: "size-20",
  "2xl": "size-24",
} as const;

export function HaniMark({
  href,
  size = "md",
  pulse,
  className,
  src,
  alt = "Hani",
  previewable,
}: Props) {
  const imageSrc = src?.trim() || HANI_BRAND_LOGO;
  const isBrand = !src?.trim();
  const canPreview = previewable ?? !href;

  const avatar = isBrand ? (
    <div
      className={cn(
        "hani-brand-mark-shell shrink-0 overflow-hidden rounded-2xl",
        sizeClass[size],
        pulse && "hani-avatar-pulse",
        canPreview && "pointer-events-none",
        !canPreview ? className : undefined
      )}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={imageSrc}
        alt={alt}
        className="size-full object-contain"
        draggable={false}
      />
    </div>
  ) : (
    <Avatar
      className={cn(
        sizeClass[size],
        "hani-avatar-shell hani-avatar-ring",
        pulse && "hani-avatar-pulse",
        canPreview && "pointer-events-none",
        !canPreview ? className : undefined
      )}
    >
      <AvatarImage src={imageSrc} alt={alt} className="object-cover" />
      <AvatarFallback className="bg-primary font-display text-sm font-bold text-primary-foreground">
        H
      </AvatarFallback>
    </Avatar>
  );

  const wrapped = canPreview ? (
    <AvatarPreview src={imageSrc} alt={alt} className={className}>
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
