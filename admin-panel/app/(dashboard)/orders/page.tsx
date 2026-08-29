"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface Order {
  id: string;
  status: string;
  total_amount: number;
  currency: string;
  created_at: string;
}

export default function OrdersPage() {
  const queryClient = useQueryClient();

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.orders,
    queryFn: async () => {
      const res = await api.get<Order[]>("/orders");
      return res.data ?? [];
    },
  });

  const checkoutMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/orders/${id}/checkout`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.orders }),
  });

  const refundMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/orders/${id}/refund`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.orders }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.orders.title")}</h1>
        <p className="text-muted-foreground">{t("pages.orders.subtitle")}</p>
      </div>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "id", header: t("common.id"), className: "font-mono text-xs" },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => <Badge variant="outline">{String(row.status)}</Badge>,
          },
          {
            key: "total_amount",
            header: t("common.total"),
            render: (row) =>
              `${row.currency} ${((row.total_amount as number) / 100).toFixed(2)}`,
          },
          {
            key: "created_at",
            header: t("common.createdAt"),
            render: (row) => formatJalaliDateTime(row.created_at as string),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => checkoutMutation.mutate(row.id as string)}
                >
                  {t("common.checkout")}
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => refundMutation.mutate(row.id as string)}
                >
                  {t("common.refund")}
                </Button>
              </div>
            ),
          },
        ]}
      />
    </div>
  );
}
