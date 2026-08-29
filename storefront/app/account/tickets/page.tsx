"use client";

import Link from "next/link";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t, statusLabel } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Ticket {
  id: string;
  subject: string;
  status: string;
  priority: string;
  created_at: string;
}

interface TicketDetail {
  ticket: Ticket;
  messages: { id: string; body: string; created_at: string }[];
}

export default function TicketsPage() {
  const queryClient = useQueryClient();
  const [subject, setSubject] = useState("");
  const [priority, setPriority] = useState("normal");
  const [departmentId, setDepartmentId] = useState("");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [replyBody, setReplyBody] = useState("");
  const [error, setError] = useState("");

  const { data: tickets = [], isLoading } = useQuery({
    queryKey: queryKeys.tickets,
    queryFn: async () => {
      const res = await api.get<Ticket[]>("/support/tickets");
      return res.data ?? [];
    },
  });

  const { data: ticketDetail } = useQuery({
    queryKey: ["tickets", selectedId],
    enabled: !!selectedId,
    queryFn: async () => {
      const res = await api.get<TicketDetail>(`/support/tickets/${selectedId}`);
      return res.data!;
    },
  });

  const createTicket = useMutation({
    mutationFn: async () => {
      await api.post("/support/tickets", {
        subject,
        priority,
        department_id: departmentId || undefined,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      setSubject("");
      setPriority("normal");
      setDepartmentId("");
      setError("");
    },
    onError: (err: Error) => setError(err.message),
  });

  const replyTicket = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${selectedId}/reply`, {
        body: replyBody,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tickets", selectedId] });
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      setReplyBody("");
    },
    onError: (err: Error) => setError(err.message),
  });

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">{t("tickets.title")}</h1>
          <p className="text-muted-foreground">{t("tickets.subtitle")}</p>
        </div>
        <Button asChild variant="outline">
          <Link href="/account/orders">{t("tickets.backOrders")}</Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("tickets.create")}</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            className="space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              createTicket.mutate();
            }}
          >
            <div className="space-y-2">
              <Label htmlFor="subject">{t("tickets.subject")}</Label>
              <Input
                id="subject"
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="priority">{t("tickets.priority")}</Label>
              <select
                id="priority"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
              >
                <option value="low">{t("tickets.priorityLow")}</option>
                <option value="normal">{t("tickets.priorityNormal")}</option>
                <option value="high">{t("tickets.priorityHigh")}</option>
                <option value="urgent">{t("tickets.priorityUrgent")}</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="department">{t("tickets.department")}</Label>
              <Input
                id="department"
                value={departmentId}
                onChange={(e) => setDepartmentId(e.target.value)}
                placeholder={t("common.optional")}
                className="ltr-field"
              />
            </div>
            {error && <p className="text-sm text-destructive">{error}</p>}
            <Button type="submit" disabled={createTicket.isPending}>
              {createTicket.isPending ? t("tickets.creating") : t("tickets.createBtn")}
            </Button>
          </form>
        </CardContent>
      </Card>

      <div className="space-y-4">
        <h2 className="text-xl font-semibold">{t("tickets.yourTickets")}</h2>
        {isLoading ? (
          <p className="text-muted-foreground">{t("tickets.loading")}</p>
        ) : (
          <>
            <DataTable
              rows={tickets}
              rowKey={(tk) => tk.id}
              emptyMessage={t("tickets.empty")}
              columns={[
                {
                  key: "subject",
                  header: t("tickets.subject"),
                  cell: (tk) => (
                    <button
                      type="button"
                      className="text-start font-medium hover:text-primary"
                      onClick={() => setSelectedId(tk.id)}
                    >
                      {tk.subject}
                    </button>
                  ),
                },
                {
                  key: "status",
                  header: t("orders.status"),
                  cell: (tk) => statusLabel(tk.status),
                },
                {
                  key: "priority",
                  header: t("tickets.priority"),
                  cell: (tk) => statusLabel(tk.priority),
                },
                {
                  key: "date",
                  header: t("orders.date"),
                  cell: (tk) => formatJalaliDateTime(tk.created_at),
                },
              ]}
            />
            {selectedId && ticketDetail && (
              <Card>
                <CardHeader>
                  <CardTitle>{ticketDetail.ticket.subject}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="max-h-64 space-y-3 overflow-y-auto">
                    {ticketDetail.messages?.map((msg) => (
                      <div
                        key={msg.id}
                        className="rounded-md border border-border p-3 text-sm"
                      >
                        <p className="leading-7">{msg.body}</p>
                        <p className="mt-1 text-xs text-muted-foreground">
                          {formatJalaliDateTime(msg.created_at)}
                        </p>
                      </div>
                    ))}
                  </div>
                  <form
                    className="space-y-2"
                    onSubmit={(e) => {
                      e.preventDefault();
                      replyTicket.mutate();
                    }}
                  >
                    <Label htmlFor="reply">{t("tickets.reply")}</Label>
                    <textarea
                      id="reply"
                      value={replyBody}
                      onChange={(e) => setReplyBody(e.target.value)}
                      required
                      rows={3}
                      className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm"
                    />
                    <Button type="submit" disabled={replyTicket.isPending}>
                      {replyTicket.isPending
                        ? t("tickets.sending")
                        : t("tickets.sendReply")}
                    </Button>
                  </form>
                </CardContent>
              </Card>
            )}
          </>
        )}
      </div>
    </div>
  );
}
