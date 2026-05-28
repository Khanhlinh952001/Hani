const STREAK_KEY = "hani_streak";
const STREAK_DATE_KEY = "hani_streak_date";
const VOICE_MIN_KEY = "hani_voice_minutes";
const VOICE_DATE_KEY = "hani_voice_date";
const PREVIEW_KEY = "hani_last_preview";
const UNREAD_KEY = "hani_chat_unread";
const PREVIEW_AT_KEY = "hani_last_preview_at";

export type HomePreview = {
  text: string;
  at: string;
};

const DAILY_PHRASES = [
  { ko: "오늘의 표현", phrase: "잘 자요", vi: "Chúc ngủ ngon" },
  { ko: "오늘의 표현", phrase: "보고 싶어요", vi: "Em nhớ anh" },
  { ko: "오늘의 표현", phrase: "맛있게 드세요", vi: "Ăn ngon nhé" },
  { ko: "오늘의 표현", phrase: "힘내요", vi: "Cố lên nhé" },
  { ko: "오늘의 표현", phrase: "천천히 말해요", vi: "Nói chậm thôi" },
] as const;

const HANI_MOODS = [
  { ko: "기분 좋아요", vi: "Vui vẻ", emoji: "🌸" },
  { ko: "조금 졸려요", vi: "Buồn ngủ", emoji: "🌙" },
  { ko: "너 생각 중", vi: "Đang nghĩ bạn", emoji: "💭" },
  { ko: "설레요", vi: "Hồi hộp", emoji: "💕" },
] as const;

const PROACTIVE_LINES = [
  { ko: "오늘도 같이 연습해볼까요?", vi: "Hôm nay luyện cùng nhau nhé?" },
  { ko: "지금 기다리고 있어요…", vi: "Hani đang đợi bạn…" },
  { ko: "메시지 보내줄래요?", vi: "Nhắn cho Hani đi" },
  { ko: "목소리도 듣고 싶어요", vi: "Muốn nghe giọng bạn" },
] as const;

function todayKey() {
  return new Date().toISOString().slice(0, 10);
}

export function getDynamicSubtitle(): { ko: string; vi: string } {
  const h = new Date().getHours();
  if (h < 6) {
    return { ko: "밤이 깊었네요…", vi: "Đêm khuya rồi nhé" };
  }
  if (h < 12) {
    return { ko: "좋은 아침이에요", vi: "Chào buổi sáng" };
  }
  if (h < 18) {
    return { ko: "오늘도 같이 연습해볼까요?", vi: "Hôm nay luyện cùng nhau nhé?" };
  }
  return { ko: "지금 기다리고 있어요…", vi: "Hani đang đợi bạn…" };
}

export function getDailyPhrase() {
  const day = Math.floor(Date.now() / 86_400_000);
  return DAILY_PHRASES[day % DAILY_PHRASES.length];
}

export function getHaniMood() {
  const slot = Math.floor(Date.now() / 3_600_000);
  return HANI_MOODS[slot % HANI_MOODS.length];
}

export function getProactiveLine() {
  const slot = Math.floor(Date.now() / 7_200_000);
  return PROACTIVE_LINES[slot % PROACTIVE_LINES.length];
}

export function readStreak(): number {
  if (typeof window === "undefined") return 0;
  try {
    const n = parseInt(localStorage.getItem(STREAK_KEY) ?? "0", 10);
    return Number.isFinite(n) ? n : 0;
  } catch {
    return 0;
  }
}

export function touchStreak(): number {
  if (typeof window === "undefined") return 0;
  const today = todayKey();
  const last = localStorage.getItem(STREAK_DATE_KEY);
  let streak = readStreak();
  if (last === today) return streak;
  if (last) {
    const prev = new Date(last);
    const now = new Date(today);
    const diff = (now.getTime() - prev.getTime()) / 86_400_000;
    streak = diff === 1 ? streak + 1 : 1;
  } else {
    streak = 1;
  }
  localStorage.setItem(STREAK_KEY, String(streak));
  localStorage.setItem(STREAK_DATE_KEY, today);
  return streak;
}

export function readVoiceMinutesToday(): number {
  if (typeof window === "undefined") return 0;
  try {
    if (localStorage.getItem(VOICE_DATE_KEY) !== todayKey()) return 0;
    const n = parseInt(localStorage.getItem(VOICE_MIN_KEY) ?? "0", 10);
    return Number.isFinite(n) ? n : 0;
  } catch {
    return 0;
  }
}

export function addVoiceMinutes(delta: number) {
  if (typeof window === "undefined" || delta <= 0) return;
  const today = todayKey();
  if (localStorage.getItem(VOICE_DATE_KEY) !== today) {
    localStorage.setItem(VOICE_MIN_KEY, "0");
    localStorage.setItem(VOICE_DATE_KEY, today);
  }
  const next = readVoiceMinutesToday() + delta;
  localStorage.setItem(VOICE_MIN_KEY, String(next));
}

export function readLastPreview(): HomePreview | null {
  if (typeof window === "undefined") return null;
  try {
    const text = localStorage.getItem(PREVIEW_KEY);
    const at = localStorage.getItem(PREVIEW_AT_KEY);
    if (!text) return null;
    return { text, at: at ?? "" };
  } catch {
    return null;
  }
}

export function saveLastPreview(text: string) {
  if (typeof window === "undefined" || !text.trim()) return;
  const trimmed = text.trim().slice(0, 80);
  localStorage.setItem(PREVIEW_KEY, trimmed);
  localStorage.setItem(PREVIEW_AT_KEY, new Date().toISOString());
  localStorage.setItem(UNREAD_KEY, "1");
}

export function readUnread(): boolean {
  if (typeof window === "undefined") return false;
  return localStorage.getItem(UNREAD_KEY) === "1";
}

export function clearUnread() {
  if (typeof window === "undefined") return;
  localStorage.removeItem(UNREAD_KEY);
}

export const DEFAULT_CHAT_PREVIEW = "잘 잤어요? ☀️";
