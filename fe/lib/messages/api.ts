import { API_URL } from "../config";
import { authHeaders } from "../auth/api";

export type ApiMessage = {
  id: string;
  session_id: string;
  role: string;
  content: string;
  translation_vi?: string;
  created_at: string;
};

export type MessagePage = {
  messages: ApiMessage[];
  has_more: boolean;
};

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

export async function fetchMessagePage(
  sessionId: string,
  opts: { before?: string; limit?: number } = {}
): Promise<MessagePage> {
  const params = new URLSearchParams({ session_id: sessionId });
  if (opts.before) params.set("before", opts.before);
  params.set("limit", String(opts.limit ?? 30));

  const res = await fetch(`${API_URL}/api/messages?${params}`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}
