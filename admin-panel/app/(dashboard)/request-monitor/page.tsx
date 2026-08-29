"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatGregorianDateTime } from "@/lib/date";
import { datetimeLocalToRFC3339 } from "@/lib/datetime-input";
import { DataTable } from "@/components/table/data-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface RequestLog {
  id: string;
  ip: string;
  method: string;
  path: string;
  status: number;
  user_id: string;
  duration_ms: number;
  country?: string;
  city?: string;
  suspicious: boolean;
  created_at: string;
}

interface RequestStats {
  total: number;
  avg_duration_ms: number;
  error_count: number;
  suspicious_count: number;
}

export default function RequestMonitorPage() {
  const [ip, setIp] = useState("");
  const [path, setPath] = useState("");
  const [status, setStatus] = useState("");
  const [userId, setUserId] = useState("");
  const [fromLocal, setFromLocal] = useState("");
  const [toLocal, setToLocal] = useState("");
  const [autoRefresh, setAutoRefresh] = useState(true);

  const filters = useMemo(
    () => ({ ip, path, status, user_id: userId, from: fromLocal, to: toLocal }),
    [ip, path, status, userId, fromLocal, toLocal]
  );

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    if (ip) params.set("ip", ip);
    if (path) params.set("path", path);
    if (status) params.set("status", status);
    if (userId) params.set("user_id", userId);
    const from = datetimeLocalToRFC3339(fromLocal);
    const to = datetimeLocalToRFC3339(toLocal);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    params.set("limit", "100");
    return params.toString();
  }, [ip, path, status, userId, fromLocal, toLocal]);

  const statsQueryString = useMemo(() => {
    const params = new URLSearchParams();
    if (ip) params.set("ip", ip);
    if (path) params.set("path", path);
    if (userId) params.set("user_id", userId);
    const from = datetimeLocalToRFC3339(fromLocal);
    const to = datetimeLocalToRFC3339(toLocal);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    return params.toString();
  }, [ip, path, userId, fromLocal, toLocal]);

  const { data = [], isLoading, refetch, isFetching } = useQuery({
    queryKey: queryKeys.requestLogs(filters),
    queryFn: async () => {
      const res = await api.get<RequestLog[]>(`/request-logs?${queryString}`);
      return res.data ?? [];
    },
    refetchInterval: autoRefresh ? 5000 : false,
  });

  const { data: stats } = useQuery({
    queryKey: queryKeys.requestStats(filters),
    queryFn: async () => {
      const res = await api.get<RequestStats>(`/request-logs/stats?${statsQueryString}`);
      return res.data;
    },
    refetchInterval: autoRefresh ? 5000 : false,
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.requestMonitor.title")}</h1>
          <p className="text-muted-foreground">
            {t("pages.requestMonitor.subtitle")}
            {isFetching && autoRefresh && t("common.refreshing")}
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant={autoRefresh ? "default" : "outline"} size="sm" onClick={() => setAutoRefresh(!autoRefresh)}>
            {autoRefresh ? t("common.autoRefreshOn") : t("common.autoRefreshOff")}
          </Button>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            {t("common.refreshNow")}
          </Button>
        </div>
      </div>

      {stats && (
        <div className="grid gap-4 md:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">{t("pages.requestMonitor.totalRequests")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.total}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">{t("pages.requestMonitor.avgDuration")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.avg_duration_ms} ms</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">{t("pages.requestMonitor.errors")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.error_count}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">{t("pages.requestMonitor.suspicious")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.suspicious_count}</p>
            </CardContent>
          </Card>
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("common.filters")}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
          <div className="space-y-2">
            <Label>{t("common.ip")}</Label>
            <Input value={ip} onChange={(e) => setIp(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.path")}</Label>
            <Input value={path} onChange={(e) => setPath(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.status")}</Label>
            <Input value={status} onChange={(e) => setStatus(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.userId")}</Label>
            <Input value={userId} onChange={(e) => setUserId(e.target.value)} />
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
          {
            key: "created_at",
            header: t("common.time"),
            render: (row) => formatGregorianDateTime(row.created_at as string),
          },
          { key: "ip", header: t("common.ip") },
          {
            key: "geo",
            header: t("pages.requestMonitor.geo"),
            render: (row) => {
              const country = row.country as string;
              const city = row.city as string;
              return country || city ? `${country || "—"} / ${city || "—"}` : "—";
            },
          },
          { key: "method", header: t("common.method") },
          { key: "path", header: t("common.path") },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => {
              const s = row.status as number;
              const variant = s >= 500 ? "destructive" : s >= 400 ? "secondary" : "success";
              return <Badge variant={variant}>{s}</Badge>;
            },
          },
          { key: "user_id", header: t("common.userId"), className: "font-mono text-xs" },
          { key: "duration_ms", header: t("common.durationMs") },
          {
            key: "suspicious",
            header: t("common.suspicious"),
            render: (row) =>
              row.suspicious ? <Badge variant="destructive">{t("common.yes")}</Badge> : t("common.no"),
          },
        ]}
      />
    </div>
  );
}
