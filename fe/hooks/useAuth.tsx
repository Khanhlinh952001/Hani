"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import {
  fetchMe,
  login as apiLogin,
  register as apiRegister,
  updateProfile as apiUpdateProfile,
  uploadAvatar as apiUploadAvatar,
} from "@/lib/auth/api";
import type { UpdateProfileInput } from "@/lib/auth/api";
import {
  clearAuth,
  getStoredUser,
  getToken,
  setAuth,
} from "@/lib/auth/storage";
import { isAdmin, type AuthUser } from "@/lib/auth/types";

type AuthContextValue = {
  user: AuthUser | null;
  token: string | null;
  loading: boolean;
  isAdmin: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  updateProfile: (input: UpdateProfileInput) => Promise<void>;
  uploadAvatar: (file: File) => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const t = getToken();
    const u = getStoredUser();
    if (!t || !u) {
      setLoading(false);
      return;
    }
    setToken(t);
    setUser(u);
    fetchMe()
      .then((fresh) => {
        setUser(fresh);
        setAuth(t, fresh);
      })
      .catch(() => {
        clearAuth();
        setToken(null);
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await apiLogin(email, password);
    setAuth(res.token, res.user);
    setToken(res.token);
    setUser(res.user);
  }, []);

  const register = useCallback(
    async (name: string, email: string, password: string) => {
      const res = await apiRegister(name, email, password);
      setAuth(res.token, res.user);
      setToken(res.token);
      setUser(res.user);
    },
    []
  );

  const logout = useCallback(() => {
    clearAuth();
    setToken(null);
    setUser(null);
  }, []);

  const updateProfile = useCallback(async (input: UpdateProfileInput) => {
    const t = getToken();
    if (!t) throw new Error("not logged in");
    const fresh = await apiUpdateProfile(input);
    setUser(fresh);
    setAuth(t, fresh);
  }, []);

  const uploadAvatar = useCallback(async (file: File) => {
    const t = getToken();
    if (!t) throw new Error("not logged in");
    const fresh = await apiUploadAvatar(file);
    setUser(fresh);
    setAuth(t, fresh);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        loading,
        isAdmin: isAdmin(user),
        login,
        register,
        logout,
        updateProfile,
        uploadAvatar,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
