"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
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
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [autoRefresh, setAutoRefresh] = useState(true);

  const filters = useMemo(
    () => ({
      ip,
      path,
      status,
      user_id: userId,
      from,
      to,
    }),
    [ip, path, status, userId, from, to]
  );

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    if (ip) params.set("ip", ip);
    if (path) params.set("path", path);
    if (status) params.set("status", status);
    if (userId) params.set("user_id", userId);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    params.set("limit", "100");
    return params.toString();
  }, [ip, path, status, userId, from, to]);

  const statsQueryString = useMemo(() => {
    const params = new URLSearchParams();
    if (ip) params.set("ip", ip);
    if (path) params.set("path", path);
    if (userId) params.set("user_id", userId);
    if (from) params.set("from", from);
    if (to) params.set("to", to);
    return params.toString();
  }, [ip, path, userId, from, to]);

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
      const res = await api.get<RequestStats>(
        `/request-logs/stats?${statsQueryString}`
      );
      return res.data;
    },
    refetchInterval: autoRefresh ? 5000 : false,
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Request Monitor</h1>
          <p className="text-muted-foreground">
            Live HTTP request log with filters
            {isFetching && autoRefresh && " (refreshing...)"}
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant={autoRefresh ? "default" : "outline"}
            size="sm"
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            Auto-refresh {autoRefresh ? "ON" : "OFF"}
          </Button>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Refresh Now
          </Button>
        </div>
      </div>

      {stats && (
        <div className="grid gap-4 md:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Total Requests</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.total}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Avg Duration</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.avg_duration_ms} ms</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Errors</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.error_count}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Suspicious</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{stats.suspicious_count}</p>
            </CardContent>
          </Card>
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
          <div className="space-y-2">
            <Label>IP</Label>
            <Input value={ip} onChange={(e) => setIp(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Path</Label>
            <Input value={path} onChange={(e) => setPath(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Status</Label>
            <Input value={status} onChange={(e) => setStatus(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>User ID</Label>
            <Input value={userId} onChange={(e) => setUserId(e.target.value)} />
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
          { key: "created_at", header: "Time" },
          { key: "ip", header: "IP" },
          { key: "method", header: "Method" },
          { key: "path", header: "Path" },
          {
            key: "status",
            header: "Status",
            render: (row) => {
              const s = row.status as number;
              const variant =
                s >= 500 ? "destructive" : s >= 400 ? "secondary" : "success";
              return <Badge variant={variant}>{s}</Badge>;
            },
          },
          { key: "user_id", header: "User" },
          { key: "duration_ms", header: "Duration (ms)" },
          {
            key: "suspicious",
            header: "Suspicious",
            render: (row) =>
              row.suspicious ? (
                <Badge variant="destructive">Yes</Badge>
              ) : (
                "No"
              ),
          },
        ]}
      />
    </div>
  );
}
