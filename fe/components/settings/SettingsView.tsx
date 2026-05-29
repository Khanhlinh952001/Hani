"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { ArrowLeft, Bell, Trash2, Volume2 } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import { useToast } from "@/hooks/useToast";
import {
  disablePushNotifications,
  enablePushNotifications,
} from "@/components/pwa/PushRegistrar";
import {
  detectDeviceType,
  isPushSupported,
  isStandalonePWA,
  PUSH_ENABLED_KEY,
} from "@/lib/push/config";
import { CompanionLayout } from "@/components/layout/CompanionLayout";
import { HaniMark } from "@/components/brand/HaniMark";
import { AvatarUpload } from "@/components/settings/AvatarUpload";
import { SONIOX_VOICE_OPTIONS, TTS_LANGUAGE_OPTIONS } from "@/lib/settings/types";
import { fetchUsage } from "@/lib/billing/api";
import type { UsageSnapshot } from "@/lib/billing/types";
import { clearConversationHistory } from "@/lib/sessions/api";
import { sendTestPush } from "@/lib/push/api";
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
    ttsVoice,
    setTtsVoice,
    ttsLanguage,
    setTtsLanguage,
  } = useSettings();
  const { toast } = useToast();
  const [busy, setBusy] = useState(false);
  const [usage, setUsage] = useState<UsageSnapshot | null>(null);
  const [pushEnabled, setPushEnabled] = useState(false);
  const [pushBusy, setPushBusy] = useState(false);

  useEffect(() => {
    setPushEnabled(localStorage.getItem(PUSH_ENABLED_KEY) === "true");
  }, []);

  useEffect(() => {
    fetchUsage()
      .then(setUsage)
      .catch(() => setUsage(null));
  }, []);

  const voiceOptions = SONIOX_VOICE_OPTIONS;

  const handlePushToggle = useCallback(async (on: boolean) => {
    setPushBusy(true);
    try {
      if (on) {
        const ok = await enablePushNotifications();
        if (!ok) {
          toast("Không bật được thông báo — hãy cho phép trong trình duyệt.", "destructive");
          setPushEnabled(false);
          return;
        }
        setPushEnabled(true);
        toast("Hani sẽ nhắn khi nhớ anh 💕");
      } else {
        await disablePushNotifications();
        setPushEnabled(false);
        toast("Đã tắt thông báo");
      }
    } catch (e) {
      toast(e instanceof Error ? e.message : "Thông báo thất bại", "destructive");
      setPushEnabled(false);
    } finally {
      setPushBusy(false);
    }
  }, [toast]);

  const handleTestPush = useCallback(async () => {
    setPushBusy(true);
    try {
      const res = await sendTestPush();
      toast(`Đã gửi: "${res.title}" — ${res.body}`);
    } catch (e) {
      toast(e instanceof Error ? e.message : "Gửi test thất bại", "destructive");
    } finally {
      setPushBusy(false);
    }
  }, [toast]);

  const handleClearHistory = useCallback(async () => {
    if (
      !confirm(
        "Xóa toàn bộ lịch sử và ký ức của Hani? Cuộc trò chuyện và vector nhớ sẽ bị xóa hết."
      )
    ) {
      return;
    }

    setBusy(true);
    try {
      await clearConversationHistory();
      localStorage.removeItem(SESSION_KEY);
      toast("Đã xóa lịch sử trò chuyện");
    } catch (e) {
      toast(e instanceof Error ? e.message : "Xóa thất bại", "destructive");
    } finally {
      setBusy(false);
    }
  }, [toast]);

  return (
    <CompanionLayout>
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
        {usage && (
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">Gói & hạn mức hôm nay</CardTitle>
              <CardDescription>
                Gói{" "}
                <Badge variant="secondary" className="uppercase">
                  {usage.plan}
                </Badge>
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <p>
                Tin nhắn:{" "}
                <strong>{usage.daily_messages}</strong>
                {usage.daily_messages_limit != null
                  ? ` / ${usage.daily_messages_limit}`
                  : " / không giới hạn"}
              </p>
              <p>
                Giọng nói:{" "}
                <strong>{Math.round(usage.daily_voice_seconds / 60)}</strong> phút
                {usage.daily_voice_limit != null
                  ? ` / ${Math.round(usage.daily_voice_limit / 60)} phút`
                  : " / không giới hạn"}
              </p>
              {usage.warning && (
                <p className="text-primary">
                  Sắp hết lượt — cân nhắc nâng cấp Plus/Premium 💕
                </p>
              )}
            </CardContent>
          </Card>
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
              onSuccess={(msg) => toast(msg)}
              onError={(msg) => toast(msg, "destructive")}
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between gap-2">
              <CardTitle className="text-base">Thông báo</CardTitle>
              <Bell className="size-4 text-primary" />
            </div>
            <CardDescription>
              Hani nhắn khi lâu không thấy anh online
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {!isPushSupported() ? (
              <p className="text-sm text-muted-foreground">
                Trình duyệt này không hỗ trợ push notification.
              </p>
            ) : (
              <>
                <div className="flex items-center justify-between gap-4">
                  <Label htmlFor="push-toggle" className="flex-1 cursor-pointer">
                    Nhận thông báo từ Hani
                  </Label>
                  <Switch
                    id="push-toggle"
                    checked={pushEnabled}
                    disabled={pushBusy}
                    onCheckedChange={(v) => void handlePushToggle(v)}
                  />
                </div>
                {detectDeviceType() === "ios" && !isStandalonePWA() && (
                  <p className="text-xs text-muted-foreground">
                    iPhone: mở bằng Safari → <strong>Chia sẻ</strong> →{" "}
                    <strong>Thêm vào Màn hình chính</strong>, rồi mở app từ icon
                    ngoài màn hình để nhận push.
                  </p>
                )}
                {pushEnabled && (
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    disabled={pushBusy}
                    onClick={() => void handleTestPush()}
                  >
                    {pushBusy ? "Đang gửi…" : "Gửi thử thông báo"}
                  </Button>
                )}
              </>
            )}
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
              <Badge variant="secondary">Soniox</Badge>
            </div>
            <CardDescription>
              Nhận dạng giọng nói dùng Soniox STT (khóa tạm từ server, nhấn giữ để nói).
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-xs text-muted-foreground">
              Giọng đọc qua <strong>Soniox TTS</strong> (cần SONIOX_API_KEY trên server).
            </p>

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

            <div className="space-y-2">
                <Label htmlFor="tts-lang">Ngôn ngữ đọc</Label>
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
    </CompanionLayout>
  );
}
