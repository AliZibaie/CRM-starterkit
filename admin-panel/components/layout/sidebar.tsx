"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  LayoutDashboard,
  Users,
  Shield,
  Clock,
  Megaphone,
  Headphones,
  Image,
  Package,
  ShoppingCart,
  CreditCard,
  Tag,
  Contact,
  FileText,
  FolderOpen,
  Settings,
  ScrollText,
  Activity,
  Percent,
  LogOut,
  X,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { t } from "@/lib/i18n";
import { useAuthStore } from "@/stores/auth-store";
import { apiClient } from "@/lib/api-client";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

const navGroups = [
  {
    labelKey: "nav.overview",
    items: [
      { href: "/dashboard", labelKey: "nav.dashboard", icon: LayoutDashboard },
      { href: "/request-monitor", labelKey: "nav.requestMonitor", icon: Activity },
      { href: "/audit-log", labelKey: "nav.auditLog", icon: ScrollText },
    ],
  },
  {
    labelKey: "nav.administration",
    items: [
      { href: "/users", labelKey: "nav.users", icon: Users },
      { href: "/roles-permissions", labelKey: "nav.rolesPermissions", icon: Shield },
      { href: "/settings", labelKey: "nav.settings", icon: Settings },
    ],
  },
  {
    labelKey: "nav.automation",
    items: [
      { href: "/cron-jobs/builder", labelKey: "nav.cronBuilder", icon: Clock },
      { href: "/cron-jobs/logs", labelKey: "nav.cronLogs", icon: Clock },
    ],
  },
  {
    labelKey: "nav.content",
    items: [
      { href: "/announcements", labelKey: "nav.announcements", icon: Megaphone },
      { href: "/banners", labelKey: "nav.banners", icon: Image },
      { href: "/cms/pages", labelKey: "nav.cmsPages", icon: FileText },
      { href: "/cms/menus", labelKey: "nav.cmsMenus", icon: FileText },
      { href: "/media-library", labelKey: "nav.mediaLibrary", icon: FolderOpen },
    ],
  },
  {
    labelKey: "nav.support",
    items: [
      { href: "/support/tickets", labelKey: "nav.tickets", icon: Headphones },
      { href: "/support/departments", labelKey: "nav.departments", icon: Headphones },
    ],
  },
  {
    labelKey: "nav.commerce",
    items: [
      { href: "/catalog/products", labelKey: "nav.products", icon: Package },
      { href: "/catalog/categories", labelKey: "nav.categories", icon: Package },
      { href: "/catalog/inventory", labelKey: "nav.inventory", icon: Package },
      { href: "/orders", labelKey: "nav.orders", icon: ShoppingCart },
      { href: "/payments/gateways-status", labelKey: "nav.paymentGateways", icon: CreditCard },
      { href: "/discounts", labelKey: "nav.discounts", icon: Percent },
    ],
  },
  {
    labelKey: "nav.crm",
    items: [
      { href: "/crm/leads", labelKey: "nav.leads", icon: Contact },
      { href: "/crm/pipeline", labelKey: "nav.pipeline", icon: Contact },
      { href: "/crm/contacts", labelKey: "nav.contacts", icon: Contact },
      { href: "/crm/tasks", labelKey: "nav.tasks", icon: Tag },
    ],
  },
];

interface SidebarProps {
  mobileOpen?: boolean;
  onMobileClose?: () => void;
  className?: string;
}

export function Sidebar({ mobileOpen, onMobileClose, className }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const roles = useAuthStore((s) => s.roles);
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const logout = useAuthStore((s) => s.logout);

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
    <aside
      className={cn(
        "sidebar-shell flex h-full w-72 shrink-0 flex-col border-s border-sidebar-border text-sidebar-foreground shadow-2xl transition-transform duration-300 lg:shadow-none",
        className
      )}
    >
      <div className="flex items-center justify-between border-b border-sidebar-border px-5 py-5">
        <div>
          <p className="text-lg font-bold text-white">{t("brand")}</p>
          <p className="text-xs text-sidebar-muted">{t("pages.dashboard.subtitle")}</p>
        </div>
        {mobileOpen && (
          <Button
            variant="ghost"
            size="icon"
            className="text-sidebar-foreground hover:bg-white/10 lg:hidden"
            onClick={onMobileClose}
          >
            <X className="h-5 w-5" />
          </Button>
        )}
      </div>
      <div className="sidebar-scroll flex-1 overflow-y-auto">
        <nav className="sidebar-scroll-inner p-3">
          {navGroups.map((group) => (
            <div key={group.labelKey} className="mb-5">
              <p className="mb-2 px-3 text-[11px] font-semibold tracking-wide text-sidebar-muted">
                {t(group.labelKey)}
              </p>
              <ul className="space-y-1">
                {group.items.map((item) => {
                  const Icon = item.icon;
                  const active =
                    pathname === item.href || pathname.startsWith(item.href + "/");
                  return (
                    <li key={item.href}>
                      <Link
                        href={item.href}
                        onClick={onMobileClose}
                        className={cn(
                          "group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm transition-all",
                          active
                            ? "bg-primary text-primary-foreground shadow-lg shadow-primary/25"
                            : "text-sidebar-foreground/90 hover:bg-white/10 hover:text-white"
                        )}
                      >
                        <Icon className={cn("h-4 w-4 shrink-0", active && "text-white")} />
                        <span>{t(item.labelKey)}</span>
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          ))}
        </nav>
      </div>
      <div className="border-t border-sidebar-border p-4">
        <div className="rounded-2xl bg-white/8 p-3 backdrop-blur-sm">
          <div className="flex items-center gap-3">
            <Avatar className="h-11 w-11 ring-white/20">
              <AvatarFallback>{initials}</AvatarFallback>
            </Avatar>
            <div className="min-w-0 flex-1 text-start">
              <p className="truncate text-sm font-semibold text-white">
                {user?.email || t("common.admin")}
              </p>
              <p className="truncate text-xs text-sidebar-muted">
                {roles[0] || t("common.admin")}
              </p>
            </div>
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="mt-3 w-full justify-start gap-2 text-sidebar-foreground hover:bg-white/10 hover:text-white"
            onClick={handleLogout}
          >
            <LogOut className="h-4 w-4 rtl-icon-mirror" />
            {t("common.logout")}
          </Button>
        </div>
      </div>
    </aside>
  );
}
