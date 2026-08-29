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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Contact {
  id: string;
  name: string;
  email: string;
  phone: string;
  company?: string;
}

interface Activity {
  id: string;
  type: string;
  body: string;
  created_at: string;
}

interface Customer360 {
  contact: Contact;
  orders: Record<string, unknown>[];
  tickets: Record<string, unknown>[];
  activities: Activity[];
  deals: Record<string, unknown>[];
  tags: string[];
}

export default function ContactsPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [company, setCompany] = useState("");
  const [editContact, setEditContact] = useState<Contact | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [detailId, setDetailId] = useState<string | null>(null);
  const [activityContactId, setActivityContactId] = useState<string | null>(null);
  const [activityType, setActivityType] = useState("call");
  const [activityBody, setActivityBody] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.contacts,
    queryFn: async () => {
      const res = await api.get<Contact[]>("/crm/contacts");
      return res.data ?? [];
    },
  });

  const { data: activities = [] } = useQuery({
    queryKey: queryKeys.activities("contact", activityContactId ?? detailId ?? ""),
    queryFn: async () => {
      const id = activityContactId ?? detailId ?? "";
      const res = await api.get<Activity[]>(`/crm/activities?entity_type=contact&entity_id=${id}`);
      return res.data ?? [];
    },
    enabled: !!(activityContactId || detailId),
  });

  const { data: view360 } = useQuery({
    queryKey: queryKeys.contact(detailId ?? ""),
    queryFn: async () => {
      const res = await api.get<Customer360>(`/crm/customers/${detailId}/360`);
      return res.data!;
    },
    enabled: !!detailId,
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/crm/contacts", { name, email, phone, company });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.contacts });
      setName("");
      setEmail("");
      setPhone("");
      setCompany("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async (c: Contact) => {
      await api.put(`/crm/contacts/${c.id}`, {
        name: c.name,
        email: c.email,
        phone: c.phone,
        company: c.company ?? "",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.contacts });
      setEditContact(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/crm/contacts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.contacts });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const logActivityMutation = useMutation({
    mutationFn: async () => {
      if (!activityContactId) return;
      await api.post("/crm/activities", {
        entity_type: "contact",
        entity_id: activityContactId,
        type: activityType,
        body: activityBody,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.activities("contact", activityContactId ?? ""),
      });
      setActivityBody("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.contacts.title")}</h1>
        <p className="text-muted-foreground">{t("pages.contacts.subtitle")}</p>
      </div>

      {feedback && (
        <div
          className={`rounded-md border px-4 py-3 text-sm ${
            feedback.type === "error"
              ? "border-destructive/50 bg-destructive/10 text-destructive"
              : "border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400"
          }`}
        >
          {feedback.text}
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.contacts.create")}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-2">
          <Input placeholder={t("common.name")} value={name} onChange={(e) => setName(e.target.value)} />
          <Input placeholder={t("common.email")} value={email} onChange={(e) => setEmail(e.target.value)} />
          <Input placeholder={t("common.phone")} value={phone} onChange={(e) => setPhone(e.target.value)} />
          <Input placeholder={t("common.company")} value={company} onChange={(e) => setCompany(e.target.value)} />
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending || !name.trim()}>
            {t("common.create")}
          </Button>
        </CardContent>
      </Card>

      {activityContactId && (
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.contacts.logActivity")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Select value={activityType} onValueChange={setActivityType}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="call">{t("pages.contacts.activityCall")}</SelectItem>
                <SelectItem value="email">{t("common.email")}</SelectItem>
                <SelectItem value="meeting">{t("pages.contacts.activityMeeting")}</SelectItem>
                <SelectItem value="note">{t("common.notes")}</SelectItem>
              </SelectContent>
            </Select>
            <Input value={activityBody} onChange={(e) => setActivityBody(e.target.value)} />
            <div className="flex gap-2">
              <Button onClick={() => logActivityMutation.mutate()} disabled={logActivityMutation.isPending}>
                {t("common.save")}
              </Button>
              <Button variant="outline" onClick={() => setActivityContactId(null)}>
                {t("common.cancel")}
              </Button>
            </div>
            <ul className="space-y-2 text-sm">
              {activities.map((a) => (
                <li key={a.id} className="rounded border px-3 py-2">
                  <span className="font-medium">{a.type}</span> — {a.body}
                  <span className="ms-2 text-muted-foreground">{formatJalaliDateTime(a.created_at)}</span>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          { key: "email", header: t("common.email") },
          { key: "phone", header: t("common.phone") },
          { key: "company", header: t("common.company") },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex flex-wrap gap-2">
                <Button size="sm" variant="outline" onClick={() => setDetailId(row.id as string)}>
                  360
                </Button>
                <Button size="sm" variant="outline" onClick={() => setActivityContactId(row.id as string)}>
                  {t("pages.contacts.activity")}
                </Button>
                <Button size="sm" variant="outline" onClick={() => setEditContact(row as unknown as Contact)}>
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

      <Dialog open={!!editContact} onOpenChange={() => setEditContact(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.edit")}</DialogTitle>
          </DialogHeader>
          {editContact && (
            <div className="space-y-3">
              <Input value={editContact.name} onChange={(e) => setEditContact({ ...editContact, name: e.target.value })} />
              <Input value={editContact.email} onChange={(e) => setEditContact({ ...editContact, email: e.target.value })} />
              <Input value={editContact.phone} onChange={(e) => setEditContact({ ...editContact, phone: e.target.value })} />
              <Input
                value={editContact.company ?? ""}
                onChange={(e) => setEditContact({ ...editContact, company: e.target.value })}
              />
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditContact(null)}>
              {t("common.cancel")}
            </Button>
            <Button onClick={() => editContact && updateMutation.mutate(editContact)} disabled={updateMutation.isPending}>
              {t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!detailId} onOpenChange={() => setDetailId(null)}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("pages.contacts.customer360")}</DialogTitle>
          </DialogHeader>
          {view360 && (
            <div className="space-y-4 text-sm">
              <p className="font-medium">{view360.contact.name} — {view360.contact.email}</p>
              <div>
                <h4 className="font-medium">{t("pages.orders.title")} ({view360.orders.length})</h4>
                {view360.orders.length === 0 && <p className="text-muted-foreground">{t("common.noData")}</p>}
                {view360.orders.map((o, i) => (
                  <p key={i}>{String(o.id)} — {String(o.status)} — {String(o.total_amount)}</p>
                ))}
              </div>
              <div>
                <h4 className="font-medium">{t("nav.tickets")} ({view360.tickets.length})</h4>
                {view360.tickets.map((tk, i) => (
                  <p key={i}>{String(tk.subject)} — {String(tk.status)}</p>
                ))}
              </div>
              <div>
                <h4 className="font-medium">{t("pages.contacts.activity")} ({view360.activities.length})</h4>
                {view360.activities.map((a) => (
                  <p key={a.id}>{a.type}: {a.body}</p>
                ))}
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Dialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteId(null)}>{t("common.cancel")}</Button>
            <Button variant="destructive" onClick={() => deleteId && deleteMutation.mutate(deleteId)}>
              {t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
