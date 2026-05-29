"use client";

import { useAuth } from "@/hooks/useAuth";
import {
  detectDeviceType,
  FCM_TOKEN_KEY,
  isPushSupported,
  PUSH_ENABLED_KEY,
} from "@/lib/push/config";
import { heartbeatDevice, registerDevice, revokeDevice } from "@/lib/push/api";
import { fetchFCMToken, listenForegroundMessages } from "@/lib/push/firebase";
import { useCallback, useEffect, useRef } from "react";

const HEARTBEAT_MS = 5 * 60 * 1000;

export function PushRegistrar() {
  const { user, loading } = useAuth();
  const tokenRef = useRef<string | null>(null);

  const syncToken = useCallback(async () => {
    if (!user || !isPushSupported()) return;
    if (localStorage.getItem(PUSH_ENABLED_KEY) !== "true") return;

    const token = await fetchFCMToken();
    if (!token) return;

    const prev = localStorage.getItem(FCM_TOKEN_KEY);
    if (prev && prev !== token) {
      try {
        await revokeDevice(prev);
      } catch {
        /* ignore */
      }
    }

    await registerDevice(token, detectDeviceType());
    localStorage.setItem(FCM_TOKEN_KEY, token);
    tokenRef.current = token;
  }, [user]);

  useEffect(() => {
    if (loading || !user) return;
    if (localStorage.getItem(PUSH_ENABLED_KEY) !== "true") return;
    void syncToken().catch((e) => console.warn("push sync:", e));
  }, [loading, user, syncToken]);

  useEffect(() => {
    if (!user) return;
    listenForegroundMessages((_title, _body) => {
      /* optional toast — notification already shown */
    });
  }, [user]);

  useEffect(() => {
    if (!user || localStorage.getItem(PUSH_ENABLED_KEY) !== "true") return;

    const tick = () => {
      const t = tokenRef.current ?? localStorage.getItem(FCM_TOKEN_KEY);
      if (t) void heartbeatDevice(t).catch(() => {});
    };

    tick();
    const id = window.setInterval(tick, HEARTBEAT_MS);

    const onVisible = () => {
      if (document.visibilityState === "visible") tick();
    };
    document.addEventListener("visibilitychange", onVisible);

    return () => {
      clearInterval(id);
      document.removeEventListener("visibilitychange", onVisible);
    };
  }, [user]);

  return null;
}

/** Called from Settings when user toggles push on/off. */
export async function enablePushNotifications(): Promise<boolean> {
  if (!isPushSupported()) return false;
  const token = await fetchFCMToken();
  if (!token) return false;
  await registerDevice(token, detectDeviceType());
  localStorage.setItem(PUSH_ENABLED_KEY, "true");
  localStorage.setItem(FCM_TOKEN_KEY, token);
  return true;
}

export async function disablePushNotifications(): Promise<void> {
  const token = localStorage.getItem(FCM_TOKEN_KEY);
  if (token) {
    try {
      await revokeDevice(token);
    } catch {
      /* ignore */
    }
  }
  localStorage.removeItem(PUSH_ENABLED_KEY);
  localStorage.removeItem(FCM_TOKEN_KEY);
}
