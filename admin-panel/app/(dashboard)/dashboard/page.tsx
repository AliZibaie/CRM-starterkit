"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface RequestStats {
  total: number;
  suspicious: number;
}

interface Ticket {
  id: string;
}

interface Order {
  id: string;
}

export default function DashboardPage() {
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

  const cards = [
    {
      title: t("pages.dashboard.users"),
      value: usersLoading ? "..." : String(usersMeta?.total ?? 0),
    },
    {
      title: t("pages.dashboard.orders"),
      value: ordersLoading ? "..." : String(orders.length),
    },
    {
      title: t("pages.dashboard.tickets"),
      value: ticketsLoading ? "..." : String(tickets.length),
    },
    {
      title: t("pages.dashboard.requests"),
      value: statsLoading ? "..." : String(stats?.total ?? 0),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.dashboard.title")}</h1>
        <p className="text-muted-foreground">{t("pages.dashboard.subtitle")}</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {cards.map((card) => (
          <Card key={card.title}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {card.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{card.value}</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
