import { API_URL } from "@/lib/config";
import { authHeaders } from "@/lib/auth/api";
import type { PlanLimit, UsageSnapshot } from "./types";

export async function fetchPlans(): Promise<PlanLimit[]> {
  const res = await fetch(`${API_URL}/api/billing/plans`);
  if (!res.ok) throw new Error("Failed to load plans");
  const data = (await res.json()) as { plans: PlanLimit[] };
  return data.plans;
}

export async function fetchUsage(): Promise<UsageSnapshot> {
  const res = await fetch(`${API_URL}/api/auth/billing/usage`, {
    headers: authHeaders(),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error((data as { error?: string }).error ?? res.statusText);
  }
  return res.json();
}

export async function createGuestSession(): Promise<{
  token: string;
  refresh_token: string;
  guest_id: string;
  usage: UsageSnapshot;
}> {
  const res = await fetch(`${API_URL}/api/auth/guest`, { method: "POST" });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error((data as { error?: string }).error ?? res.statusText);
  }
  const data = await res.json();
  return {
    token: data.token ?? data.access_token,
    refresh_token: data.refresh_token,
    guest_id: data.guest_id,
    usage: data.usage,
  };
}
