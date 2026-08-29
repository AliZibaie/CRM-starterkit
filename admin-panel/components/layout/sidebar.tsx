"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
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
} from "lucide-react";
import { cn } from "@/lib/utils";
import { t } from "@/lib/i18n";

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

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex h-full w-64 shrink-0 flex-col border-s border-border bg-sidebar text-sidebar-foreground">
      <div className="border-b border-border px-4 py-5">
        <span className="text-lg font-bold tracking-tight">{t("brand")}</span>
      </div>
      <nav className="flex-1 overflow-y-auto p-3">
        {navGroups.map((group) => (
          <div key={group.labelKey} className="mb-4">
            <p className="mb-2 px-2 text-xs font-semibold tracking-wider text-muted-foreground">
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
                      className={cn(
                        "flex items-center gap-2 rounded-md px-2 py-2 text-sm transition-colors",
                        active
                          ? "bg-sidebar-accent font-medium text-sidebar-accent-foreground"
                          : "hover:bg-sidebar-accent/50"
                      )}
                    >
                      <Icon className="h-4 w-4 shrink-0" />
                      {t(item.labelKey)}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
}
