"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface CronJob {
  id: string;
  name: string;
  table_name: string;
  cron_expr: string;
  is_active: boolean;
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

  const { data: jobs = [], isLoading } = useQuery({
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
    onSuccess: () => {
      if (selectedJobId) {
        queryClient.invalidateQueries({
          queryKey: queryKeys.cronLogs(selectedJobId),
        });
      }
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Cron Job Logs</h1>
        <p className="text-muted-foreground">View and run scheduled jobs</p>
      </div>

      <DataTable
        loading={isLoading}
        data={jobs as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: "Name" },
          { key: "table_name", header: "Table" },
          { key: "cron_expr", header: "Schedule" },
          {
            key: "is_active",
            header: "Status",
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? "Active" : "Inactive"}
              </Badge>
            ),
          },
          {
            key: "actions",
            header: "Actions",
            render: (row) => (
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setSelectedJobId(row.id as string)}
                >
                  View Logs
                </Button>
                <Button
                  size="sm"
                  onClick={() => runMutation.mutate(row.id as string)}
                  disabled={runMutation.isPending}
                >
                  Run Now
                </Button>
              </div>
            ),
          },
        ]}
      />

      {selectedJobId && (
        <Card>
          <CardHeader>
            <CardTitle>Execution Logs</CardTitle>
          </CardHeader>
          <CardContent>
            <DataTable
              loading={logsLoading}
              data={logs as unknown as Record<string, unknown>[]}
              keyField="id"
              columns={[
                { key: "started_at", header: "Started" },
                { key: "finished_at", header: "Finished" },
                {
                  key: "success",
                  header: "Success",
                  render: (row) => (
                    <Badge variant={row.success ? "success" : "destructive"}>
                      {row.success ? "Yes" : "No"}
                    </Badge>
                  ),
                },
                { key: "affected_rows", header: "Affected Rows" },
                { key: "error_message", header: "Error" },
              ]}
            />
          </CardContent>
        </Card>
      )}
    </div>
  );
}
