"use client";

import { FormEvent, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { HaniMark } from "@/components/brand/HaniMark";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { GENDER_OPTIONS, type UserGender } from "@/lib/auth/gender";
import { cn } from "@/lib/utils";

type Mode = "login" | "register";

export function AuthForm() {
  const { login, register } = useAuth();
  const [mode, setMode] = useState<Mode>("login");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [gender, setGender] = useState<UserGender | "">("");
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
        if (!gender) {
          setError("Chọn giới tính của bạn");
          return;
        }
        await register(name, email, password, gender);
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
          <HaniMark size="lg" />
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
                onClick={() => {
                  setMode(m);
                  setError(null);
                }}
              >
                {m === "login" ? "Đăng nhập" : "Đăng ký"}
              </Button>
            ))}
          </div>

          <form className="space-y-3.5" onSubmit={onSubmit}>
            {mode === "register" && (
              <>
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
                <fieldset className="space-y-2">
                  <legend className="text-sm font-medium leading-none">
                    Giới tính của bạn
                  </legend>
                  <p className="text-xs text-muted-foreground">
                    Hani sẽ xưng hô phù hợp hơn khi trò chuyện
                  </p>
                  <div className="grid grid-cols-3 gap-2">
                    {GENDER_OPTIONS.map((opt) => (
                      <button
                        key={opt.id}
                        type="button"
                        onClick={() => setGender(opt.id)}
                        className={cn(
                          "flex flex-col items-center gap-0.5 rounded-xl border px-2 py-2.5 text-center transition-all",
                          gender === opt.id
                            ? "border-primary bg-primary/10 text-primary shadow-sm"
                            : "border-border bg-card/80 text-muted-foreground hover:border-primary/30"
                        )}
                      >
                        <span className="text-sm font-semibold">
                          {opt.label}
                        </span>
                        <span className="text-[0.625rem] leading-tight opacity-80">
                          {opt.desc}
                        </span>
                      </button>
                    ))}
                  </div>
                </fieldset>
              </>
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
