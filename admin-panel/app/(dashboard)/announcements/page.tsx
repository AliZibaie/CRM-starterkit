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
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Announcement {
  id: string;
  title: string;
  body: string;
  audience: string;
  audience_ref?: string;
  status: string;
  created_at: string;
  broadcast_at?: string;
}

interface Role {
  id: string;
  name: string;
}

function audienceLabel(a: Announcement) {
  if (a.audience === "role") return `${t("audience.role")}: ${a.audience_ref || "—"}`;
  if (a.audience === "group") return `${t("audience.group")}: ${a.audience_ref || "—"}`;
  return t("audience.all");
}

function statusBadge(status: string) {
  const key = `announcementStatus.${status}`;
  const label = t(key);
  const text = label === key ? status : label;
  const variant =
    status === "broadcasted" ? "success" : status === "draft" ? "secondary" : "outline";
  return <Badge variant={variant}>{text}</Badge>;
}

export default function AnnouncementsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [values, setValues] = useState({
    title: "",
    body: "",
    audience: "all",
    audience_ref: "",
  });
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.announcements,
    queryFn: async () => {
      const res = await api.get<Announcement[]>("/announcements");
      return res.data ?? [];
    },
  });

  const { data: roles = [] } = useQuery({
    queryKey: queryKeys.roles,
    queryFn: async () => {
      const res = await api.get<Role[]>("/rbac/roles");
      return res.data ?? [];
    },
  });

  const resetForm = () => {
    setValues({ title: "", body: "", audience: "all", audience_ref: "" });
    setEditId(null);
  };

  const payload = () => ({
    title: values.title,
    body: values.body,
    audience: values.audience,
    audience_ref: values.audience === "all" ? "" : values.audience_ref,
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/announcements", payload());
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.announcements });
      setShowForm(false);
      resetForm();
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/announcements/${editId}`, payload());
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.announcements });
      setShowForm(false);
      resetForm();
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/announcements/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.announcements });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const broadcastMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/announcements/${id}/broadcast`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.announcements });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const startEdit = (row: Announcement) => {
    setEditId(row.id);
    setShowForm(true);
    setValues({
      title: row.title,
      body: row.body,
      audience: row.audience || "all",
      audience_ref: row.audience_ref || "",
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.announcements.title")}</h1>
          <p className="text-muted-foreground">{t("pages.announcements.subtitle")}</p>
        </div>
        <Button
          onClick={() => {
            if (showForm) {
              setShowForm(false);
              resetForm();
            } else {
              setShowForm(true);
            }
          }}
        >
          {showForm ? t("common.cancel") : t("pages.announcements.newAnnouncement")}
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

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>
              {editId ? t("common.edit") : t("pages.announcements.create")}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label>{t("common.title")}</Label>
              <Input
                value={values.title}
                onChange={(e) => setValues((p) => ({ ...p, title: e.target.value }))}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("common.body")}</Label>
              <Textarea
                value={values.body}
                onChange={(e) => setValues((p) => ({ ...p, body: e.target.value }))}
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label>{t("common.audience")}</Label>
                <Select
                  value={values.audience}
                  onValueChange={(v) => setValues((p) => ({ ...p, audience: v, audience_ref: "" }))}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">{t("audience.all")}</SelectItem>
                    <SelectItem value="role">{t("audience.role")}</SelectItem>
                    <SelectItem value="group">{t("audience.group")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              {values.audience === "role" && (
                <div className="space-y-2">
                  <Label>{t("audience.refRole")}</Label>
                  <Select
                    value={values.audience_ref}
                    onValueChange={(v) => setValues((p) => ({ ...p, audience_ref: v }))}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("common.selectRole")} />
                    </SelectTrigger>
                    <SelectContent>
                      {roles.map((r) => (
                        <SelectItem key={r.id} value={r.id}>
                          {r.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}
              {values.audience === "group" && (
                <div className="space-y-2">
                  <Label>{t("audience.refGroup")}</Label>
                  <Input
                    value={values.audience_ref}
                    onChange={(e) => setValues((p) => ({ ...p, audience_ref: e.target.value }))}
                  />
                </div>
              )}
            </div>
            <Button
              onClick={() => (editId ? updateMutation.mutate() : createMutation.mutate())}
              disabled={
                createMutation.isPending ||
                updateMutation.isPending ||
                !values.title ||
                !values.body
              }
            >
              {createMutation.isPending || updateMutation.isPending
                ? t("common.saving")
                : editId
                  ? t("common.save")
                  : t("common.create")}
            </Button>
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: t("common.title") },
          {
            key: "audience",
            header: t("common.audience"),
            render: (row) => audienceLabel(row as unknown as Announcement),
          },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => statusBadge(String(row.status)),
          },
          {
            key: "created_at",
            header: t("common.createdAt"),
            render: (row) => formatJalaliDateTime(row.created_at as string),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => {
              const item = row as unknown as Announcement;
              return (
                <div className="flex flex-wrap gap-2">
                  {item.status !== "broadcasted" && (
                    <>
                      <Button size="sm" variant="outline" onClick={() => startEdit(item)}>
                        {t("common.edit")}
                      </Button>
                      <Button
                        size="sm"
                        onClick={() => broadcastMutation.mutate(item.id)}
                        disabled={broadcastMutation.isPending}
                      >
                        {t("common.broadcast")}
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => setDeleteId(item.id)}
                      >
                        {t("common.delete")}
                      </Button>
                    </>
                  )}
                </div>
              );
            },
          },
        ]}
      />

      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
            <DialogDescription>{t("common.confirmDeleteMessage")}</DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
            >
              {deleteMutation.isPending ? t("common.saving") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
