"use client";

import Link from "next/link";
import { useCallback, useState } from "react";
import { ArrowLeft, Trash2, Volume2 } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import { AppShell } from "@/components/layout/AppShell";
import { HaniMark } from "@/components/brand/HaniMark";
import { AvatarUpload } from "@/components/settings/AvatarUpload";
import {
  SONIOX_VOICE_OPTIONS,
  TTS_LANGUAGE_OPTIONS,
  TTS_PROVIDER_OPTIONS,
  TTS_VOICE_OPTIONS,
} from "@/lib/settings/types";
import { clearConversationHistory } from "@/lib/sessions/api";
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
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";

const SESSION_KEY = "hani_session_id";

export function SettingsView() {
  const { user, logout } = useAuth();
  const {
    showVietnamese,
    setShowVietnamese,
    ttsProvider,
    setTtsProvider,
    ttsVoice,
    setTtsVoice,
    ttsLanguage,
    setTtsLanguage,
  } = useSettings();
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const isSoniox = ttsProvider === "soniox";
  const voiceOptions = isSoniox ? SONIOX_VOICE_OPTIONS : TTS_VOICE_OPTIONS;

  const handleClearHistory = useCallback(async () => {
    if (
      !confirm(
        "Xóa toàn bộ lịch sử và ký ức của Hani? Cuộc trò chuyện và vector nhớ sẽ bị xóa hết."
      )
    ) {
      return;
    }

    setBusy(true);
    setError(null);
    setMessage(null);
    try {
      await clearConversationHistory();
      localStorage.removeItem(SESSION_KEY);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Xóa thất bại");
    } finally {
      setBusy(false);
    }
  }, []);

  return (
    <AppShell>
      <header className="hani-header flex items-center gap-3 px-4 py-3">
        <Button variant="ghost" size="icon-sm" asChild>
          <Link href="/" aria-label="Về trang chủ">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <HaniMark size="sm" />
        <div>
          <h1 className="font-display text-lg font-bold">Cài đặt</h1>
          <p className="text-xs text-muted-foreground">{user?.name}</p>
        </div>
      </header>

      <main className="flex-1 space-y-4 overflow-y-auto p-4 pb-8">
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        {message && (
          <Alert>
            <AlertDescription>{message}</AlertDescription>
          </Alert>
        )}

        <Alert>
          <Volume2 className="size-4" />
          <AlertDescription>
            Đổi <strong>TTS</strong> hoặc giọng xong hãy{" "}
            <strong>thoát và mở lại</strong> trang{" "}
            <Link href="/speak" className="text-primary underline-offset-2 hover:underline">
              luyện nói
            </Link>{" "}
            để kết nối lại.
          </AlertDescription>
        </Alert>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Ảnh đại diện</CardTitle>
            <CardDescription>Hiển thị trên tin nhắn và thanh trên cùng</CardDescription>
          </CardHeader>
          <CardContent>
            <AvatarUpload
              onSuccess={(msg) => {
                setMessage(msg);
                setError(null);
              }}
              onError={(msg) => {
                setError(msg);
                setMessage(null);
              }}
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Hiển thị</CardTitle>
            <CardDescription>Dòng tiếng Việt dưới câu Hàn</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="vi-toggle" className="flex-1 cursor-pointer">
                Tiếng Việt bên dưới
              </Label>
              <Switch
                id="vi-toggle"
                checked={showVietnamese}
                onCheckedChange={setShowVietnamese}
              />
            </div>
          </CardContent>
        </Card>

        <Card className="border-primary/20">
          <CardHeader>
            <div className="flex items-center justify-between gap-2">
              <CardTitle className="text-base">Giọng Hani (TTS)</CardTitle>
              <Badge variant={isSoniox ? "secondary" : "default"}>
                {isSoniox ? "Soniox" : "OpenAI"}
              </Badge>
            </div>
            <CardDescription>
              Nhận dạng giọng nói dùng Soniox STT (khóa tạm từ server, nhấn giữ để nói).
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="tts-provider">Nhà cung cấp giọng</Label>
              <Select
                value={ttsProvider}
                onValueChange={(v) => setTtsProvider(v as typeof ttsProvider)}
              >
                <SelectTrigger id="tts-provider" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TTS_PROVIDER_OPTIONS.map((p) => (
                    <SelectItem key={p.id} value={p.id}>
                      {p.label} — {p.desc}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="tts-voice">Giọng đọc</Label>
              <Select value={ttsVoice} onValueChange={setTtsVoice}>
                <SelectTrigger id="tts-voice" className="w-full">
                  <SelectValue placeholder="Chọn giọng" />
                </SelectTrigger>
                <SelectContent>
                  {voiceOptions.map((v) => (
                    <SelectItem key={v.id} value={v.id}>
                      {v.label} — {v.desc}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {isSoniox ? (
              <div className="space-y-2">
                <Label htmlFor="tts-lang">Ngôn ngữ đọc (Soniox)</Label>
                <Select
                  value={ttsLanguage}
                  onValueChange={(v) =>
                    setTtsLanguage(v as typeof ttsLanguage)
                  }
                >
                  <SelectTrigger id="tts-lang" className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {TTS_LANGUAGE_OPTIONS.map((l) => (
                      <SelectItem key={l.id} value={l.id}>
                        {l.label} — {l.desc}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            ) : (
              <p className="text-xs text-muted-foreground">
                OpenAI TTS tự nhận tiếng Hàn từ nội dung câu.
              </p>
            )}
          </CardContent>
        </Card>

        <Card className="border-destructive/30">
          <CardHeader>
            <CardTitle className="text-base">Cuộc trò chuyện</CardTitle>
            <CardDescription>
              Một cuộc chuyện chung. Xóa cả tin nhắn và ký ức (vector).
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button
              variant="destructive"
              className="w-full h-10"
              disabled={busy}
              onClick={() => void handleClearHistory()}
            >
              <Trash2 className="size-4" />
              {busy ? "Đang xóa…" : "Xóa toàn bộ lịch sử"}
            </Button>
          </CardContent>
        </Card>

        <Separator />

        <Button variant="outline" className="w-full h-10" onClick={() => void logout()}>
          Đăng xuất
        </Button>
      </main>
    </AppShell>
  );
}
