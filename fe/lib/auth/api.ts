import { API_URL } from "../config";
import type { UserGender } from "./gender";
import type { AuthResponse, AuthUser } from "./types";
import { getRefreshToken, getStoredUser, getToken, setAuth } from "./storage";

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
}

function accessTokenFromResponse(res: AuthResponse): string {
  return res.access_token ?? res.token ?? "";
}

let refreshInFlight: Promise<AuthResponse> | null = null;

async function requestRefreshSession(): Promise<AuthResponse> {
  const refresh = getRefreshToken();
  if (!refresh) throw new Error("no refresh token");

  const res = await fetch(`${API_URL}/api/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refresh }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function refreshSession(): Promise<AuthResponse> {
  if (!refreshInFlight) {
    refreshInFlight = requestRefreshSession()
      .then((res) => {
        const token = accessTokenFromResponse(res);
        const user = res.user ?? getStoredUser();
        if (user) setAuth(token, user, res.refresh_token);
        return res;
      })
      .finally(() => {
        refreshInFlight = null;
      });
  }
  return refreshInFlight;
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

export async function restoreSession(): Promise<AuthUser | null> {
  const token = getToken();
  const user = getStoredUser();
  if (!token || !user) return null;

  try {
    return await fetchMe();
  } catch {
    return null;
  }
}

export async function register(
  name: string,
  email: string,
  password: string,
  gender: UserGender
): Promise<AuthResponse> {
  const res = await fetch(`${API_URL}/api/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, password, gender }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function login(
  email: string,
  password: string
): Promise<AuthResponse> {
  const res = await fetch(`${API_URL}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

type MeResponse = {
  user?: AuthUser;
  guest?: boolean;
};

export async function fetchMe(): Promise<AuthUser> {
  const res = await fetchWithAuth(`${API_URL}/api/auth/me`);
  if (!res.ok) throw new Error(await parseError(res));
  const data = (await res.json()) as MeResponse & AuthUser;
  if (data.user) return data.user;
  if (typeof data.id === "number") return data;
  throw new Error("invalid profile response");
}

export async function uploadAvatar(file: File): Promise<AuthUser> {
  const body = new FormData();
  body.append("avatar", file);

  const res = await fetchWithAuth(`${API_URL}/api/auth/me/avatar`, {
    method: "POST",
    body,
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export type UpdateProfileInput = {
  name?: string;
  avatar?: string;
};

export async function updateProfile(
  input: UpdateProfileInput
): Promise<AuthUser> {
  const res = await fetchWithAuth(`${API_URL}/api/auth/me`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export function authHeaders(): HeadersInit {
  const token = getToken();
  if (!token) return {};
  return { Authorization: `Bearer ${token}` };
}
