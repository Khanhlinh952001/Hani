/* Hani PWA + FCM push — cache + background notifications */
importScripts("https://www.gstatic.com/firebasejs/10.14.1/firebase-app-compat.js");
importScripts("https://www.gstatic.com/firebasejs/10.14.1/firebase-messaging-compat.js");

firebase.initializeApp({
  apiKey: "AIzaSyAumWf1YVCdIYhX_YVCzdqkqvMr1iURrt0",
  authDomain: "hani-ad3e8.firebaseapp.com",
  projectId: "hani-ad3e8",
  storageBucket: "hani-ad3e8.firebasestorage.app",
  messagingSenderId: "647699315918",
  appId: "1:647699315918:web:1bf9a7a4bd718dddcd249f",
});

try {
  const messaging = firebase.messaging();
  messaging.onBackgroundMessage((payload) => {
    const title = payload.notification?.title || "Hani";
    const body = payload.notification?.body || "";
    self.registration.showNotification(title, {
      body,
      icon: "/logo/android-chrome-192x192.png",
      badge: "/logo/android-chrome-192x192.png",
      data: { url: payload.fcmOptions?.link || "/" },
    });
  });
} catch (_) {
  /* FCM unavailable in this context */
}

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const url = event.notification.data?.url || "/";
  event.waitUntil(
    clients.matchAll({ type: "window", includeUncontrolled: true }).then((list) => {
      for (const client of list) {
        if ("focus" in client) {
          client.navigate(url);
          return client.focus();
        }
      }
      return clients.openWindow(url);
    })
  );
});

const CACHE_VERSION = "hani-pwa-v4";
const CACHE = CACHE_VERSION;
const PRECACHE = [
  "/offline.html",
  "/logo/android-chrome-192x192.png",
  "/logo/android-chrome-512x512.png",
  "/logo/apple-touch-icon.png",
  "/logo/favicon-32x32.png",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE).then((cache) => cache.addAll(PRECACHE)).then(() => self.skipWaiting())
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)))
      )
      .then(() => self.clients.claim())
  );
});

self.addEventListener("fetch", (event) => {
  const { request } = event;
  if (request.method !== "GET") return;

  const url = new URL(request.url);
  if (url.origin !== self.location.origin) return;

  if (PRECACHE.includes(url.pathname)) {
    event.respondWith(
      caches.match(request).then((cached) => cached || fetch(request))
    );
    return;
  }

  if (request.mode === "navigate") {
    event.respondWith(
      fetch(request).catch(() =>
        caches.match("/offline.html").then((r) => r || new Response("Offline", { status: 503 }))
      )
    );
    return;
  }

  if (url.pathname.startsWith("/_next/static/")) {
    event.respondWith(
      caches.open(CACHE).then(async (cache) => {
        const hit = await cache.match(request);
        const network = fetch(request).then((res) => {
          if (res.ok) cache.put(request, res.clone());
          return res;
        });
        return hit || network;
      })
    );
  }
});
