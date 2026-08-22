"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  ConditionBuilder,
  ConditionNode,
} from "@/components/form-builder/condition-builder";

interface TableMeta {
  name: string;
}

interface ColumnMeta {
  name: string;
  type: string;
}

interface ActionInfo {
  type: string;
  description: string;
}

export default function CronBuilderPage() {
  const queryClient = useQueryClient();
  const [tableName, setTableName] = useState("");
  const [jobName, setJobName] = useState("");
  const [cronExpr, setCronExpr] = useState("0 * * * *");
  const [actionType, setActionType] = useState("");
  const [condition, setCondition] = useState<ConditionNode>({
    field: "",
    operator: "eq",
    value: "",
  });
  const [buildResult, setBuildResult] = useState("");

  const { data: tables = [] } = useQuery({
    queryKey: queryKeys.cronTables,
    queryFn: async () => {
      const res = await api.get<TableMeta[]>("/cronbuilder/tables");
      return res.data ?? [];
    },
  });

  const { data: columns = [] } = useQuery({
    queryKey: ["cron-columns", tableName],
    queryFn: async () => {
      const res = await api.get<ColumnMeta[]>(
        `/cronbuilder/tables/${tableName}/columns`
      );
      return res.data ?? [];
    },
    enabled: !!tableName,
  });

  const { data: actions = [] } = useQuery({
    queryKey: queryKeys.cronActions,
    queryFn: async () => {
      const res = await api.get<ActionInfo[]>("/cronbuilder/actions");
      return res.data ?? [];
    },
  });

  const buildMutation = useMutation({
    mutationFn: async () => {
      const res = await api.post<Record<string, unknown>>(
        "/cronbuilder/conditions/build",
        { table_name: tableName, condition }
      );
      return res.data;
    },
    onSuccess: (data) => setBuildResult(JSON.stringify(data, null, 2)),
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/cronbuilder/jobs", {
        name: jobName,
        table_name: tableName,
        condition,
        action: { type: actionType, params: {} },
        cron_expr: cronExpr,
        is_active: true,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cronJobs });
      setJobName("");
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Cron Job Builder</h1>
        <p className="text-muted-foreground">Visual cron job configuration</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Job Configuration</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Job Name</Label>
              <Input value={jobName} onChange={(e) => setJobName(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Cron Expression</Label>
              <Input value={cronExpr} onChange={(e) => setCronExpr(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Table</Label>
              <Select value={tableName} onValueChange={setTableName}>
                <SelectTrigger>
                  <SelectValue placeholder="Select table" />
                </SelectTrigger>
                <SelectContent>
                  {tables.map((t) => (
                    <SelectItem key={t.name} value={t.name}>
                      {t.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Action</Label>
              <Select value={actionType} onValueChange={setActionType}>
                <SelectTrigger>
                  <SelectValue placeholder="Select action" />
                </SelectTrigger>
                <SelectContent>
                  {actions.map((a) => (
                    <SelectItem key={a.type} value={a.type}>
                      {a.type} - {a.description}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {tableName && (
        <ConditionBuilder
          columns={columns}
          value={condition}
          onChange={setCondition}
          onBuild={() => buildMutation.mutate()}
          buildResult={buildResult}
        />
      )}

      <Button
        onClick={() => createMutation.mutate()}
        disabled={createMutation.isPending || !jobName || !tableName || !actionType}
      >
        Create Cron Job
      </Button>
    </div>
  );
}
