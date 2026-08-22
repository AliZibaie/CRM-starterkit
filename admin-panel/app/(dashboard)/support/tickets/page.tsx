"use client";

import { useEffect, useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
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

interface Ticket {
  id: string;
  subject: string;
  status: string;
  priority: string;
  created_at: string;
}

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

interface Department {
  id: string;
  name: string;
}

const WS_BASE =
  process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1";

const fields: FormField[] = [
  { name: "subject", label: "Subject", type: "text", required: true },
  {
    name: "priority",
    label: "Priority",
    type: "select",
    options: [
      { label: "Low", value: "low" },
      { label: "Medium", value: "medium" },
      { label: "High", value: "high" },
    ],
  },
  { name: "department_id", label: "Department", type: "select", options: [] },
];

export default function SupportTicketsPage() {
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);
  const [showForm, setShowForm] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [reply, setReply] = useState("");
  const [wsMessages, setWsMessages] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const [values, setValues] = useState<Record<string, unknown>>({
    subject: "",
    priority: "medium",
    department_id: "",
  });

  const { data: departments = [] } = useQuery({
    queryKey: queryKeys.departments,
    queryFn: async () => {
      const res = await api.get<Department[]>("/support/departments");
      return res.data ?? [];
    },
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.tickets,
    queryFn: async () => {
      const res = await api.get<Ticket[]>("/support/tickets");
      return res.data ?? [];
    },
  });

  const { data: ticketDetail, isLoading: detailLoading } = useQuery({
    queryKey: queryKeys.ticket(selectedId ?? ""),
    queryFn: async () => {
      const res = await api.get<TicketDetail>(`/support/tickets/${selectedId}`);
      return res.data!;
    },
    enabled: !!selectedId,
  });

  useEffect(() => {
    if (!selectedId || !accessToken) return;
    const ws = new WebSocket(
      `${WS_BASE}/support/ws/${selectedId}?token=${accessToken}`
    );
    wsRef.current = ws;
    ws.onmessage = (event) => {
      setWsMessages((prev) => [...prev, event.data]);
    };
    return () => {
      ws.close();
      setWsMessages([]);
    };
  }, [selectedId, accessToken]);

  const formFields = fields.map((f) =>
    f.name === "department_id"
      ? {
          ...f,
          options: departments.map((d) => ({ label: d.name, value: d.id })),
        }
      : f
  );

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/support/tickets", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      setShowForm(false);
    },
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
      queryClient.invalidateQueries({
        queryKey: queryKeys.ticket(selectedId ?? ""),
      });
    },
  });

  const statusMutation = useMutation({
    mutationFn: async (status: string) => {
      await api.put(`/support/tickets/${selectedId}/status`, { status });
    },
    onSuccess: () =>
      queryClient.invalidateQueries({
        queryKey: queryKeys.ticket(selectedId ?? ""),
      }),
  });

  const closeMutation = useMutation({
    mutationFn: async () => {
      await api.post(`/support/tickets/${selectedId}/close`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets });
      queryClient.invalidateQueries({
        queryKey: queryKeys.ticket(selectedId ?? ""),
      });
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Support Tickets</h1>
          <p className="text-muted-foreground">Manage customer support tickets</p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "New Ticket"}
        </Button>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>Create Ticket</CardTitle>
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
            { key: "subject", header: "Subject" },
            {
              key: "status",
              header: "Status",
              render: (row) => (
                <Badge variant="outline">{String(row.status)}</Badge>
              ),
            },
            { key: "priority", header: "Priority" },
            { key: "created_at", header: "Created" },
            {
              key: "actions",
              header: "Actions",
              render: (row) => (
                <Button
                  size="sm"
                  variant={selectedId === row.id ? "default" : "outline"}
                  onClick={() => setSelectedId(row.id as string)}
                >
                  View
                </Button>
              ),
            },
          ]}
        />

        {selectedId && (
          <Card>
            <CardHeader>
              <div className="flex items-start justify-between gap-2">
                <CardTitle className="text-base">
                  {detailLoading ? "Loading..." : ticketDetail?.subject}
                </CardTitle>
                <div className="flex gap-2">
                  <Select
                    value={ticketDetail?.status}
                    onValueChange={(v) => statusMutation.mutate(v)}
                  >
                    <SelectTrigger className="w-32">
                      <SelectValue placeholder="Status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="open">Open</SelectItem>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="resolved">Resolved</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={() => closeMutation.mutate()}
                  >
                    Close
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="max-h-64 space-y-2 overflow-y-auto">
                {(ticketDetail?.messages ?? []).map((msg) => (
                  <div
                    key={msg.id}
                    className="rounded-md border border-border p-3 text-sm"
                  >
                    <p>{msg.body}</p>
                    <p className="mt-1 text-xs text-muted-foreground">
                      {msg.created_at}
                      {msg.is_internal && " (internal)"}
                    </p>
                  </div>
                ))}
                {wsMessages.map((msg, i) => (
                  <div
                    key={`ws-${i}`}
                    className="rounded-md border border-primary/30 bg-primary/5 p-3 text-sm"
                  >
                    <p className="text-xs text-muted-foreground">Live</p>
                    <p>{msg}</p>
                  </div>
                ))}
              </div>
              <div className="space-y-2">
                <Label>Reply</Label>
                <Textarea
                  value={reply}
                  onChange={(e) => setReply(e.target.value)}
                />
                <Button
                  onClick={() => replyMutation.mutate()}
                  disabled={replyMutation.isPending || !reply}
                >
                  Send Reply
                </Button>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
