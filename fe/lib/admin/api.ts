import { API_URL } from "../config";
import { authHeaders } from "../auth/api";

export type AdminStats = {
  users: number;
  sessions: number;
  messages: number;
  memories: number;
};

export type AdminUser = {
  id: number;
  name: string;
  email: string;
  avatar?: string;
  role: number;
  status: number;
  subscription_plan?: string;
  is_active?: boolean;
  level?: number;
  created_at: string;
  updated_at: string;
};

export type UsageSnapshot = {
  plan: string;
  daily_messages: number;
  daily_messages_limit?: number | null;
  daily_voice_seconds: number;
  daily_voice_limit?: number | null;
  warning?: boolean;
};

export type AdminSession = {
  id: string;
  user_id: number;
  started_at: string;
  ended_at?: string | null;
};

export type AdminMessage = {
  id: string;
  session_id: string;
  role: string;
  content: string;
  translation_vi?: string;
  created_at: string;
};

export type AdminMemory = {
  id: string;
  user_id: number;
  content: string;
  memory_type?: string;
  importance_score: number;
  created_at: string;
};

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

async function adminFetch(path: string, init?: RequestInit) {
  const res = await fetch(`${API_URL}/api/admin${path}`, {
    ...init,
    headers: {
      ...(await authHeaders()),
      ...(init?.headers as Record<string, string>),
    },
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res;
}

export async function fetchAdminStats(): Promise<AdminStats> {
  const res = await adminFetch("/stats");
  return res.json();
}

export async function fetchAdminUsers(): Promise<AdminUser[]> {
  const res = await adminFetch("/users");
  return res.json();
}

export async function patchAdminUser(
  id: number,
  body: {
    name?: string;
    status?: number;
    role?: number;
    subscription_plan?: string;
    is_active?: boolean;
  }
): Promise<AdminUser> {
  const res = await adminFetch(`/users/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return res.json();
}

export async function deleteAdminUser(id: number): Promise<void> {
  await adminFetch(`/users/${id}`, { method: "DELETE" });
}

export async function fetchUserMemories(userId: number): Promise<AdminMemory[]> {
  const res = await adminFetch(`/users/${userId}/memories`);
  return res.json();
}

export async function fetchSessionMessages(
  sessionId: string
): Promise<AdminMessage[]> {
  const res = await adminFetch(`/sessions/${sessionId}/messages`);
  return res.json();
}

export async function clearUserMemories(userId: number): Promise<void> {
  await adminFetch(`/users/${userId}/memories`, { method: "DELETE" });
}

export async function clearUserConversation(userId: number): Promise<void> {
  await adminFetch(`/users/${userId}/clear-conversation`, { method: "POST" });
}

export async function resetUserUsage(userId: number): Promise<UsageSnapshot> {
  const res = await adminFetch(`/users/${userId}/reset-usage`, {
    method: "POST",
  });
  return res.json();
}
