"use client";

import { useRouter } from "next/navigation";
import { LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAuthStore } from "@/stores/auth-store";
import { apiClient } from "@/lib/api-client";

export function Header() {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const roles = useAuthStore((s) => s.roles);
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const logout = useAuthStore((s) => s.logout);

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
    <header className="flex h-14 items-center justify-between border-b border-border bg-background px-6">
      <div />
      <div className="flex items-center gap-3">
        <div className="text-right text-sm">
          <p className="font-medium">{user?.email || "Admin"}</p>
          <div className="flex flex-wrap justify-end gap-1">
            {roles.slice(0, 2).map((role) => (
              <Badge key={role} variant="secondary" className="text-xs">
                {role}
              </Badge>
            ))}
          </div>
        </div>
        <Button variant="ghost" size="icon" onClick={handleLogout}>
          <LogOut className="h-4 w-4" />
        </Button>
      </div>
    </header>
  );
}
