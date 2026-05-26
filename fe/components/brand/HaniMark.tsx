import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

export const HANI_AVATAR_SRC = "/avatar.jpg";

type Props = {
  href?: string;
  size?: "sm" | "md" | "lg" | "xl";
  pulse?: boolean;
  className?: string;
};

const sizeClass = {
  sm: "size-9",
  md: "size-10",
  lg: "size-14",
  xl: "size-20",
} as const;

export function HaniMark({ href, size = "md", pulse, className }: Props) {
  const avatar = (
    <Avatar
      className={cn(
        sizeClass[size],
        "hani-avatar-glow ring-2 ring-primary/25",
        pulse && "animate-pulse ring-primary/50",
        className
      )}
    >
      <AvatarImage src={HANI_AVATAR_SRC} alt="Hani" />
      <AvatarFallback className="bg-primary font-display text-sm font-bold text-primary-foreground">
        하
      </AvatarFallback>
    </Avatar>
  );

  if (href) {
    return (
      <Link
        href={href}
        className="shrink-0 rounded-full outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
        aria-label="Trang chủ"
      >
        {avatar}
      </Link>
    );
  }
  return avatar;
}
