"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Gateway {
  name: string;
  enabled: boolean;
  configured: boolean;
}

interface FailedPayment {
  id: string;
  order_id: string;
  gateway: string;
  amount: number;
  currency: string;
  status: string;
  authority?: string;
  created_at: string;
}

export default function GatewaysStatusPage() {
  const { data: gateways = [], isLoading } = useQuery({
    queryKey: queryKeys.gateways,
    queryFn: async () => {
      const res = await api.get<Gateway[]>("/payments/gateways");
      return res.data ?? [];
    },
  });

  const { data: failed = [], isLoading: failedLoading } = useQuery({
    queryKey: queryKeys.failedPayments,
    queryFn: async () => {
      const res = await api.get<FailedPayment[]>("/payments/failed");
      return res.data ?? [];
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.payments.title")}</h1>
        <p className="text-muted-foreground">{t("pages.payments.subtitle")}</p>
      </div>

      <DataTable
        loading={isLoading}
        data={gateways as unknown as Record<string, unknown>[]}
        keyField="name"
        columns={[
          { key: "name", header: t("common.gateway") },
          {
            key: "configured",
            header: t("common.configured"),
            render: (row) => (
              <Badge variant={row.configured ? "success" : "destructive"}>
                {row.configured ? t("common.yes") : t("common.no")}
              </Badge>
            ),
          },
          {
            key: "enabled",
            header: t("common.enabled"),
            render: (row) => (
              <Badge variant={row.enabled ? "success" : "secondary"}>
                {row.enabled ? t("common.yes") : t("common.no")}
              </Badge>
            ),
          },
        ]}
      />

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.payments.failedTransactions")}</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable
            loading={failedLoading}
            data={failed as unknown as Record<string, unknown>[]}
            keyField="id"
            columns={[
              { key: "order_id", header: t("common.id"), className: "font-mono text-xs" },
              { key: "gateway", header: t("common.gateway") },
              {
                key: "amount",
                header: t("common.amount"),
                render: (row) =>
                  `${row.currency} ${((row.amount as number) / 100).toFixed(2)}`,
              },
              {
                key: "status",
                header: t("common.status"),
                render: () => (
                  <Badge variant="destructive">{t("common.error")}</Badge>
                ),
              },
              {
                key: "created_at",
                header: t("common.createdAt"),
                render: (row) => formatJalaliDateTime(row.created_at as string),
              },
            ]}
          />
        </CardContent>
      </Card>
    </div>
  );
}
