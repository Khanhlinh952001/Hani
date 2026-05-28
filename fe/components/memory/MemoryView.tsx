"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { ArrowLeft, Brain, Heart, Loader2, Sparkles } from "lucide-react";
import { CompanionLayout } from "@/components/layout/CompanionLayout";
import { HaniMark } from "@/components/brand/HaniMark";
import { useCompanionHome } from "@/hooks/useCompanionHome";
import { useSettings } from "@/hooks/useSettings";
import { fetchMemories, type Memory } from "@/lib/memories/api";
import { BilingualText } from "@/components/chat/BilingualText";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

function memoryIcon(type?: string) {
  if (type === "emotional") {
    return <Heart className="size-3.5" />;
  }
  return <Sparkles className="size-3.5" />;
}

function formatWhen(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const days = Math.floor(diffMs / 86_400_000);
  if (days === 0) return "Hôm nay";
  if (days === 1) return "Hôm qua";
  if (days < 7) return `${days} ngày trước`;
  return d.toLocaleDateString("vi-VN", {
    day: "numeric",
    month: "short",
  });
}

export function MemoryView() {
  const { showVietnamese } = useSettings();
  const home = useCompanionHome();
  const [memories, setMemories] = useState<Memory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const list = await fetchMemories();
      setMemories(list);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không tải được ký ức");
      setMemories([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  return (
    <CompanionLayout>
      <header className="memory-header">
        <Button variant="ghost" size="icon-sm" asChild>
          <Link href="/" aria-label="Về trang chủ">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <div className="flex min-w-0 flex-1 items-center gap-2">
          <HaniMark size="md" />
          <div>
            <h1 className="font-display text-lg font-bold">Ký ức</h1>
            <p className="text-xs text-muted-foreground">Hani nhớ về bạn</p>
          </div>
        </div>
      </header>

      <main className="memory-main">
        <div className="memory-hero">
          <Brain className="size-8 text-primary" />
          <p className="font-display text-lg font-bold">
            {home.mood.emoji} {home.mood.ko}
          </p>
          <p className="text-sm text-muted-foreground">
            Những khoảnh khắc Hani giữ lại để trò chuyện tự nhiên hơn
          </p>
          {!loading && !error ? (
            <p className="text-xs font-medium text-primary">
              {memories.length} ký ức
            </p>
          ) : null}
        </div>

        {error ? (
          <Alert variant="destructive">
            <AlertDescription className="flex flex-col gap-2">
              {error}
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="w-fit"
                onClick={() => void load()}
              >
                Thử lại
              </Button>
            </AlertDescription>
          </Alert>
        ) : null}

        {loading ? (
          <div className="memory-loading" role="status">
            <Loader2 className="size-6 animate-spin text-primary" />
            <p className="text-sm text-muted-foreground">Đang tải ký ức…</p>
          </div>
        ) : null}

        {!loading && !error && memories.length === 0 ? (
          <div className="memory-empty">
            <p className="font-medium text-foreground">
              Hani chưa có ký ức nào về bạn
            </p>
            <p className="text-sm text-muted-foreground">
              Trò chuyện thêm — Hani sẽ nhớ sở thích và khoảnh khắc của bạn.
            </p>
            <Button asChild className="mt-2">
              <Link href="/chat">Bắt đầu trò chuyện</Link>
            </Button>
          </div>
        ) : null}

        {!loading && !error && memories.length > 0 ? (
          <ul className="memory-list">
            {memories.map((m) => (
              <li key={m.id} className="memory-item">
                <span className="memory-item-icon">
                  {memoryIcon(m.memory_type)}
                </span>
                <div className="min-w-0 flex-1">
                  <BilingualText
                    ko={m.content}
                    vi={
                      showVietnamese
                        ? m.translation_vi || undefined
                        : undefined
                    }
                  />
                  <p className="memory-item-meta">
                    {m.memory_type === "emotional"
                      ? "Cảm xúc"
                      : m.memory_type === "life"
                        ? "Cuộc sống"
                        : "Ký ức"}
                    {m.created_at ? ` · ${formatWhen(m.created_at)}` : null}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        ) : null}

        <p className="memory-footnote">
          Ký ức mới lưu bằng tiếng Hàn + tiếng Việt. Bản cũ (tiếng Anh) vẫn hiện
          cho đến khi bạn xóa và trò chuyện lại.{" "}
          <Link href="/settings" className="text-primary underline-offset-2 hover:underline">
            Cài đặt
          </Link>{" "}
          để xóa toàn bộ.
        </p>
      </main>
    </CompanionLayout>
  );
}
