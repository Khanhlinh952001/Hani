import type { UserGender } from "./gender";

export type AuthUser = {
  id: number;
  name: string;
  email: string;
  gender?: UserGender | string;
  ai_profile_id?: string;
  selected_character_id?: string;
  avatar?: string;
  level?: number;
  role?: number;
  status?: number;
  updated_at?: string;
  created_at?: string;
};

export type AuthResponse = {
  token: string;
  user: AuthUser;
};

export const ROLE_ADMIN = 1;

export function isAdmin(user: AuthUser | null | undefined): boolean {
  return user?.role === ROLE_ADMIN;
}
