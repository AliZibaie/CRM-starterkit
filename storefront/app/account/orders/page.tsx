"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { formatPrice } from "@/lib/utils";
import { t, statusLabel } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";

interface Order {
  id: string;
  status: string;
  total_amount: number;
  currency: string;
  created_at: string;
}

export default function OrdersPage() {
  const { data: orders = [], isLoading, error } = useQuery({
    queryKey: queryKeys.orders.all,
    queryFn: async () => {
      const res = await api.get<Order[]>("/orders");
      return res.data ?? [];
    },
  });

  if (isLoading) {
    return <p className="text-muted-foreground">{t("orders.loading")}</p>;
  }

  if (error) {
    return <p className="text-destructive">{error.message}</p>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">{t("orders.title")}</h1>
          <p className="text-muted-foreground">{t("orders.subtitle")}</p>
        </div>
        <div className="flex gap-2">
          <Button asChild variant="outline">
            <Link href="/account/tickets">{t("nav.support")}</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/account/licenses">{t("nav.licenses")}</Link>
          </Button>
        </div>
      </div>
      <DataTable
        rows={orders}
        rowKey={(o) => o.id}
        emptyMessage={t("orders.empty")}
        columns={[
          {
            key: "id",
            header: "شناسه",
            cell: (o) => (
              <span className="font-mono text-xs ltr-field">{o.id.slice(0, 8)}…</span>
            ),
          },
          {
            key: "date",
            header: t("orders.date"),
            cell: (o) => formatJalaliDateTime(o.created_at),
          },
          {
            key: "status",
            header: t("orders.status"),
            cell: (o) => (
              <span className="rounded-full border px-3 py-1 text-xs">
                {statusLabel(o.status)}
              </span>
            ),
          },
          {
            key: "amount",
            header: t("orders.amount"),
            cell: (o) => (
              <span className="font-semibold">
                {formatPrice(o.total_amount, o.currency)}
              </span>
            ),
          },
        ]}
      />
    </div>
  );
}
