"use client";

import { useMemo, useState } from "react";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t, statusLabel } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { parseTicketDetail, TicketMessage } from "@/lib/support-api";
import { useTicketWebSocket } from "@/lib/use-ticket-websocket";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface UserOption {
  id: string;
  email: string;
}

export default function TicketDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const queryClient = useQueryClient();
  const [reply, setReply] = useState("");
  const [liveMessages, setLiveMessages] = useState<TicketMessage[]>([]);
  const [assignAgentId, setAssignAgentId] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: ticketDetail, isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.ticket(id),
    queryFn: async () => {
      const res = await api.get(`/support/tickets/${id}`);
      return parseTicketDetail(res.data);
    },
  });

  const { data: agents = [] } = useQuery({
    queryKey: ["users", "agents"],
    queryFn: async () => {
      const res = await api.get<UserOption[]>("/users?page=1&per_page=100");
      return res.data ?? [];
    },
  });

  useTicketWebSocket(id, (message) => {
    setLiveMessages((prev) => {
      if (prev.some((m) => m.id === message.id)) return prev;
      return [...prev, message];
    });
  });

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

  const replyMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${id}/reply`, { body: reply, is_internal: false });
    },
    onSuccess: () => {
      setReply("");
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const statusMutation = useMutation({
    mutationFn: async (status: string) => {
      await api.put(`/support/tickets/${id}/status`, { status });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const assignMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${id}/assign`, { agent_id: assignAgentId });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const closeMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${id}/close`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  if (isLoading) {
    return <p className="text-muted-foreground">{t("common.loading")}</p>;
  }

  if (isError || !ticketDetail) {
    return (
      <div className="rounded-xl border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
        {getApiErrorMessage(error)}
        <Button size="sm" variant="outline" className="ms-3" onClick={() => refetch()}>
          {t("common.refresh")}
        </Button>
      </div>
    );
  }

  const ticket = ticketDetail.ticket;

  return (
    <div className="space-y-6">
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

      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{ticket.subject}</h1>
          <div className="mt-2 flex gap-2">
            <Badge>{statusLabel(ticket.status)}</Badge>
            <Badge variant="secondary">{ticket.priority}</Badge>
          </div>
        </div>
        <div className="flex flex-wrap gap-2">
          <Select value={ticket.status} onValueChange={(v) => statusMutation.mutate(v)}>
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
            variant="outline"
            disabled={!assignAgentId || assignMutation.isPending}
            onClick={() => assignMutation.mutate()}
          >
            {t("common.assign")}
          </Button>
          <Button variant="destructive" onClick={() => closeMutation.mutate()}>
            {t("common.close")}
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.ticketDetail.messages")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {allMessages.map((msg) => (
            <div key={msg.id} className="rounded-md border border-border p-3 text-sm">
              <p>{msg.body}</p>
              <p className="mt-1 text-xs text-muted-foreground">
                {formatJalaliDateTime(msg.created_at)}
                {msg.is_internal && ` (${t("common.internal")})`}
              </p>
            </div>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.ticketDetail.reply")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>{t("common.message")}</Label>
            <Textarea value={reply} onChange={(e) => setReply(e.target.value)} />
          </div>
          <Button
            onClick={() => replyMutation.mutate()}
            disabled={replyMutation.isPending || !reply}
          >
            {replyMutation.isPending ? t("common.saving") : t("common.sendReply")}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
