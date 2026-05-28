"use client";

import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  ArrowLeft,
  Brain,
  Crown,
  Lock,
  MessageSquare,
  RefreshCw,
  RotateCcw,
  Search,
  Shield,
  Sparkles,
  Trash2,
  Unlock,
  UserCog,
  Users,
  Zap,
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
  resetUserUsage,
  type AdminMemory,
  type AdminMessage,
  type AdminStats,
  type AdminUser,
  type UsageSnapshot,
} from "@/lib/admin/api";
import { ROLE_ADMIN } from "@/lib/auth/types";
import { AppShell } from "@/components/layout/AppShell";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const PLAN_OPTIONS = [
  {
    id: "free",
    label: "Free",
    desc: "30 tin / ngày",
    icon: Users,
  },
  {
    id: "plus",
    label: "Plus",
    desc: "1000 tin · 60p voice",
    icon: Zap,
  },
  {
    id: "premium",
    label: "Premium",
    desc: "Không giới hạn",
    icon: Crown,
  },
] as const;

function planPillClass(plan?: string) {
  const p = plan || "free";
  return cn(
    "hani-admin-plan-pill",
    p === "premium" && "hani-admin-plan-pill--premium",
    p === "plus" && "hani-admin-plan-pill--plus",
    p === "free" && "hani-admin-plan-pill--free"
  );
}

function userInitials(name: string) {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }
  return (name.slice(0, 2) || "H").toUpperCase();
}

function StatCard({
  label,
  value,
  icon: Icon,
}: {
  label: string;
  value: number;
  icon: typeof Users;
}) {
  return (
    <div className="hani-admin-stat">
      <div className="hani-admin-stat-icon">
        <Icon className="size-5" strokeWidth={2} />
      </div>
      <p className="text-xs font-medium text-muted-foreground">{label}</p>
      <p className="hani-admin-stat-value">{value.toLocaleString("vi-VN")}</p>
    </div>
  );
}

