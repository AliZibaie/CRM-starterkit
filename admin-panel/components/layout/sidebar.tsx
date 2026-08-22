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

const navGroups = [
  {
    label: "Overview",
    items: [
      { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
      { href: "/request-monitor", label: "Request Monitor", icon: Activity },
      { href: "/audit-log", label: "Audit Log", icon: ScrollText },
    ],
  },
  {
    label: "Administration",
    items: [
      { href: "/users", label: "Users", icon: Users },
      { href: "/roles-permissions", label: "Roles & Permissions", icon: Shield },
      { href: "/settings", label: "Settings", icon: Settings },
    ],
  },
  {
    label: "Automation",
    items: [
      { href: "/cron-jobs/builder", label: "Cron Builder", icon: Clock },
      { href: "/cron-jobs/logs", label: "Cron Logs", icon: Clock },
    ],
  },
  {
    label: "Content",
    items: [
      { href: "/announcements", label: "Announcements", icon: Megaphone },
      { href: "/banners", label: "Banners", icon: Image },
      { href: "/cms/pages", label: "CMS Pages", icon: FileText },
      { href: "/cms/menus", label: "CMS Menus", icon: FileText },
      { href: "/media-library", label: "Media Library", icon: FolderOpen },
    ],
  },
  {
    label: "Support",
    items: [
      { href: "/support/tickets", label: "Tickets", icon: Headphones },
      { href: "/support/departments", label: "Departments", icon: Headphones },
    ],
  },
  {
    label: "Commerce",
    items: [
      { href: "/catalog/products", label: "Products", icon: Package },
      { href: "/catalog/categories", label: "Categories", icon: Package },
      { href: "/catalog/inventory", label: "Inventory", icon: Package },
      { href: "/orders", label: "Orders", icon: ShoppingCart },
      { href: "/payments/gateways-status", label: "Payment Gateways", icon: CreditCard },
      { href: "/discounts", label: "Discounts", icon: Percent },
    ],
  },
  {
    label: "CRM",
    items: [
      { href: "/crm/leads", label: "Leads", icon: Contact },
      { href: "/crm/pipeline", label: "Pipeline", icon: Contact },
      { href: "/crm/contacts", label: "Contacts", icon: Contact },
      { href: "/crm/tasks", label: "Tasks", icon: Tag },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex h-full w-64 flex-col border-r border-border bg-sidebar text-sidebar-foreground">
      <div className="border-b border-border px-4 py-5">
        <span className="text-lg font-bold tracking-tight">CRM Admin</span>
      </div>
      <nav className="flex-1 overflow-y-auto p-3">
        {navGroups.map((group) => (
          <div key={group.label} className="mb-4">
            <p className="mb-2 px-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              {group.label}
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
                          ? "bg-sidebar-accent text-sidebar-accent-foreground font-medium"
                          : "hover:bg-sidebar-accent/50"
                      )}
                    >
                      <Icon className="h-4 w-4 shrink-0" />
                      {item.label}
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
