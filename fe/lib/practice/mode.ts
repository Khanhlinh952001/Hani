/** speak = Soniox push-to-talk + AI TTS; chat = text only */
export type PracticeMode = "speak" | "chat";

export const PRACTICE_MODE_OPTIONS: {
  id: PracticeMode;
  href: string;
  label: string;
  desc: string;
  ko: string;
}[] = [
  {
    id: "speak",
    href: "/speak",
    label: "Luyện nói",
    desc: "Giữ nút nói — Hani trả lời bằng giọng",
    ko: "말하기 연습",
  },
  {
    id: "chat",
    href: "/chat",
    label: "Luyện nhắn tin",
    desc: "Chỉ gõ chữ — không đọc, không giọng nói",
    ko: "채팅 연습",
  },
];
