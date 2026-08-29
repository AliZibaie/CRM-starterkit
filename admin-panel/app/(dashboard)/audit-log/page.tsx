"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { datetimeLocalToRFC3339 } from "@/lib/datetime-input";
import { DataTable } from "@/components/table/data-table";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface AuditLog {
  id: string;
  action: string;
  module: string;
  entity_type: string;
  entity_id: string;
  user_id: string;
  created_at: string;
}

export default function AuditLogPage() {
  const [actorId, setActorId] = useState("");
  const [module, setModule] = useState("");
  const [entityType, setEntityType] = useState("");
  const [entityId, setEntityId] = useState("");
  const [fromLocal, setFromLocal] = useState("");
  const [toLocal, setToLocal] = useState("");

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    if (actorId) params.set("actor_id", actorId);
    if (module) params.set("module", module);
    if (entityType) params.set("entity_type", entityType);
    if (entityId) params.set("entity_id", entityId);
    const from = datetimeLocalToRFC3339(fromLocal);
    const to = datetimeLocalToRFC3339(toLocal);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    params.set("limit", "100");
    return params.toString();
  }, [actorId, module, entityType, entityId, fromLocal, toLocal]);

  const { data = [], isLoading } = useQuery({
    queryKey: [...queryKeys.auditLogs, queryString],
    queryFn: async () => {
      const res = await api.get<AuditLog[]>(`/audit-logs?${queryString}`);
      return res.data ?? [];
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.auditLog.title")}</h1>
        <p className="text-muted-foreground">{t("pages.auditLog.subtitle")}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("common.filters")}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
          <div className="space-y-2">
            <Label>{t("common.actorId")}</Label>
            <Input value={actorId} onChange={(e) => setActorId(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.module")}</Label>
            <Input value={module} onChange={(e) => setModule(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.entityType")}</Label>
            <Input value={entityType} onChange={(e) => setEntityType(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.entityId")}</Label>
            <Input value={entityId} onChange={(e) => setEntityId(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.from")}</Label>
            <Input type="datetime-local" value={fromLocal} onChange={(e) => setFromLocal(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.to")}</Label>
            <Input type="datetime-local" value={toLocal} onChange={(e) => setToLocal(e.target.value)} />
          </div>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "module", header: t("common.module") },
          { key: "action", header: t("common.action") },
          { key: "entity_type", header: t("common.entityType") },
          { key: "entity_id", header: t("common.entityId"), className: "font-mono text-xs" },
          { key: "user_id", header: t("common.userId"), className: "font-mono text-xs" },
          {
            key: "created_at",
            header: t("common.createdAt"),
            render: (row) => formatJalaliDateTime(row.created_at as string),
          },
        ]}
      />
    </div>
  );
}
