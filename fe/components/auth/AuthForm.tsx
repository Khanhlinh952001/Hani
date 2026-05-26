"use client";

import { FormEvent, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { HaniMark } from "@/components/brand/HaniMark";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

type Mode = "login" | "register";

export function AuthForm() {
  const { login, register } = useAuth();
  const [mode, setMode] = useState<Mode>("login");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      if (mode === "login") {
        await login(email, password);
      } else {
        await register(name, email, password);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="hani-auth-wrap">
      <Card className="hani-auth-card relative gap-0 pt-1">
        <div className="hani-auth-header">
          <HaniMark size="lg" className="hani-avatar-glow" />
          <h1 className="font-display text-3xl font-bold tracking-tight text-foreground">
            Hani
          </h1>
          <p className="text-sm text-muted-foreground">
            네 연인 · người yêu Hàn của bạn
          </p>
        </div>

        <CardContent className="space-y-4 px-6 pb-8 pt-2">
          <div
            className="grid grid-cols-2 gap-1 rounded-xl bg-muted/70 p-1"
            role="tablist"
          >
            {(["login", "register"] as const).map((m) => (
              <Button
                key={m}
                type="button"
                variant="ghost"
                size="sm"
                className={cn(
                  "h-10 w-full rounded-lg font-medium shadow-none",
                  mode === m
                    ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground"
                    : "bg-transparent text-muted-foreground hover:bg-background/90 hover:text-foreground"
                )}
                onClick={() => setMode(m)}
              >
                {m === "login" ? "Đăng nhập" : "Đăng ký"}
              </Button>
            ))}
          </div>

          <form className="space-y-3.5" onSubmit={onSubmit}>
            {mode === "register" && (
              <div className="space-y-1.5">
                <Label htmlFor="name">Anh muốn em gọi anh là gì?</Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Minh, 오빠…"
                  required
                  minLength={2}
                />
              </div>
            )}
            <div className="space-y-1.5">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@email.com"
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="password">Mật khẩu</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••"
                required
                minLength={6}
              />
            </div>

            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <Button
              type="submit"
              className="mt-1 h-11 w-full shadow-md shadow-primary/25"
              size="lg"
              disabled={submitting}
            >
              {submitting
                ? "…"
                : mode === "login"
                  ? "Vào trò chuyện"
                  : "Tạo tài khoản"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
