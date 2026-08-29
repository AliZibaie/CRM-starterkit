"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
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
  const [buildError, setBuildError] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: tables = [], isError: tablesError, error: tablesErr, refetch: refetchTables } =
    useQuery({
      queryKey: queryKeys.cronTables,
      queryFn: async () => {
        const res = await api.get<TableMeta[]>("/cronbuilder/tables");
        return res.data ?? [];
      },
    });

  const { data: columns = [] } = useQuery({
    queryKey: ["cron-columns", tableName],
    queryFn: async () => {
      const res = await api.get<ColumnMeta[]>(`/cronbuilder/tables/${tableName}/columns`);
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
      const res = await api.post<Record<string, unknown>>("/cronbuilder/conditions/build", {
        table_name: tableName,
        condition,
      });
      return res.data;
    },
    onSuccess: (data) => {
      setBuildError("");
      setBuildResult(JSON.stringify(data, null, 2));
    },
    onError: (err) => {
      setBuildResult("");
      setBuildError(getApiErrorMessage(err));
    },
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
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.cronBuilder.title")}</h1>
        <p className="text-muted-foreground">{t("pages.cronBuilder.subtitle")}</p>
      </div>

      {feedback && (
        <div
          className={`rounded-xl border px-4 py-3 text-sm ${
            feedback.type === "error"
              ? "border-destructive/40 bg-destructive/10 text-destructive"
              : "border-emerald-500/40 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
          }`}
        >
          {feedback.text}
        </div>
      )}

      {tablesError && (
        <div className="rounded-xl border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {getApiErrorMessage(tablesErr)}
          <Button size="sm" variant="outline" className="ms-3" onClick={() => refetchTables()}>
            {t("common.refresh")}
          </Button>
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.cronBuilder.jobConfiguration")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>{t("common.jobName")}</Label>
              <Input value={jobName} onChange={(e) => setJobName(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.cronExpression")}</Label>
              <Input value={cronExpr} onChange={(e) => setCronExpr(e.target.value)} className="ltr-field" />
            </div>
            <div className="space-y-2">
              <Label>{t("common.table")}</Label>
              <Select value={tableName} onValueChange={setTableName}>
                <SelectTrigger>
                  <SelectValue placeholder={t("common.selectTable")} />
                </SelectTrigger>
                <SelectContent>
                  {tables.map((tbl) => (
                    <SelectItem key={tbl.name} value={tbl.name}>
                      {tbl.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.action")}</Label>
              <Select value={actionType} onValueChange={setActionType}>
                <SelectTrigger>
                  <SelectValue placeholder={t("common.selectAction")} />
                </SelectTrigger>
                <SelectContent>
                  {actions.map((a) => (
                    <SelectItem key={a.type} value={a.type}>
                      {a.type} — {a.description}
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
          buildError={buildError}
          building={buildMutation.isPending}
        />
      )}

      <Button
        onClick={() => createMutation.mutate()}
        disabled={createMutation.isPending || !jobName || !tableName || !actionType}
      >
        {createMutation.isPending ? t("common.saving") : t("pages.cronBuilder.create")}
      </Button>
    </div>
  );
}
