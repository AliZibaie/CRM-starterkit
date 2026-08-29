"use client";

import { useRouter } from "next/navigation";
import { Bell, LogOut, Menu, Moon, Search, Sun, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuthStore } from "@/stores/auth-store";
import { useThemeStore } from "@/stores/theme-store";
import { apiClient } from "@/lib/api-client";
import { t } from "@/lib/i18n";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";

interface HeaderProps {
  onMenuClick?: () => void;
}

export function Header({ onMenuClick }: HeaderProps) {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const roles = useAuthStore((s) => s.roles);
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const logout = useAuthStore((s) => s.logout);
  const theme = useThemeStore((s) => s.theme);
  const toggleTheme = useThemeStore((s) => s.toggleTheme);

  const { data: tickets = [] } = useQuery({
    queryKey: queryKeys.tickets,
    queryFn: async () => {
      const res = await api.get<{ status: string }[]>("/support/tickets");
      return res.data ?? [];
    },
    staleTime: 60_000,
  });

  const openTickets = tickets.filter((tk) => tk.status === "open").length;
  const initials = (user?.email?.[0] ?? "A").toUpperCase();

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await apiClient("/auth/logout", {
          method: "POST",
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      }
    } catch {
    } finally {
      logout();
      router.push("/login");
    }
  };

  return (
    <header className="glass-header sticky top-0 z-30 flex h-16 items-center justify-between gap-4 border-b border-border px-4 md:px-6">
      <div className="flex min-w-0 flex-1 items-center gap-3">
        <Button
          variant="outline"
          size="icon"
          className="shrink-0 lg:hidden"
          onClick={onMenuClick}
        >
          <Menu className="h-5 w-5" />
        </Button>
        <div className="relative hidden max-w-md flex-1 md:block">
          <Search className="pointer-events-none absolute start-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder={t("header.searchPlaceholder")}
            className="h-10 rounded-xl border-border/70 bg-background/70 ps-10"
          />
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="icon"
          className="rounded-xl"
          onClick={toggleTheme}
          aria-label={theme === "dark" ? "light" : "dark"}
        >
          {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </Button>
        <Button
          variant="outline"
          size="icon"
          className="relative rounded-xl"
          onClick={() => router.push("/support/tickets")}
        >
          <Bell className="h-4 w-4" />
          {openTickets > 0 && (
            <Badge className="absolute -top-1 end-0 flex h-5 min-w-5 items-center justify-center rounded-full px-1 text-[10px]">
              {openTickets}
            </Badge>
          )}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-10 gap-2 rounded-xl px-2">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="text-xs">{initials}</AvatarFallback>
              </Avatar>
              <span className="hidden max-w-[140px] truncate text-sm md:inline">
                {user?.email || t("common.admin")}
              </span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-56">
            <div className="px-3 py-2">
              <p className="truncate text-sm font-medium">{user?.email}</p>
              <p className="truncate text-xs text-muted-foreground">{roles.join(", ")}</p>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => router.push("/users")}>
              <User className="ms-2 h-4 w-4" />
              {t("nav.users")}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => router.push("/settings")}>
              <Search className="ms-2 h-4 w-4" />
              {t("nav.settings")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} className="text-destructive">
              <LogOut className="ms-2 h-4 w-4 rtl-icon-mirror" />
              {t("common.logout")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