export function AdminView() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [selected, setSelected] = useState<AdminUser | null>(null);
  const [memories, setMemories] = useState<AdminMemory[]>([]);
  const [messages, setMessages] = useState<AdminMessage[]>([]);
  const [query, setQuery] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [usage, setUsage] = useState<UsageSnapshot | null>(null);

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
        String(u.id).includes(q) ||
        (u.subscription_plan || "free").includes(q)
    );
  }, [users, query]);

  const selectUser = useCallback(async (u: AdminUser) => {
    setSelected(u);
    setMessages([]);
    setUsage(null);
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

  const currentPlan = selected?.subscription_plan || "free";

  return (
    <AppShell className="hani-admin-shell">
      <header className="hani-admin-header">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon-sm" asChild className="shrink-0">
            <Link href="/" aria-label="Về trang chủ">
              <ArrowLeft className="size-4" />
            </Link>
          </Button>
          <div className="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-primary/10 text-primary shadow-sm">
            <Shield className="size-5" strokeWidth={2.2} />
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-[0.625rem] font-semibold uppercase tracking-widest text-primary/80">
              Hani Console
            </p>
            <h1 className="font-display text-xl font-bold tracking-tight">
              Quản trị
            </h1>
          </div>
          <Button
            variant="outline"
            size="sm"
            disabled={busy}
            onClick={() => void load()}
            className="shrink-0 gap-1.5 rounded-full border-primary/20"
          >
            <RefreshCw className={cn("size-3.5", busy && "animate-spin")} />
            Làm mới
          </Button>
        </div>

        {stats && (
          <div className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
            <StatCard label="Người dùng" value={stats.users} icon={Users} />
            <StatCard label="Phiên chat" value={stats.sessions} icon={MessageSquare} />
            <StatCard label="Tin nhắn" value={stats.messages} icon={Sparkles} />
            <StatCard label="Ký ức" value={stats.memories} icon={Brain} />
          </div>
        )}
      </header>

      <div className="hani-admin-layout">
        {error && (
          <Alert variant="destructive" className="col-span-full">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Sidebar — user list */}
        <aside className="hani-admin-sidebar">
          <div className="border-b border-primary/8 p-4">
            <h2 className="font-display text-sm font-bold">Người dùng</h2>
            <p className="text-xs text-muted-foreground">
              {filtered.length} / {users.length} tài khoản
            </p>
            <div className="relative mt-3">
              <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                className="rounded-xl border-primary/12 bg-background pl-9"
                placeholder="Tìm tên, email, gói…"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
            </div>
          </div>
          <ul className="hani-admin-scroll space-y-1 p-2">
            {filtered.length === 0 ? (
              <li className="px-3 py-8 text-center text-sm text-muted-foreground">
                Không tìm thấy user
              </li>
            ) : (
              filtered.map((u) => (
                <li key={u.id}>
                  <button
                    type="button"
                    onClick={() => void selectUser(u)}
                    className={cn(
                      "hani-admin-user-row",
                      selected?.id === u.id && "hani-admin-user-row--active"
                    )}
                  >
                    <span className="hani-admin-avatar" aria-hidden>
                      {userInitials(u.name)}
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="block truncate text-sm font-semibold">
                        {u.name}
                      </span>
                      <span className="block truncate text-xs text-muted-foreground">
                        {u.email}
                      </span>
                    </span>
                    <span className={planPillClass(u.subscription_plan)}>
                      {u.subscription_plan || "free"}
                    </span>
                  </button>
                </li>
              ))
            )}
          </ul>
        </aside>

        {/* Detail panel */}
        <section className="min-h-0">
          {!selected ? (
            <div className="hani-admin-detail-empty">
              <div className="flex size-14 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                <UserCog className="size-7" />
              </div>
              <p className="font-display text-base font-semibold">
                Chọn một người dùng
              </p>
              <p className="max-w-xs text-sm text-muted-foreground">
                Xem gói đăng ký, đổi plan, reset hạn mức hoặc quản lý dữ liệu chat.
              </p>
            </div>
          ) : (
            <div className="hani-admin-detail">
              {/* Profile hero */}
              <div className="border-b border-primary/10 bg-linear-to-br from-primary/8 via-transparent to-secondary/40 p-5">
                <div className="flex flex-wrap items-start gap-4">
                  <span className="hani-admin-avatar size-14 text-lg">
                    {userInitials(selected.name)}
                  </span>
                  <div className="min-w-0 flex-1">
                    <h2 className="font-display text-xl font-bold">
                      {selected.name}
                    </h2>
                    <p className="text-sm text-muted-foreground">
                      {selected.email}
                    </p>
                    <p className="mt-1 text-xs text-muted-foreground">
                      ID #{selected.id}
                      {selected.created_at
                        ? ` · tham gia ${new Date(selected.created_at).toLocaleDateString("vi-VN")}`
                        : ""}
                    </p>
                  </div>
                  <div className="flex flex-wrap gap-1.5">
                    <Badge
                      variant={
                        selected.role === ROLE_ADMIN ? "default" : "secondary"
                      }
                      className="rounded-full"
                    >
                      {selected.role === ROLE_ADMIN ? "Admin" : "User"}
                    </Badge>
                    <Badge
                      variant={
                        selected.status === 1 ? "outline" : "destructive"
                      }
                      className="rounded-full"
                    >
                      {selected.status === 1 ? "Hoạt động" : "Đã khóa"}
                    </Badge>
                    <span className={planPillClass(currentPlan)}>
                      {currentPlan}
                    </span>
                  </div>
                </div>
              </div>

              <div className="space-y-4 p-5">
                {/* Subscription */}
                <div className="hani-admin-section">
                  <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
                    <div>
                      <h3 className="text-sm font-semibold">Gói đăng ký</h3>
                      <p className="text-xs text-muted-foreground">
                        Áp dụng ngay. Nếu user vẫn báo hết lượt, bấm Reset hạn mức.
                      </p>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={busy}
                      className="gap-1.5 rounded-full"
                      onClick={() =>
                        void runAction(async () => {
                          const snap = await resetUserUsage(selected.id);
                          setUsage(snap);
                        })
                      }
                    >
                      <RotateCcw className="size-3.5" />
                      Reset hạn mức
                    </Button>
                  </div>

                  <div className="hani-admin-plan-picker">
                    {PLAN_OPTIONS.map((p) => {
                      const Icon = p.icon;
                      const active = currentPlan === p.id;
                      return (
                        <button
                          key={p.id}
                          type="button"
                          disabled={busy}
                          onClick={() => {
                            if (active) return;
                            void runAction(() =>
                              patchAdminUser(selected.id, {
                                subscription_plan: p.id,
                              }).then(() => undefined)
                            );
                          }}
                          className={cn(
                            "hani-admin-plan-option",
                            active && "hani-admin-plan-option--active"
                          )}
                        >
                          <Icon
                            className={cn(
                              "size-4",
                              active ? "text-primary" : "text-muted-foreground"
                            )}
                          />
                          <span className="font-semibold">{p.label}</span>
                          <span className="text-[0.625rem] leading-tight text-muted-foreground">
                            {p.desc}
                          </span>
                        </button>
                      );
                    })}
                  </div>

                  {usage && (
                    <p className="mt-3 rounded-lg bg-primary/5 px-3 py-2 text-xs text-muted-foreground">
                      Đã reset hôm nay — tin nhắn{" "}
                      <strong>{usage.daily_messages}</strong>
                      {usage.daily_messages_limit != null
                        ? ` / ${usage.daily_messages_limit}`
                        : ""}
                      , voice{" "}
                      <strong>
                        {Math.round(usage.daily_voice_seconds / 60)}
                      </strong>{" "}
                      phút
                    </p>
                  )}
                </div>

                {/* Account actions */}
                <div className="hani-admin-section">
                  <h3 className="mb-3 text-sm font-semibold">Tài khoản</h3>
                  <div className="hani-admin-action-grid">
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={busy}
                      className="h-10 justify-start gap-2 rounded-xl"
                      onClick={() =>
                        void runAction(() =>
                          patchAdminUser(selected.id, {
                            status: selected.status === 1 ? 0 : 1,
                          }).then(() => undefined)
                        )
                      }
                    >
                      {selected.status === 1 ? (
                        <Lock className="size-4 text-destructive" />
                      ) : (
                        <Unlock className="size-4 text-primary" />
                      )}
                      {selected.status === 1 ? "Khóa tài khoản" : "Mở khóa"}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={busy}
                      className="h-10 justify-start gap-2 rounded-xl"
                      onClick={() =>
                        void runAction(() =>
                          patchAdminUser(selected.id, {
                            role:
                              selected.role === ROLE_ADMIN ? 0 : ROLE_ADMIN,
                          }).then(() => undefined)
                        )
                      }
                    >
                      <Shield className="size-4" />
                      {selected.role === ROLE_ADMIN ? "Bỏ quyền admin" : "Làm admin"}
                    </Button>
                  </div>
                </div>

                {/* Data actions */}
                <div className="hani-admin-section">
                  <h3 className="mb-3 text-sm font-semibold">Dữ liệu chat</h3>
                  <div className="hani-admin-action-grid">
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={busy}
                      className="h-10 justify-start gap-2 rounded-xl"
                      onClick={() =>
                        void runAction(() => clearUserConversation(selected.id))
                      }
                    >
                      <MessageSquare className="size-4" />
                      Xóa hội thoại
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={busy}
                      className="h-10 justify-start gap-2 rounded-xl"
                      onClick={() =>
                        void runAction(() => clearUserMemories(selected.id))
                      }
                    >
                      <Brain className="size-4" />
                      Xóa ký ức ({memories.length})
                    </Button>
                  </div>
                </div>

                {/* Memories */}
                <div>
                  <h3 className="mb-2 text-sm font-semibold">
                    Ký ức gần đây
                  </h3>
                  <ul className="max-h-32 space-y-2 overflow-y-auto">
                    {memories.length === 0 ? (
                      <li className="text-xs text-muted-foreground">
                        Chưa có ký ức vector
                      </li>
                    ) : (
                      memories.slice(0, 12).map((m) => (
                        <li key={m.id} className="hani-admin-memory-item">
                          <Badge
                            variant="secondary"
                            className="mb-1 rounded-md text-[0.625rem]"
                          >
                            {m.memory_type || "fact"}
                          </Badge>
                          <p>{m.content}</p>
                        </li>
                      ))
                    )}
                  </ul>
                </div>

                {/* Messages */}
                <div>
                  <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
                    <h3 className="text-sm font-semibold">Tin nhắn</h3>
                    <Button
                      size="sm"
                      variant="secondary"
                      className="rounded-full text-xs"
                      onClick={() => {
                        const sid = prompt("Dán session_id (UUID):");
                        if (sid) void loadMessages(sid.trim());
                      }}
                    >
                      Tải theo session
                    </Button>
                  </div>
                  <ul className="max-h-44 space-y-2 overflow-y-auto">
                    {messages.length === 0 ? (
                      <li className="text-xs text-muted-foreground">
                        Nhập session ID để xem lịch sử tin nhắn
                      </li>
                    ) : (
                      messages.map((m) => (
                        <li
                          key={m.id}
                          className={
                            m.role === "user"
                              ? "hani-admin-msg-user"
                              : "hani-admin-msg-assistant"
                          }
                        >
                          <span className="font-semibold capitalize text-foreground/80">
                            {m.role === "user" ? "User" : "Hani"}:
                          </span>{" "}
                          {m.content}
                        </li>
                      ))
                    )}
                  </ul>
                </div>

                {/* Danger */}
                <div className="rounded-xl border border-destructive/25 bg-destructive/5 p-4">
                  <h3 className="mb-2 text-sm font-semibold text-destructive">
                    Vùng nguy hiểm
                  </h3>
                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={busy}
                    className="gap-2 rounded-xl"
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
                    <Trash2 className="size-4" />
                    Xóa user khỏi hệ thống
                  </Button>
                </div>
              </div>
            </div>
          )}
        </section>
      </div>
    </AppShell>
  );
}
