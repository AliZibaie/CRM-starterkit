"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Badge } from "@/components/ui/badge";

interface Gateway {
  name: string;
  enabled: boolean;
  configured: boolean;
}

export default function GatewaysStatusPage() {
  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.gateways,
    queryFn: async () => {
      const res = await api.get<Gateway[]>("/payments/gateways");
      return res.data ?? [];
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Payment Gateways</h1>
        <p className="text-muted-foreground">Gateway status and configuration</p>
      </div>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="name"
        columns={[
          { key: "name", header: "Gateway" },
          {
            key: "enabled",
            header: "Enabled",
            render: (row) => (
              <Badge variant={row.enabled ? "success" : "secondary"}>
                {row.enabled ? "Yes" : "No"}
              </Badge>
            ),
          },
          {
            key: "configured",
            header: "Configured",
            render: (row) => (
              <Badge variant={row.configured ? "success" : "destructive"}>
                {row.configured ? "Yes" : "No"}
              </Badge>
            ),
          },
        ]}
      />
    </div>
  );
}
