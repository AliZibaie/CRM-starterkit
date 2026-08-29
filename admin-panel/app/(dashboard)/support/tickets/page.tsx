"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t, statusLabel } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { parseTicketDetail, TicketMessage, TicketRecord } from "@/lib/support-api";
import { useTicketWebSocket } from "@/lib/use-ticket-websocket";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Department {
  id: string;
  name: string;
}

interface UserOption {
  id: string;
  email: string;
}

export default function SupportTicketsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [reply, setReply] = useState("");
  const [liveMessages, setLiveMessages] = useState<TicketMessage[]>([]);
  const [statusFilter, setStatusFilter] = useState("");
  const [deptFilter, setDeptFilter] = useState("");
  const [priorityFilter, setPriorityFilter] = useState("");
  const [assignAgentId, setAssignAgentId] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );
  const [values, setValues] = useState<Record<string, unknown>>({
    subject: "",
    priority: "normal",
    department_id: "",
  });

  const ticketQueryKey = useMemo(() => {
    const params = new URLSearchParams();
    if (statusFilter) params.set("status", statusFilter);
    if (deptFilter) params.set("department_id", deptFilter);
    if (priorityFilter) params.set("priority", priorityFilter);
    const qs = params.toString();
    return [...queryKeys.tickets, qs] as const;
  }, [statusFilter, deptFilter, priorityFilter]);

  const { data: departments = [] } = useQuery({
    queryKey: queryKeys.departments,
    queryFn: async () => {
      const res = await api.get<Department[]>("/support/departments");
      return res.data ?? [];
    },
  });

  const { data: agents = [] } = useQuery({
    queryKey: ["users", "agents"],
    queryFn: async () => {
      const res = await api.get<UserOption[]>("/users?page=1&per_page=100");
      return res.data ?? [];
    },
  });

  const { data = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: ticketQueryKey,
    queryFn: async () => {
      const params = new URLSearchParams();
      if (statusFilter) params.set("status", statusFilter);
      if (deptFilter) params.set("department_id", deptFilter);
      if (priorityFilter) params.set("priority", priorityFilter);
      const qs = params.toString();
      const res = await api.get<TicketRecord[]>(`/support/tickets${qs ? `?${qs}` : ""}`);
      return res.data ?? [];
    },
  });

  const { data: ticketDetail, isLoading: detailLoading } = useQuery({
    queryKey: queryKeys.ticket(selectedId ?? ""),
    queryFn: async () => {
      const res = await api.get(`/support/tickets/${selectedId}`);
      return parseTicketDetail(res.data);
    },
    enabled: !!selectedId,
  });

  useTicketWebSocket(selectedId, (message) => {
    setLiveMessages((prev) => {
      if (prev.some((m) => m.id === message.id)) return prev;
      return [...prev, message];
    });
  });

  const deptMap = Object.fromEntries(departments.map((d) => [d.id, d.name]));

  const formFields: FormField[] = [
    { name: "subject", label: t("common.subject"), type: "text", required: true },
    {
      name: "priority",
      label: t("common.priority"),
      type: "select",
      options: [
        { label: t("common.low"), value: "low" },
        { label: t("common.normal"), value: "normal" },
        { label: t("common.medium"), value: "medium" },
        { label: t("common.high"), value: "high" },
      ],
    },
    {
      name: "department_id",
      label: t("common.department"),
      type: "select",
      options: departments.map((d) => ({ label: d.name, value: d.id })),
    },
  ];

  const allMessages = useMemo(() => {
    const base = ticketDetail?.messages ?? [];
    const merged = [...base];
    for (const m of liveMessages) {
      if (!merged.some((x) => x.id === m.id)) merged.push(m);
    }
    return merged.sort(
      (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );
  }, [ticketDetail?.messages, liveMessages]);

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/support/tickets", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      setShowForm(false);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const replyMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${selectedId}/reply`, {
        body: reply,
        is_internal: false,
      });
    },
    onSuccess: () => {
      setReply("");
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(selectedId ?? "") });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const statusMutation = useMutation({
    mutationFn: async (status: string) => {
      await api.put(`/support/tickets/${selectedId}/status`, { status });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(selectedId ?? "") });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const assignMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${selectedId}/assign`, { agent_id: assignAgentId });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(selectedId ?? "") });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const closeMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${selectedId}/close`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(selectedId ?? "") });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const selectTicket = (id: string) => {
    setSelectedId(id);
    setLiveMessages([]);
    setAssignAgentId("");
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.tickets.title")}</h1>
          <p className="text-muted-foreground">{t("pages.tickets.subtitle")}</p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? t("common.cancel") : t("pages.tickets.newTicket")}
        </Button>
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

      <div className="grid gap-3 md:grid-cols-3">
        <Select value={statusFilter || "all"} onValueChange={(v) => setStatusFilter(v === "all" ? "" : v)}>
          <SelectTrigger>
            <SelectValue placeholder={t("common.status")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("common.allStatuses")}</SelectItem>
            <SelectItem value="open">{t("status.open")}</SelectItem>
            <SelectItem value="in_progress">{t("status.in_progress")}</SelectItem>
            <SelectItem value="resolved">{t("status.resolved")}</SelectItem>
            <SelectItem value="closed">{t("status.closed")}</SelectItem>
          </SelectContent>
        </Select>
        <Select value={deptFilter || "all"} onValueChange={(v) => setDeptFilter(v === "all" ? "" : v)}>
          <SelectTrigger>
            <SelectValue placeholder={t("common.department")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("common.allDepartments")}</SelectItem>
            {departments.map((d) => (
              <SelectItem key={d.id} value={d.id}>
                {d.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select
          value={priorityFilter || "all"}
          onValueChange={(v) => setPriorityFilter(v === "all" ? "" : v)}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("common.priority")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("common.allPriorities")}</SelectItem>
            <SelectItem value="low">{t("common.low")}</SelectItem>
            <SelectItem value="normal">{t("common.normal")}</SelectItem>
            <SelectItem value="medium">{t("common.medium")}</SelectItem>
            <SelectItem value="high">{t("common.high")}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.tickets.create")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={formFields}
              values={values}
              onChange={(n, v) => setValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() => createMutation.mutate(values)}
              loading={createMutation.isPending}
            />
          </CardContent>
        </Card>
      )}

      <div className="grid gap-6 lg:grid-cols-2">
        <DataTable
          loading={isLoading}
          data={data as unknown as Record<string, unknown>[]}
          keyField="id"
          columns={[
            { key: "subject", header: t("common.subject") },
            {
              key: "status",
              header: t("common.status"),
              render: (row) => (
                <Badge variant="outline">{statusLabel(String(row.status))}</Badge>
              ),
            },
            { key: "priority", header: t("common.priority") },
            {
              key: "department_id",
              header: t("common.department"),
              render: (row) => deptMap[row.department_id as string] ?? "—",
            },
            {
              key: "created_at",
              header: t("common.createdAt"),
              render: (row) => formatJalaliDateTime(row.created_at as string),
            },
            {
              key: "actions",
              header: t("common.actions"),
              render: (row) => (
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant={selectedId === row.id ? "default" : "outline"}
                    onClick={() => selectTicket(row.id as string)}
                  >
                    {t("common.view")}
                  </Button>
                  <Button size="sm" variant="ghost" asChild>
                    <Link href={`/support/tickets/${row.id}`}>{t("common.open")}</Link>
                  </Button>
                </div>
              ),
            },
          ]}
        />

        {selectedId && ticketDetail && (
          <Card>
            <CardHeader>
              <div className="flex flex-col gap-3">
                <CardTitle className="text-base">
                  {detailLoading ? t("common.loading") : ticketDetail.ticket.subject}
                </CardTitle>
                <div className="flex flex-wrap gap-2">
                  <Select
                    value={ticketDetail.ticket.status}
                    onValueChange={(v) => statusMutation.mutate(v)}
                  >
                    <SelectTrigger className="w-40">
                      <SelectValue placeholder={t("common.status")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="open">{t("status.open")}</SelectItem>
                      <SelectItem value="in_progress">{t("status.in_progress")}</SelectItem>
                      <SelectItem value="resolved">{t("status.resolved")}</SelectItem>
                      <SelectItem value="closed">{t("status.closed")}</SelectItem>
                    </SelectContent>
                  </Select>
                  <Select value={assignAgentId} onValueChange={setAssignAgentId}>
                    <SelectTrigger className="w-48">
                      <SelectValue placeholder={t("common.assignAgent")} />
                    </SelectTrigger>
                    <SelectContent>
                      {agents.map((u) => (
                        <SelectItem key={u.id} value={u.id}>
                          {u.email}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={!assignAgentId || assignMutation.isPending}
                    onClick={() => assignMutation.mutate()}
                  >
                    {t("common.assign")}
                  </Button>
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={() => closeMutation.mutate()}
                    disabled={closeMutation.isPending}
                  >
                    {t("common.close")}
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="max-h-64 space-y-2 overflow-y-auto">
                {allMessages.map((msg) => (
                  <div
                    key={msg.id}
                    className="rounded-md border border-border p-3 text-sm"
                  >
                    <p>{msg.body}</p>
                    <p className="mt-1 text-xs text-muted-foreground">
                      {formatJalaliDateTime(msg.created_at)}
                      {msg.is_internal && ` (${t("common.internal")})`}
                    </p>
                  </div>
                ))}
              </div>
              <div className="space-y-2">
                <Label>{t("common.reply")}</Label>
                <Textarea value={reply} onChange={(e) => setReply(e.target.value)} />
                <Button
                  onClick={() => replyMutation.mutate()}
                  disabled={replyMutation.isPending || !reply}
                >
                  {replyMutation.isPending ? t("common.saving") : t("common.sendReply")}
                </Button>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
