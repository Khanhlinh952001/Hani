"use client";

import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  ArrowLeft,
  Brain,
  MessageSquare,
  RefreshCw,
  Shield,
  Trash2,
  Users,
} from "lucide-react";
import {
  clearUserConversation,
  clearUserMemories,
  deleteAdminUser,
  fetchAdminStats,
  fetchAdminUsers,
  fetchSessionMessages,
  fetchUserMemories,
  patchAdminUser,
  type AdminMemory,
  type AdminMessage,
  type AdminStats,
  type AdminUser,
} from "@/lib/admin/api";
import { ROLE_ADMIN } from "@/lib/auth/types";
import { AppShell } from "@/components/layout/AppShell";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

export function AdminView() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [selected, setSelected] = useState<AdminUser | null>(null);
  const [memories, setMemories] = useState<AdminMemory[]>([]);
  const [messages, setMessages] = useState<AdminMessage[]>([]);
  const [query, setQuery] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const [s, u] = await Promise.all([fetchAdminStats(), fetchAdminUsers()]);
    setStats(s);
    setUsers(u);
  }, []);

  useEffect(() => {
    void load().catch((e) =>
      setError(e instanceof Error ? e.message : "Tải thất bại")
    );
  }, [load]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return users;
    return users.filter(
      (u) =>
        u.name.toLowerCase().includes(q) ||
        u.email.toLowerCase().includes(q) ||
        String(u.id).includes(q)
    );
  }, [users, query]);

  const selectUser = useCallback(async (u: AdminUser) => {
    setSelected(u);
    setMessages([]);
    try {
      const mem = await fetchUserMemories(u.id);
      setMemories(mem);
    } catch {
      setMemories([]);
    }
  }, []);

  const refreshSelected = useCallback(async () => {
    if (!selected) return;
    await selectUser(selected);
    const fresh = await fetchAdminUsers();
    setUsers(fresh);
    const updated = fresh.find((u) => u.id === selected.id);
    if (updated) setSelected(updated);
  }, [selected, selectUser]);

  const runAction = useCallback(
    async (fn: () => Promise<void>) => {
      setBusy(true);
      setError(null);
      try {
        await fn();
        await load();
        await refreshSelected();
      } catch (e) {
        setError(e instanceof Error ? e.message : "Thao tác thất bại");
      } finally {
        setBusy(false);
      }
    },
    [load, refreshSelected]
  );

  const loadMessages = useCallback(async (sessionId: string) => {
    try {
      const list = await fetchSessionMessages(sessionId);
      setMessages(list.slice(-20));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không tải được tin nhắn");
    }
  }, []);

  return (
    <AppShell className="max-w-2xl">
      <header className="hani-header flex items-center gap-2 px-3 py-2.5">
        <Button variant="ghost" size="icon-sm" asChild>
          <Link href="/" aria-label="Về trang chủ">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <Shield className="size-5 text-primary" />
        <div className="min-w-0 flex-1">
          <h1 className="font-display text-lg font-bold">Quản trị Hani</h1>
          <p className="text-xs text-muted-foreground">Users · ký ức · hội thoại</p>
        </div>
        <Button
          variant="ghost"
          size="icon-sm"
          disabled={busy}
          onClick={() => void load()}
          aria-label="Làm mới"
        >
          <RefreshCw className={cn("size-4", busy && "animate-spin")} />
        </Button>
      </header>

      <main className="flex-1 space-y-4 overflow-y-auto p-4 pb-8">
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {stats && (
          <div className="grid grid-cols-2 gap-2">
            {[
              { label: "Người dùng", value: stats.users, icon: Users },
              { label: "Phiên chat", value: stats.sessions, icon: MessageSquare },
              { label: "Tin nhắn", value: stats.messages, icon: MessageSquare },
              { label: "Ký ức", value: stats.memories, icon: Brain },
            ].map((item) => (
              <Card key={item.label} className="border-primary/10 py-3">
                <CardContent className="flex items-center gap-3 px-4 py-0">
                  <item.icon className="size-4 text-primary" />
                  <div>
                    <p className="text-xs text-muted-foreground">{item.label}</p>
                    <p className="text-xl font-bold">{item.value}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base">Người dùng</CardTitle>
            <CardDescription>Chọn để xem chi tiết và thao tác</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Input
              placeholder="Tìm tên, email, id…"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
            <ul className="max-h-48 space-y-1 overflow-y-auto">
              {filtered.map((u) => (
                <li key={u.id}>
                  <button
                    type="button"
                    onClick={() => void selectUser(u)}
                    className={cn(
                      "flex w-full items-center justify-between gap-2 rounded-lg border px-3 py-2 text-left text-sm transition",
                      selected?.id === u.id
                        ? "border-primary/40 bg-primary/10"
                        : "border-transparent hover:bg-muted/60"
                    )}
                  >
                    <span className="min-w-0 truncate font-medium">{u.name}</span>
                    <span className="shrink-0 text-xs text-muted-foreground">
                      #{u.id}
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>

        {selected && (
          <Card className="border-primary/20">
            <CardHeader>
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <CardTitle>{selected.name}</CardTitle>
                  <CardDescription>{selected.email}</CardDescription>
                </div>
                <div className="flex flex-wrap gap-1">
                  <Badge variant={selected.role === ROLE_ADMIN ? "default" : "secondary"}>
                    {selected.role === ROLE_ADMIN ? "Admin" : "User"}
                  </Badge>
                  <Badge variant={selected.status === 1 ? "outline" : "destructive"}>
                    {selected.status === 1 ? "Hoạt động" : "Khóa"}
                  </Badge>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex flex-wrap gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  disabled={busy}
                  onClick={() =>
                    void runAction(() =>
                      patchAdminUser(selected.id, {
                        status: selected.status === 1 ? 0 : 1,
                      }).then(() => undefined)
                    )
                  }
                >
                  {selected.status === 1 ? "Khóa tài khoản" : "Mở khóa"}
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  disabled={busy}
                  onClick={() =>
                    void runAction(() =>
                      patchAdminUser(selected.id, {
                        role: selected.role === ROLE_ADMIN ? 0 : ROLE_ADMIN,
                      }).then(() => undefined)
                    )
                  }
                >
                  {selected.role === ROLE_ADMIN ? "Bỏ admin" : "Làm admin"}
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  disabled={busy}
                  onClick={() =>
                    void runAction(() => clearUserConversation(selected.id))
                  }
                >
                  Xóa hội thoại
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  disabled={busy}
                  onClick={() =>
                    void runAction(() => clearUserMemories(selected.id))
                  }
                >
                  Xóa ký ức
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  disabled={busy}
                  onClick={() => {
                    if (
                      !confirm(
                        `Xóa vĩnh viễn ${selected.name}? Không hoàn tác được.`
                      )
                    ) {
                      return;
                    }
                    void runAction(async () => {
                      await deleteAdminUser(selected.id);
                      setSelected(null);
                      setMemories([]);
                      setMessages([]);
                    });
                  }}
                >
                  <Trash2 className="size-3.5" />
                  Xóa user
                </Button>
              </div>

              <Separator />

              <div>
                <p className="mb-2 text-sm font-medium">
                  Ký ức ({memories.length})
                </p>
                <ul className="max-h-36 space-y-1 overflow-y-auto text-xs text-muted-foreground">
                  {memories.length === 0 ? (
                    <li>Chưa có ký ức</li>
                  ) : (
                    memories.slice(0, 15).map((m) => (
                      <li key={m.id} className="rounded-md bg-muted/50 px-2 py-1">
                        <span className="text-primary/80">{m.memory_type || "fact"}</span>
                        {" · "}
                        {m.content}
                      </li>
                    ))
                  )}
                </ul>
              </div>

              <div>
                <p className="mb-2 text-sm font-medium">Tin nhắn gần đây</p>
                <Button
                  size="sm"
                  variant="secondary"
                  className="mb-2"
                  onClick={() => {
                    const sid = prompt(
                      "Dán session_id (UUID) để xem tin nhắn:"
                    );
                    if (sid) void loadMessages(sid.trim());
                  }}
                >
                  Tải theo session ID
                </Button>
                <ul className="max-h-40 space-y-1 overflow-y-auto text-xs">
                  {messages.length === 0 ? (
                    <li className="text-muted-foreground">
                      Chưa tải — dùng session ID từ DB
                    </li>
                  ) : (
                    messages.map((m) => (
                      <li
                        key={m.id}
                        className={cn(
                          "rounded-md px-2 py-1",
                          m.role === "user" ? "bg-primary/10" : "bg-muted/50"
                        )}
                      >
                        <strong>{m.role}:</strong> {m.content}
                      </li>
                    ))
                  )}
                </ul>
              </div>
            </CardContent>
          </Card>
        )}
      </main>
    </AppShell>
  );
}
