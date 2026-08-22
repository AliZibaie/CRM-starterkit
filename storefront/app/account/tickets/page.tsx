"use client";

import Link from "next/link";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
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
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Support Tickets</h1>
          <p className="text-muted-foreground">Get help from our team</p>
        </div>
        <Button asChild variant="outline">
          <Link href="/account/orders">Back to Orders</Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Ticket</CardTitle>
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
              <Label htmlFor="subject">Subject</Label>
              <Input
                id="subject"
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="priority">Priority</Label>
              <select
                id="priority"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
              >
                <option value="low">Low</option>
                <option value="normal">Normal</option>
                <option value="high">High</option>
                <option value="urgent">Urgent</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="department">Department ID</Label>
              <Input
                id="department"
                value={departmentId}
                onChange={(e) => setDepartmentId(e.target.value)}
                placeholder="Optional"
              />
            </div>
            {error && <p className="text-sm text-destructive">{error}</p>}
            <Button type="submit" disabled={createTicket.isPending}>
              {createTicket.isPending ? "Creating..." : "Create Ticket"}
            </Button>
          </form>
        </CardContent>
      </Card>

      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Your Tickets</h2>
        {isLoading ? (
          <p className="text-muted-foreground">Loading tickets...</p>
        ) : tickets.length === 0 ? (
          <p className="text-muted-foreground">No tickets yet.</p>
        ) : (
          <div className="grid gap-6 lg:grid-cols-2">
            <div className="space-y-2">
              {tickets.map((ticket) => (
                <Card
                  key={ticket.id}
                  className={`cursor-pointer transition-colors ${
                    selectedId === ticket.id ? "border-primary" : ""
                  }`}
                  onClick={() => setSelectedId(ticket.id)}
                >
                  <CardContent className="py-4">
                    <p className="font-medium">{ticket.subject}</p>
                    <div className="mt-1 flex gap-2 text-xs text-muted-foreground">
                      <span className="capitalize">{ticket.status}</span>
                      <span>·</span>
                      <span className="capitalize">{ticket.priority}</span>
                      <span>·</span>
                      <span>
                        {new Date(ticket.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
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
                        className="rounded-md border p-3 text-sm"
                      >
                        <p>{msg.body}</p>
                        <p className="mt-1 text-xs text-muted-foreground">
                          {new Date(msg.created_at).toLocaleString()}
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
                    <Label htmlFor="reply">Reply</Label>
                    <textarea
                      id="reply"
                      value={replyBody}
                      onChange={(e) => setReplyBody(e.target.value)}
                      required
                      rows={3}
                      className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm"
                    />
                    <Button type="submit" disabled={replyTicket.isPending}>
                      {replyTicket.isPending ? "Sending..." : "Send Reply"}
                    </Button>
                  </form>
                </CardContent>
              </Card>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
