import { initializeApp, getApps, type FirebaseApp } from "firebase/app";
import { getAnalytics, isSupported, type Analytics } from "firebase/analytics";
import { firebaseConfig } from "@/lib/push/config";

let app: FirebaseApp | undefined;
let analytics: Analytics | undefined;

/** Initialize Firebase app (singleton). */
export function getFirebaseApp(): FirebaseApp {
  if (!app) {
    app = getApps().length ? getApps()[0]! : initializeApp(firebaseConfig);
  }
  return app;
}

/** Firebase Analytics — browser only. */
export async function getFirebaseAnalytics(): Promise<Analytics | null> {
  if (typeof window === "undefined") return null;
  if (analytics) return analytics;

  const supported = await isSupported();
  if (!supported) return null;

  analytics = getAnalytics(getFirebaseApp());
  return analytics;
}
