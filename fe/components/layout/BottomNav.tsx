"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, MessageCircle, Mic, Brain, Settings } from "lucide-react";
import { cn } from "@/lib/utils";

const tabs = [
  { href: "/", label: "Home", icon: Home },
  { href: "/chat", label: "Chat", icon: MessageCircle },
  { href: "/speak", label: "Voice", icon: Mic },
  { href: "/memory", label: "Memory", icon: Brain },
  { href: "/settings", label: "Settings", icon: Settings },
] as const;

export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="hani-bottom-nav" aria-label="Điều hướng chính">
      <div className="hani-bottom-nav-inner">
        {tabs.map(({ href, label, icon: Icon }) => {
          const active =
            href === "/"
              ? pathname === "/"
              : pathname === href || pathname.startsWith(`${href}/`);

          return (
            <Link
              key={href}
              href={href}
              className={cn(
                "hani-tab",
                active && "hani-tab-active"
              )}
              aria-current={active ? "page" : undefined}
            >
              <span className="hani-tab-icon-wrap">
                <Icon className="size-[1.125rem]" strokeWidth={active ? 2.5 : 2} />
                {active ? <span className="hani-tab-glow" aria-hidden /> : null}
              </span>
              <span className="hani-tab-label">{label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
