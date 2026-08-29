"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ConditionNode } from "@/components/form-builder/condition-builder";

interface CronJob {
  id: string;
  name: string;
  table_name: string;
  cron_expr: string;
  is_active: boolean;
  condition?: ConditionNode;
  action?: { type: string; params?: Record<string, unknown> };
}

interface CronLog {
  id: string;
  job_id: string;
  started_at: string;
  finished_at: string;
  success: boolean;
  affected_rows: number;
  error_message: string;
}

export default function CronLogsPage() {
  const queryClient = useQueryClient();
  const [selectedJobId, setSelectedJobId] = useState<string | null>(null);
  const [editJob, setEditJob] = useState<CronJob | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: jobs = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.cronJobs,
    queryFn: async () => {
      const res = await api.get<CronJob[]>("/cronbuilder/jobs");
      return res.data ?? [];
    },
  });

  const { data: logs = [], isLoading: logsLoading } = useQuery({
    queryKey: queryKeys.cronLogs(selectedJobId ?? ""),
    queryFn: async () => {
      const res = await api.get<CronLog[]>(
        `/cronbuilder/jobs/${selectedJobId}/logs?limit=50`
      );
      return res.data ?? [];
    },
    enabled: !!selectedJobId,
  });

  const runMutation = useMutation({
    mutationFn: async (jobId: string) => {
      await api.post(`/cronbuilder/jobs/${jobId}/run`);
    },
    onSuccess: (_, jobId) => {
      setFeedback({ type: "success", text: t("common.operationSuccess") });
      if (selectedJobId === jobId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.cronLogs(jobId) });
      }
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async (job: CronJob) => {
      await api.put(`/cronbuilder/jobs/${job.id}`, {
        name: job.name,
        table_name: job.table_name,
        cron_expr: job.cron_expr,
        is_active: job.is_active,
        condition: job.condition ?? { field: "", operator: "eq", value: "" },
        action: job.action ?? { type: "noop", params: {} },
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cronJobs });
      setEditJob(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/cronbuilder/jobs/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cronJobs });
      if (selectedJobId === deleteId) setSelectedJobId(null);
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const openEdit = async (id: string) => {
    try {
      const res = await api.get<CronJob>(`/cronbuilder/jobs/${id}`);
      setEditJob(res.data!);
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.cronLogs.title")}</h1>
        <p className="text-muted-foreground">{t("pages.cronLogs.subtitle")}</p>
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

      {isError && (
        <div className="rounded-xl border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {getApiErrorMessage(error)}
          <Button size="sm" variant="outline" className="ms-3" onClick={() => refetch()}>
            {t("common.refresh")}
          </Button>
        </div>
      )}

      <DataTable
        loading={isLoading}
        data={jobs as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          { key: "table_name", header: t("common.table"), className: "ltr-field font-mono text-xs" },
          { key: "cron_expr", header: t("common.schedule"), className: "ltr-field" },
          {
            key: "is_active",
            header: t("common.status"),
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? t("common.active") : t("common.inactive")}
              </Badge>
            ),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex flex-wrap gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setSelectedJobId(row.id as string)}
                >
                  {t("common.viewLogs")}
                </Button>
                <Button
                  size="sm"
                  onClick={() => runMutation.mutate(row.id as string)}
                  disabled={runMutation.isPending}
                >
                  {t("common.runNow")}
                </Button>
                <Button size="sm" variant="secondary" onClick={() => openEdit(row.id as string)}>
                  {t("common.edit")}
                </Button>
                <Button size="sm" variant="destructive" onClick={() => setDeleteId(row.id as string)}>
                  {t("common.delete")}
                </Button>
              </div>
            ),
          },
        ]}
      />

      {selectedJobId && (
        <Card>
          <CardHeader>
            <CardTitle>{t("common.executionLogs")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DataTable
              loading={logsLoading}
              data={logs as unknown as Record<string, unknown>[]}
              keyField="id"
              columns={[
                {
                  key: "started_at",
                  header: t("common.started"),
                  render: (row) => formatJalaliDateTime(row.started_at as string),
                },
                {
                  key: "finished_at",
                  header: t("common.finished"),
                  render: (row) => formatJalaliDateTime(row.finished_at as string),
                },
                {
                  key: "success",
                  header: t("common.success"),
                  render: (row) => (
                    <Badge variant={row.success ? "success" : "destructive"}>
                      {row.success ? t("common.yes") : t("common.no")}
                    </Badge>
                  ),
                },
                { key: "affected_rows", header: t("common.affectedRows") },
                {
                  key: "error_message",
                  header: t("common.error"),
                  render: (row) => String(row.error_message || "—"),
                },
              ]}
            />
          </CardContent>
        </Card>
      )}

      {editJob && (
        <Card>
          <CardHeader>
            <CardTitle>{t("common.editJob")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.name")}</Label>
              <Input
                value={editJob.name}
                onChange={(e) => setEditJob({ ...editJob, name: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("common.cronExpression")}</Label>
              <Input
                className="ltr-field"
                value={editJob.cron_expr}
                onChange={(e) => setEditJob({ ...editJob, cron_expr: e.target.value })}
              />
            </div>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={editJob.is_active}
                onChange={(e) => setEditJob({ ...editJob, is_active: e.target.checked })}
                className="h-4 w-4"
              />
              {t("common.active")}
            </label>
            <div className="flex gap-2">
              <Button onClick={() => updateMutation.mutate(editJob)} disabled={updateMutation.isPending}>
                {updateMutation.isPending ? t("common.saving") : t("common.save")}
              </Button>
              <Button variant="ghost" onClick={() => setEditJob(null)}>
                {t("common.cancel")}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
            <DialogDescription>{t("common.confirmDeleteMessage")}</DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
            >
              {deleteMutation.isPending ? t("common.saving") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
