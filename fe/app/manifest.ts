import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    id: "/",
    name: "Hani — 네 연인",
    short_name: "Hani",
    description:
      "Trò chuyện với người yêu AI Hàn — học tiếng Hàn qua cảm xúc thật",
    start_url: "/",
    scope: "/",
    display: "standalone",
    orientation: "portrait",
    background_color: "#fff7fa",
    theme_color: "#ff5c8a",
    lang: "ko",
    dir: "ltr",
    categories: ["social", "lifestyle", "education"],
    icons: [
      {
        src: "/logo/favicon-16x16.png",
        sizes: "16x16",
        type: "image/png",
      },
      {
        src: "/logo/favicon-32x32.png",
        sizes: "32x32",
        type: "image/png",
      },
      {
        src: "/logo/android-chrome-192x192.png",
        sizes: "192x192",
        type: "image/png",
        purpose: "any",
      },
      {
        src: "/logo/android-chrome-512x512.png",
        sizes: "512x512",
        type: "image/png",
        purpose: "any",
      },
      {
        src: "/logo/android-chrome-512x512.png",
        sizes: "512x512",
        type: "image/png",
        purpose: "maskable",
      },
      {
        src: "/logo/apple-touch-icon.png",
        sizes: "180x180",
        type: "image/png",
        purpose: "any",
      },
    ],
  };
}
