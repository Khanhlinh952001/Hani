import {
  getMessaging,
  getToken,
  onMessage,
  type Messaging,
} from "firebase/messaging";
import { getFirebaseApp } from "@/lib/firebase/app";
import { FCM_VAPID_KEY } from "./config";

let messaging: Messaging | undefined;

export async function getFirebaseMessaging(): Promise<Messaging | null> {
  if (typeof window === "undefined") return null;
  if (!("serviceWorker" in navigator)) return null;
  if (!messaging) {
    messaging = getMessaging(getFirebaseApp());
  }
  return messaging;
}

export async function registerMessagingServiceWorker(): Promise<ServiceWorkerRegistration> {
  let reg = await navigator.serviceWorker.getRegistration("/");
  if (!reg) {
    reg = await navigator.serviceWorker.register("/sw.js", { scope: "/" });
  }
  await navigator.serviceWorker.ready;
  return reg;
}

export async function fetchFCMToken(): Promise<string | null> {
  if (!FCM_VAPID_KEY) {
    console.warn("push: set NEXT_PUBLIC_FIREBASE_VAPID_KEY in .env.local");
    return null;
  }

  const permission = await Notification.requestPermission();
  if (permission !== "granted") return null;

  const swReg = await registerMessagingServiceWorker();
  const msg = await getFirebaseMessaging();
  if (!msg) return null;

  return getToken(msg, { vapidKey: FCM_VAPID_KEY, serviceWorkerRegistration: swReg });
}

export function listenForegroundMessages(
  onPayload: (title: string, body: string) => void
): (() => void) | undefined {
  void getFirebaseMessaging().then((msg) => {
    if (!msg) return;
    onMessage(msg, (payload) => {
      const title = payload.notification?.title ?? "Hani";
      const body = payload.notification?.body ?? "";
      onPayload(title, body);
      if (Notification.permission === "granted") {
        new Notification(title, {
          body,
          icon: "/logo/android-chrome-192x192.png",
        });
      }
    });
  });
  return undefined;
}
