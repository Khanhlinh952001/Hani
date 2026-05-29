/** Firebase web config — public client values (override via NEXT_PUBLIC_FIREBASE_*). */
export const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY ?? "AIzaSyAumWf1YVCdIYhX_YVCzdqkqvMr1iURrt0",
  authDomain:
    process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN ?? "hani-ad3e8.firebaseapp.com",
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID ?? "hani-ad3e8",
  storageBucket:
    process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET ??
    "hani-ad3e8.firebasestorage.app",
  messagingSenderId:
    process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID ?? "647699315918",
  appId:
    process.env.NEXT_PUBLIC_FIREBASE_APP_ID ??
    "1:647699315918:web:1bf9a7a4bd718dddcd249f",
  measurementId:
    process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID ?? "G-25G5NLXJKQ",
};

/** Web Push VAPID key — Firebase Console → Cloud Messaging → Web Push certificates */
export const FCM_VAPID_KEY = process.env.NEXT_PUBLIC_FIREBASE_VAPID_KEY ?? "";

export const PUSH_ENABLED_KEY = "hani_push_enabled";
export const FCM_TOKEN_KEY = "hani_fcm_token";

export function isPushSupported(): boolean {
  if (typeof window === "undefined") return false;
  return "Notification" in window && "serviceWorker" in navigator;
}

export function detectDeviceType(): "android" | "ios" | "web" {
  if (typeof navigator === "undefined") return "web";
  const ua = navigator.userAgent;
  if (/iPhone|iPad|iPod/i.test(ua)) return "ios";
  if (/Android/i.test(ua)) return "android";
  return "web";
}

export function isStandalonePWA(): boolean {
  if (typeof window === "undefined") return false;
  return (
    window.matchMedia("(display-mode: standalone)").matches ||
    // @ts-expect-error iOS Safari
    window.navigator.standalone === true
  );
}
