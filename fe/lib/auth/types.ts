import type { UserGender } from "./gender";

import type { UsageSnapshot } from "@/lib/billing/types";

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
  subscription_plan?: string;
  is_active?: boolean;
  updated_at?: string;
  created_at?: string;
};

export type AuthResponse = {
  token: string;
  access_token?: string;
  refresh_token?: string;
  expires_in?: number;
  user: AuthUser;
  usage?: UsageSnapshot;
};

export const ROLE_ADMIN = 1;

export function isAdmin(user: AuthUser | null | undefined): boolean {
  return user?.role === ROLE_ADMIN;
}
