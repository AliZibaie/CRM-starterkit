"use client";

import { create } from "zustand";
import { jwtDecode } from "jwt-decode";

const STORAGE_KEY = "storefront-auth";

interface JwtPayload {
  uid: string;
  email?: string;
  permissions?: string[];
  roles?: string[];
}

export interface AuthUser {
  id: string;
  email: string;
}

interface StoredAuth {
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthUser | null;
}

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthUser | null;
  permissions: string[];
  roles: string[];
  login: (accessToken: string, refreshToken: string, user: AuthUser) => void;
  logout: () => void;
  setTokens: (accessToken: string, refreshToken: string) => void;
  loadFromStorage: () => void;
}

function decodeToken(token: string) {
  const payload = jwtDecode<JwtPayload>(token);
  return {
    permissions: payload.permissions ?? [],
    roles: payload.roles ?? [],
  };
}

function setAuthCookie(token: string | null) {
  if (typeof document === "undefined") return;
  if (token) {
    document.cookie = `access_token=${token}; path=/; max-age=86400; SameSite=Lax`;
  } else {
    document.cookie = "access_token=; path=/; max-age=0";
  }
}

function saveToStorage(data: StoredAuth) {
  if (typeof window === "undefined") return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
}

function readFromStorage(): StoredAuth | null {
  if (typeof window === "undefined") return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as StoredAuth;
  } catch {
    return null;
  }
}

export const useAuthStore = create<AuthState>()((set) => ({
  accessToken: null,
  refreshToken: null,
  user: null,
  permissions: [],
  roles: [],
  login: (accessToken, refreshToken, user) => {
    const decoded = decodeToken(accessToken);
    setAuthCookie(accessToken);
    saveToStorage({ accessToken, refreshToken, user });
    set({
      accessToken,
      refreshToken,
      user,
      permissions: decoded.permissions,
      roles: decoded.roles,
    });
  },
  logout: () => {
    setAuthCookie(null);
    if (typeof window !== "undefined") {
      localStorage.removeItem(STORAGE_KEY);
    }
    set({
      accessToken: null,
      refreshToken: null,
      user: null,
      permissions: [],
      roles: [],
    });
  },
  setTokens: (accessToken, refreshToken) => {
    const decoded = decodeToken(accessToken);
    const user = useAuthStore.getState().user;
    setAuthCookie(accessToken);
    saveToStorage({ accessToken, refreshToken, user });
    set({
      accessToken,
      refreshToken,
      permissions: decoded.permissions,
      roles: decoded.roles,
    });
  },
  loadFromStorage: () => {
    const stored = readFromStorage();
    if (!stored?.accessToken) return;
    const decoded = decodeToken(stored.accessToken);
    setAuthCookie(stored.accessToken);
    set({
      accessToken: stored.accessToken,
      refreshToken: stored.refreshToken,
      user: stored.user,
      permissions: decoded.permissions,
      roles: decoded.roles,
    });
  },
}));

export function getStoredTokens(): StoredAuth | null {
  return readFromStorage();
}
