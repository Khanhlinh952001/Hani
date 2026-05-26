import { API_URL } from "@/lib/config";

const TEMP_KEY_TTL_MS = 45_000;

type Cache = {
  key: string | null;
  token: string | null;
  fetchedAt: number;
};

const cache: Cache = {
  key: null,
  token: null,
  fetchedAt: 0,
};

async function fetchTemporaryApiKey(token: string): Promise<string> {
  const base = API_URL.replace(/\/$/, "");
  const res = await fetch(`${base}/api/soniox/temporary-key`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });
  const data = (await res.json().catch(() => ({}))) as {
    apiKey?: string;
    error?: string;
  };
  if (!res.ok || !data.apiKey) {
    throw new Error(data.error ?? `HTTP ${res.status}`);
  }
  return data.apiKey;
}

/** Giống Memora: cache key ngắn, lấy từ backend (không lộ SONIOX_API_KEY). */
export function createSonioxKeyFetcher(token: string) {
  return async function getSonioxApiKey(): Promise<string> {
    if (!token) {
      throw new Error("Chưa đăng nhập");
    }
    const now = Date.now();
    if (
      cache.key &&
      cache.token === token &&
      now - cache.fetchedAt < TEMP_KEY_TTL_MS
    ) {
      return cache.key;
    }
    const apiKey = await fetchTemporaryApiKey(token);
    cache.key = apiKey;
    cache.token = token;
    cache.fetchedAt = now;
    return apiKey;
  };
}
