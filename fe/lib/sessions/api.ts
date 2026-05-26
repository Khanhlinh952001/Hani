import { API_URL } from "../config";
import { authHeaders } from "../auth/api";

export type Session = {
  id: string;
  user_id: number;
  started_at: string;
  ended_at?: string | null;
};

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

export async function listSessions(): Promise<Session[]> {
  const res = await fetch(`${API_URL}/api/sessions`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function deleteSession(sessionId: string): Promise<void> {
  const res = await fetch(`${API_URL}/api/sessions/${sessionId}`, {
    method: "DELETE",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function clearConversationHistory(): Promise<void> {
  const res = await fetch(`${API_URL}/api/sessions/current/clear`, {
    method: "POST",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
}
