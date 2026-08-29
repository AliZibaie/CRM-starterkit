"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import {
  Activity,
  Headphones,
  ShoppingCart,
  Users,
} from "lucide-react";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatRelativeTime } from "@/lib/relative-time";
import { StatCard } from "@/components/dashboard/stat-card";
import { DashboardAreaChart } from "@/components/charts/dashboard-area-chart";
import { DashboardPieChart } from "@/components/charts/dashboard-pie-chart";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface RequestStats {
  total: number;
  suspicious: number;
}

interface RequestLog {
  created_at: string;
}

interface Ticket {
  id: string;
  subject: string;
  status: string;
  created_at: string;
}

interface Order {
  id: string;
  status: string;
  created_at: string;
}

interface AuditLog {
  id: string;
  action: string;
  entity_type: string;
  created_at: string;
}

function dayKey(date: Date) {
  return date.toLocaleDateString("fa-IR", { month: "short", day: "numeric" });
}

function groupByDay(items: { created_at: string }[], days = 7) {
  const map = new Map<string, number>();
  const now = new Date();
  for (let i = days - 1; i >= 0; i -= 1) {
    const d = new Date(now);
    d.setDate(now.getDate() - i);
    map.set(dayKey(d), 0);
  }
  items.forEach((item) => {
    const key = dayKey(new Date(item.created_at));
    if (map.has(key)) map.set(key, (map.get(key) ?? 0) + 1);
  });
  return Array.from(map.entries()).map(([name, value]) => ({ name, value }));
}

export default function DashboardPage() {
  const weekAgo = useMemo(() => {
    const d = new Date();
    d.setDate(d.getDate() - 7);
    return d.toISOString();
  }, []);

  const { data: usersMeta, isLoading: usersLoading } = useQuery({
    queryKey: queryKeys.users(1),
    queryFn: async () => {
      const res = await api.get<unknown[]>("/users?page=1&per_page=1");
      return res.meta;
    },
  });

  const { data: orders = [], isLoading: ordersLoading } = useQuery({
    queryKey: queryKeys.orders,
    queryFn: async () => {
      const res = await api.get<Order[]>("/orders");
      return res.data ?? [];
    },
  });

  const { data: tickets = [], isLoading: ticketsLoading } = useQuery({
    queryKey: queryKeys.tickets,
    queryFn: async () => {
      const res = await api.get<Ticket[]>("/support/tickets");
      return res.data ?? [];
    },
  });

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: queryKeys.requestStats({}),
    queryFn: async () => {
      const res = await api.get<RequestStats>("/request-logs/stats");
      return res.data!;
    },
  });

  const { data: requestLogs = [] } = useQuery({
    queryKey: ["dashboard-request-logs", weekAgo],
    queryFn: async () => {
      const res = await api.get<RequestLog[]>(
        `/request-logs?from=${encodeURIComponent(weekAgo)}&limit=500`
      );
      return res.data ?? [];
    },
  });

  const { data: auditLogs = [] } = useQuery({
    queryKey: ["dashboard-audit-logs"],
    queryFn: async () => {
      const res = await api.get<AuditLog[]>("/audit-logs?limit=8");
      return res.data ?? [];
    },
  });

  const requestTrend = useMemo(() => groupByDay(requestLogs), [requestLogs]);

  const ticketDistribution = useMemo(() => {
    const map = new Map<string, number>();
    tickets.forEach((tk) => {
      map.set(tk.status, (map.get(tk.status) ?? 0) + 1);
    });
    return Array.from(map.entries()).map(([name, value]) => ({ name, value }));
  }, [tickets]);

  const recentActivities = useMemo(() => {
    const fromTickets = tickets.slice(0, 5).map((tk) => ({
      id: `ticket-${tk.id}`,
      title: tk.subject,
      subtitle: tk.status,
      href: `/support/tickets/${tk.id}`,
      created_at: tk.created_at,
      kind: "ticket" as const,
    }));
    const fromAudit = auditLogs.slice(0, 5).map((log) => ({
      id: `audit-${log.id}`,
      title: log.action,
      subtitle: log.entity_type,
      href: "/audit-log",
      created_at: log.created_at,
      kind: "audit" as const,
    }));
    return [...fromTickets, ...fromAudit]
      .sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at))
      .slice(0, 8);
  }, [tickets, auditLogs]);

  const cards = [
    {
      title: t("pages.dashboard.users"),
      value: usersLoading ? "..." : String(usersMeta?.total ?? 0),
      icon: Users,
      tone: "primary" as const,
    },
    {
      title: t("pages.dashboard.orders"),
      value: ordersLoading ? "..." : String(orders.length),
      icon: ShoppingCart,
      tone: "success" as const,
    },
    {
      title: t("pages.dashboard.tickets"),
      value: ticketsLoading ? "..." : String(tickets.length),
      icon: Headphones,
      tone: "warning" as const,
    },
    {
      title: t("pages.dashboard.requests"),
      value: statsLoading ? "..." : String(stats?.total ?? 0),
      icon: Activity,
      tone: "info" as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t("pages.dashboard.title")}</h1>
        <p className="mt-1 text-muted-foreground">{t("pages.dashboard.subtitle")}</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {cards.map((card) => (
          <StatCard key={card.title} {...card} />
        ))}
      </div>

      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.dashboard.requestTrend")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DashboardAreaChart data={requestTrend} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.dashboard.ticketDistribution")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DashboardPieChart data={ticketDistribution} />
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.dashboard.recentActivity")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {recentActivities.length === 0 ? (
            <p className="text-sm text-muted-foreground">{t("common.noData")}</p>
          ) : (
            recentActivities.map((item) => (
              <Link
                key={item.id}
                href={item.href}
                className="flex items-center justify-between gap-4 rounded-xl border border-border/70 bg-muted/20 px-4 py-3 transition-colors hover:bg-muted/40"
              >
                <div className="min-w-0 text-start">
                  <p className="truncate font-medium">{item.title}</p>
                  <p className="truncate text-xs text-muted-foreground">{item.subtitle}</p>
                </div>
                <div className="flex shrink-0 items-center gap-2">
                  <Badge variant="secondary">
                    {item.kind === "ticket" ? t("nav.tickets") : t("nav.auditLog")}
                  </Badge>
                  <span className="text-xs text-muted-foreground">
                    {formatRelativeTime(item.created_at)}
                  </span>
                </div>
              </Link>
            ))
          )}
        </CardContent>
      </Card>
    </div>
  );
}
