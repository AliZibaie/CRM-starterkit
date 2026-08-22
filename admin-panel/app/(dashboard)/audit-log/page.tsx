"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface AuditLog {
  id: string;
  action: string;
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
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    if (actorId) params.set("actor_id", actorId);
    if (module) params.set("module", module);
    if (entityType) params.set("entity_type", entityType);
    if (entityId) params.set("entity_id", entityId);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    params.set("limit", "100");
    return params.toString();
  }, [actorId, module, entityType, entityId, from, to]);

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
        <h1 className="text-2xl font-bold">Audit Log</h1>
        <p className="text-muted-foreground">System activity audit trail</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
          <div className="space-y-2">
            <Label>Actor ID</Label>
            <Input value={actorId} onChange={(e) => setActorId(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Module</Label>
            <Input value={module} onChange={(e) => setModule(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Entity Type</Label>
            <Input value={entityType} onChange={(e) => setEntityType(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Entity ID</Label>
            <Input value={entityId} onChange={(e) => setEntityId(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>From (RFC3339)</Label>
            <Input value={from} onChange={(e) => setFrom(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>To (RFC3339)</Label>
            <Input value={to} onChange={(e) => setTo(e.target.value)} />
          </div>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "action", header: "Action" },
          { key: "entity_type", header: "Entity Type" },
          { key: "entity_id", header: "Entity ID" },
          { key: "user_id", header: "User ID" },
          { key: "created_at", header: "Created" },
        ]}
      />
    </div>
  );
}
