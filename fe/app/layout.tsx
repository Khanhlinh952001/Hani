import type { Metadata, Viewport } from "next";
import { Gowun_Batang, Noto_Sans_KR } from "next/font/google";
import { Providers } from "./providers";
import "./globals.css";

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#ff5c8a" },
    { media: "(prefers-color-scheme: dark)", color: "#ff5c8a" },
  ],
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  viewportFit: "cover",
  colorScheme: "light",
};

const noto = Noto_Sans_KR({
  variable: "--font-noto",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

const gowun = Gowun_Batang({
  variable: "--font-gowun",
  subsets: ["latin"],
  weight: ["400", "700"],
});

export const metadata: Metadata = {
  title: "Hani — 네 연인",
  description: "Trò chuyện tự nhiên với Hani — người yêu Hàn Quốc của bạn",
  applicationName: "Hani",
  appleWebApp: {
    capable: true,
    statusBarStyle: "default",
    title: "Hani",
  },
  formatDetection: {
    telephone: false,
  },
  icons: {
    icon: [
      { url: "/logo/favicon.ico", sizes: "48x48" },
      { url: "/logo/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "/logo/favicon-32x32.png", sizes: "32x32", type: "image/png" },
    ],
    apple: [{ url: "/logo/apple-touch-icon.png", sizes: "180x180" }],
  },
  other: {
    "mobile-web-app-capable": "yes",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ko" className={noto.variable}>
      <body className={`${noto.variable} ${gowun.variable} antialiased`}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
