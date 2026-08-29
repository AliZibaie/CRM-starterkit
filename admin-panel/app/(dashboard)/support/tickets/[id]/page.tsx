"use client";

import { useEffect, useRef, useState } from "react";
import { useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { useAuthStore } from "@/stores/auth-store";
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

interface TicketDetail {
  id: string;
  subject: string;
  status: string;
  priority: string;
  messages?: {
    id: string;
    body: string;
    created_at: string;
    is_internal: boolean;
  }[];
}

const WS_BASE =
  process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1";

export default function TicketDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);
  const [reply, setReply] = useState("");
  const [wsMessages, setWsMessages] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  const { data: ticket, isLoading } = useQuery({
    queryKey: queryKeys.ticket(id),
    queryFn: async () => {
      const res = await api.get<TicketDetail>(`/support/tickets/${id}`);
      return res.data!;
    },
  });

  useEffect(() => {
    if (!id || !accessToken) return;
    const ws = new WebSocket(`${WS_BASE}/support/ws/${id}?token=${accessToken}`);
    wsRef.current = ws;
    ws.onmessage = (event) => {
      setWsMessages((prev) => [...prev, event.data]);
    };
    return () => ws.close();
  }, [id, accessToken]);

  const replyMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${id}/reply`, {
        body: reply,
        is_internal: false,
      });
    },
    onSuccess: () => {
      setReply("");
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) });
    },
  });

  const statusMutation = useMutation({
    mutationFn: async (status: string) => {
      await api.put(`/support/tickets/${id}/status`, { status });
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) }),
  });

  const closeMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${id}/close`);
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.ticket(id) }),
  });

  if (isLoading) {
    return <p className="text-muted-foreground">{t("common.loading")}</p>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold">{ticket?.subject}</h1>
          <div className="mt-2 flex gap-2">
            <Badge>{ticket?.status}</Badge>
            <Badge variant="secondary">{ticket?.priority}</Badge>
          </div>
        </div>
        <div className="flex gap-2">
          <Select
            value={ticket?.status}
            onValueChange={(v) => statusMutation.mutate(v)}
          >
            <SelectTrigger className="w-36">
              <SelectValue placeholder={t("common.status")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="open">{t("status.open")}</SelectItem>
              <SelectItem value="pending">{t("status.pending")}</SelectItem>
              <SelectItem value="resolved">{t("status.resolved")}</SelectItem>
            </SelectContent>
          </Select>
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
          {(ticket?.messages ?? []).map((msg) => (
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
          {wsMessages.map((msg, i) => (
            <div
              key={`ws-${i}`}
              className="rounded-md border border-primary/30 bg-primary/5 p-3 text-sm"
            >
              <p className="text-xs text-muted-foreground">{t("common.live")}</p>
              <p>{msg}</p>
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
            {t("common.sendReply")}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
