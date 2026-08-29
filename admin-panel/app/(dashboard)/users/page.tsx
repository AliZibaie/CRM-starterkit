"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { PaginationBar } from "@/components/table/pagination-bar";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  role_ids?: string[];
  created_at: string;
}

interface Role {
  id: string;
  name: string;
}

const PER_PAGE = 20;

const createFields: FormField[] = [
  { name: "email", label: t("common.email"), type: "email", required: true },
  { name: "password", label: t("common.password"), type: "password", required: true },
  { name: "first_name", label: t("common.firstName"), type: "text" },
  { name: "last_name", label: t("common.lastName"), type: "text" },
];

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [formValues, setFormValues] = useState<Record<string, unknown>>({
    email: "",
    password: "",
    first_name: "",
    last_name: "",
  });
  const [createRoleIds, setCreateRoleIds] = useState<string[]>([]);
  const [editId, setEditId] = useState<string | null>(null);
  const [editValues, setEditValues] = useState<Record<string, unknown>>({});
  const [editRoleIds, setEditRoleIds] = useState<string[]>([]);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: roles = [] } = useQuery({
    queryKey: queryKeys.roles,
    queryFn: async () => {
      const res = await api.get<Role[]>("/rbac/roles");
      return res.data ?? [];
    },
  });

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.users(page),
    queryFn: async () => {
      const res = await api.get<User[]>(`/users?page=${page}&per_page=${PER_PAGE}`);
      return {
        users: res.data ?? [],
        total: Number(res.meta?.total ?? 0),
      };
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/users", {
        ...payload,
        role_ids: createRoleIds,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setShowForm(false);
      setFormValues({ email: "", password: "", first_name: "", last_name: "" });
      setCreateRoleIds([]);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      payload,
      roleIds,
    }: {
      id: string;
      payload: Record<string, unknown>;
      roleIds: string[];
    }) => {
      const body: Record<string, unknown> = {
        email: payload.email,
        first_name: payload.first_name,
        last_name: payload.last_name,
        is_active: payload.is_active,
      };
      if (payload.password && String(payload.password).length > 0) {
        body.password = payload.password;
      }
      await api.put(`/users/${id}`, body);
      await api.put(`/users/${id}/roles`, { role_ids: roleIds });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setEditId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/users/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const openEdit = async (row: User) => {
    setFeedback(null);
    try {
      const res = await api.get<User>(`/users/${row.id}`);
      const user = res.data!;
      setEditId(user.id);
      setEditValues({
        email: user.email,
        first_name: user.first_name,
        last_name: user.last_name,
        is_active: user.is_active,
        password: "",
      });
      setEditRoleIds(user.role_ids ?? []);
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    }
  };

  const roleNameMap = Object.fromEntries(roles.map((r) => [r.id, r.name]));

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.users.title")}</h1>
          <p className="text-muted-foreground">{t("pages.users.subtitle")}</p>
        </div>
        <Button
          onClick={() => {
            setShowForm(!showForm);
            setFeedback(null);
          }}
        >
          {showForm ? t("common.cancel") : t("pages.users.addUser")}
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
            <CardTitle>{t("pages.users.create")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={createFields}
              values={formValues}
              onChange={(n, v) => setFormValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() => createMutation.mutate(formValues)}
              submitLabel={t("common.create")}
              loading={createMutation.isPending}
            />
            <div className="mt-4 space-y-2">
              <Label>{t("pages.users.assignRoles")}</Label>
              <div className="flex flex-wrap gap-3">
                {roles.map((role) => (
                  <label key={role.id} className="flex items-center gap-2 text-sm">
                    <input
                      type="checkbox"
                      checked={createRoleIds.includes(role.id)}
                      onChange={(e) => {
                        setCreateRoleIds((prev) =>
                          e.target.checked
                            ? [...prev, role.id]
                            : prev.filter((id) => id !== role.id)
                        );
                      }}
                      className="h-4 w-4"
                    />
                    {role.name}
                  </label>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={(data?.users ?? []) as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "email", header: t("common.email") },
          {
            key: "name",
            header: t("pages.users.fullName"),
            render: (row) =>
              `${row.first_name as string} ${row.last_name as string}`.trim() || "—",
          },
          {
            key: "roles",
            header: t("common.roles"),
            render: (row) => {
              const ids = (row.role_ids as string[] | undefined) ?? [];
              if (ids.length === 0) return "—";
              return (
                <div className="flex flex-wrap gap-1">
                  {ids.map((id) => (
                    <Badge key={id} variant="outline">
                      {roleNameMap[id] ?? id}
                    </Badge>
                  ))}
                </div>
              );
            },
          },
          {
            key: "is_active",
            header: t("common.status"),
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? t("common.active") : t("common.inactive")}
              </Badge>
            ),
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
              <div className="flex flex-wrap gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => openEdit(row as unknown as User)}
                >
                  {t("common.edit")}
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => setDeleteId(row.id as string)}
                >
                  {t("common.delete")}
                </Button>
              </div>
            ),
          },
        ]}
      />

      {data && (
        <PaginationBar
          page={page}
          perPage={PER_PAGE}
          total={data.total}
          onPageChange={setPage}
        />
      )}

      {editId && (
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.users.editUser")}</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={[
                { name: "email", label: t("common.email"), type: "email", required: true },
                { name: "first_name", label: t("common.firstName"), type: "text" },
                { name: "last_name", label: t("common.lastName"), type: "text" },
                {
                  name: "password",
                  label: t("pages.users.newPassword"),
                  type: "password",
                },
                { name: "is_active", label: t("common.active"), type: "checkbox" },
              ]}
              values={editValues}
              onChange={(n, v) => setEditValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() =>
                updateMutation.mutate({
                  id: editId,
                  payload: editValues,
                  roleIds: editRoleIds,
                })
              }
              submitLabel={t("common.save")}
              loading={updateMutation.isPending}
            />
            <div className="mt-4 space-y-2">
              <Label>{t("pages.users.assignRoles")}</Label>
              <div className="flex flex-wrap gap-3">
                {roles.map((role) => (
                  <label key={role.id} className="flex items-center gap-2 text-sm">
                    <input
                      type="checkbox"
                      checked={editRoleIds.includes(role.id)}
                      onChange={(e) => {
                        setEditRoleIds((prev) =>
                          e.target.checked
                            ? [...prev, role.id]
                            : prev.filter((id) => id !== role.id)
                        );
                      }}
                      className="h-4 w-4"
                    />
                    {role.name}
                  </label>
                ))}
              </div>
            </div>
            <Button variant="ghost" className="mt-4" onClick={() => setEditId(null)}>
              {t("common.cancel")}
            </Button>
          </CardContent>
        </Card>
      )}

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
