import { API_URL } from "@/lib/config";

/** Resolve avatar path/URL for display (uploads, absolute, blob). */
export function resolveAvatarSrc(
  avatar?: string,
  cacheBust?: string
): string | undefined {
  const raw = avatar?.trim();
  if (!raw) return undefined;

  let url: string;
  if (
    raw.startsWith("http://") ||
    raw.startsWith("https://") ||
    raw.startsWith("blob:")
  ) {
    url = raw;
  } else if (raw.startsWith("/uploads/")) {
    url = `${API_URL}${raw}`;
  } else if (raw.startsWith("/")) {
    url = raw;
  } else {
    url = raw;
  }

  if (cacheBust && url.includes("/uploads/")) {
    const sep = url.includes("?") ? "&" : "?";
    return `${url}${sep}v=${encodeURIComponent(cacheBust)}`;
  }
  return url;
}
