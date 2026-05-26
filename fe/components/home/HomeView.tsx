"use client";

import Link from "next/link";
import { ChevronRight, MessageCircle, Mic, Settings, Shield } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { HaniMark } from "@/components/brand/HaniMark";
import { AppShell } from "@/components/layout/AppShell";
import { PRACTICE_MODE_OPTIONS } from "@/lib/practice/mode";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

const icons = {
  speak: Mic,
  chat: MessageCircle,
} as const;

export function HomeView() {
  const { user, logout, isAdmin } = useAuth();

  return (
    <AppShell>
      <header className="hani-header flex items-center justify-between gap-3 px-4 py-3">
        <div className="flex items-center gap-3">
          <HaniMark size="md" />
          <div>
            <h1 className="font-display text-lg font-bold leading-tight tracking-tight">
              Hani
            </h1>
            <p className="text-xs font-medium text-primary/75">
              {user ? `안녕, ${user.name}` : "한국어 연습"}
            </p>
          </div>
        </div>
        <div className="flex gap-0.5">
          {isAdmin ? (
            <Button variant="ghost" size="icon-sm" asChild>
              <Link href="/admin" aria-label="Quản trị">
                <Shield className="size-4 text-primary" />
              </Link>
            </Button>
          ) : null}
          <Button variant="ghost" size="icon-sm" asChild>
            <Link href="/settings" aria-label="Cài đặt">
              <Settings className="size-4" />
            </Link>
          </Button>
        </div>
      </header>

      <main className="flex flex-1 flex-col gap-5 overflow-y-auto p-4 pb-6">
        <section className="hani-hero">
          <p className="font-display text-lg font-bold text-foreground">
            오늘도 만나서 반가워요
          </p>
          <p className="mt-1.5 text-sm text-muted-foreground">
            Hôm nay muốn trò chuyện với Hani thế nào?
          </p>
        </section>

        <div className="flex flex-col gap-3">
          {PRACTICE_MODE_OPTIONS.map((opt) => {
            const Icon = icons[opt.id];
            return (
              <Link key={opt.id} href={opt.href} className="group block">
                <Card className="hani-mode-card">
                  <CardHeader className="flex-row items-center gap-4 space-y-0 p-4">
                    <div className="hani-mode-icon">
                      <Icon className="size-5" strokeWidth={2.25} />
                    </div>
                    <div className="min-w-0 flex-1">
                      <CardTitle className="font-display text-xl font-bold tracking-tight text-primary">
                        {opt.ko}
                      </CardTitle>
                      <p className="mt-0.5 text-sm font-medium text-foreground/90">
                        {opt.label}
                      </p>
                      <CardDescription className="mt-1 leading-relaxed">
                        {opt.desc}
                      </CardDescription>
                    </div>
                  </CardHeader>
                </Card>
              </Link>
            );
          })}
        </div>

        <Separator className="bg-primary/10" />
        <p className="text-center text-xs leading-relaxed text-muted-foreground">
          Một cuộc trò chuyện liên tục — Hani nhớ bạn qua cả hai chế độ
        </p>
      </main>
    </AppShell>
  );
}
