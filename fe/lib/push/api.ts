import { API_URL } from "../config";
import { authHeaders, refreshSession } from "../auth/api";
import { getRefreshToken, getToken } from "../auth/storage";

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

async function fetchWithAuth(
  url: string,
  init: RequestInit = {}
): Promise<Response> {
  const token = getToken();
  if (!token) throw new Error("not logged in");

  const headers = new Headers(init.headers);
  headers.set("Authorization", `Bearer ${token}`);

  let res = await fetch(url, { ...init, headers });
  if (res.status === 401 && getRefreshToken()) {
    await refreshSession();
    const newToken = getToken();
    if (!newToken) throw new Error("not logged in");
    headers.set("Authorization", `Bearer ${newToken}`);
    res = await fetch(url, { ...init, headers });
  }
  return res;
}

export async function registerDevice(
  fcmToken: string,
  deviceType: "android" | "ios" | "web"
): Promise<void> {
  const res = await fetchWithAuth(`${API_URL}/api/devices`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify({ fcm_token: fcmToken, device_type: deviceType }),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function heartbeatDevice(fcmToken: string): Promise<void> {
  const res = await fetchWithAuth(`${API_URL}/api/devices/heartbeat`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify({ fcm_token: fcmToken }),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function revokeDevice(fcmToken: string): Promise<void> {
  const res = await fetchWithAuth(
    `${API_URL}/api/devices/${encodeURIComponent(fcmToken)}`,
    { method: "DELETE", headers: authHeaders() }
  );
  if (!res.ok && res.status !== 404) throw new Error(await parseError(res));
}

export async function sendTestPush(): Promise<{ title: string; body: string }> {
  const res = await fetchWithAuth(`${API_URL}/api/push/test`, {
    method: "POST",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}
