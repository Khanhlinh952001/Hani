import { API_URL } from "../config";
import type { UserGender } from "./gender";
import type { AuthResponse, AuthUser } from "./types";
import { getToken } from "./storage";

async function parseError(res: Response): Promise<string> {
  const data = await res.json().catch(() => ({}));
  return (data as { error?: string }).error ?? res.statusText;
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

export async function fetchMe(): Promise<AuthUser> {
  const token = getToken();
  if (!token) throw new Error("not logged in");

  const res = await fetch(`${API_URL}/api/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function uploadAvatar(file: File): Promise<AuthUser> {
  const token = getToken();
  if (!token) throw new Error("not logged in");

  const body = new FormData();
  body.append("avatar", file);

  const res = await fetch(`${API_URL}/api/auth/me/avatar`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
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
  const token = getToken();
  if (!token) throw new Error("not logged in");

  const res = await fetch(`${API_URL}/api/auth/me`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
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
