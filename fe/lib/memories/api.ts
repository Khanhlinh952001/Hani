import { API_URL } from "../config";
import { authHeaders } from "../auth/api";

export type Memory = {
  id: string;
  user_id: number;
  content: string;
  translation_vi?: string;
  memory_type?: string;
  importance_score: number;
  created_at: string;
};

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

/** All memories for the logged-in user (newest / most important first). */
export async function fetchMemories(memoryType?: string): Promise<Memory[]> {
  const params = memoryType ? `?memory_type=${encodeURIComponent(memoryType)}` : "";
  const res = await fetch(`${API_URL}/api/memories${params}`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  const list = (await res.json()) as Memory[];
  return Array.isArray(list) ? list : [];
}
