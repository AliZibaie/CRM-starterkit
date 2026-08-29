"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { useAuthStore } from "@/stores/auth-store";
import { useThemeStore } from "@/stores/theme-store";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30_000,
            retry: 1,
          },
        },
      })
  );
  const loadFromStorage = useAuthStore((s) => s.loadFromStorage);
  const hydrateTheme = useThemeStore((s) => s.hydrate);

  useEffect(() => {
    loadFromStorage();
    hydrateTheme();
  }, [loadFromStorage, hydrateTheme]);

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
