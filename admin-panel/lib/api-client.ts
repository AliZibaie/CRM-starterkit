import { getStoredTokens, useAuthStore } from "@/stores/auth-store";

export const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export interface ApiEnvelope<T = unknown> {
  success: boolean;
  data?: T;
  message?: string;
  meta?: Record<string, unknown>;
}

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

interface RefreshData {
  access_token: string;
  refresh_token: string;
}

let refreshPromise: Promise<boolean> | null = null;

function getAccessToken(): string | null {
  const fromStore = useAuthStore.getState().accessToken;
  if (fromStore) return fromStore;
  return getStoredTokens()?.accessToken ?? null;
}

function getRefreshToken(): string | null {
  const fromStore = useAuthStore.getState().refreshToken;
  if (fromStore) return fromStore;
  return getStoredTokens()?.refreshToken ?? null;
}

async function refreshAccessToken(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  const { setTokens, logout } = useAuthStore.getState();
  if (!refreshToken) {
    logout();
    return false;
  }
  try {
    const res = await fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    const json = (await res.json()) as ApiEnvelope<RefreshData>;
    if (!res.ok || !json.success || !json.data) {
      logout();
      return false;
    }
    setTokens(json.data.access_token, json.data.refresh_token);
    return true;
  } catch {
    logout();
    return false;
  }
}

async function ensureRefreshed(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = refreshAccessToken().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

async function parseEnvelope<T>(res: Response): Promise<ApiEnvelope<T>> {
  const text = await res.text();
  if (!text) {
    return {
      success: res.ok,
      message: res.statusText || undefined,
    };
  }
  try {
    return JSON.parse(text) as ApiEnvelope<T>;
  } catch {
    throw new ApiError(text || res.statusText || "Request failed", res.status);
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  retry = true
): Promise<ApiEnvelope<T>> {
  const token = getAccessToken();
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  };
  if (!(options.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  let res: Response;
  let json: ApiEnvelope<T>;
  try {
    res = await fetch(`${API_BASE}${path}`, { ...options, headers });
    json = await parseEnvelope<T>(res);
  } catch (err) {
    if (err instanceof ApiError) {
      throw err;
    }
    throw new ApiError("network_error", 0);
  }
  if (
    res.status === 401 &&
    retry &&
    path !== "/auth/login" &&
    path !== "/auth/refresh"
  ) {
    const refreshed = await ensureRefreshed();
    if (refreshed) {
      return request<T>(path, options, false);
    }
    throw new ApiError(json.message || "Unauthorized", res.status);
  }
  if (!res.ok || !json.success) {
    throw new ApiError(json.message || "Request failed", res.status);
  }
  return json;
}

export async function apiClient<T>(
  path: string,
  options: RequestInit = {}
): Promise<ApiEnvelope<T>> {
  return request<T>(path, options);
}

export const api = {
  get: <T>(path: string) => apiClient<T>(path),
  post: <T>(path: string, body?: unknown) =>
    apiClient<T>(path, {
      method: "POST",
      body: JSON.stringify(body ?? {}),
    }),
  put: <T>(path: string, body?: unknown) =>
    apiClient<T>(path, {
      method: "PUT",
      body: JSON.stringify(body ?? {}),
    }),
  delete: <T>(path: string) => apiClient<T>(path, { method: "DELETE" }),
};

export async function apiUpload<T>(
  path: string,
  formData: FormData
): Promise<ApiEnvelope<T>> {
  return request<T>(path, { method: "POST", body: formData });
}
