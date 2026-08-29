"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Department {
  id: string;
  name: string;
  description: string;
}

const fields: FormField[] = [
  { name: "name", label: t("common.name"), type: "text", required: true },
  { name: "description", label: t("common.description"), type: "textarea" },
];

export default function SupportDepartmentsPage() {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, unknown>>({
    name: "",
    description: "",
  });
  const [editId, setEditId] = useState<string | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.departments,
    queryFn: async () => {
      const res = await api.get<Department[]>("/support/departments");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/support/departments", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.departments });
      setValues({ name: "", description: "" });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/support/departments/${editId}`, values);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.departments });
      setEditId(null);
      setValues({ name: "", description: "" });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/support/departments/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.departments });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.departments.title")}</h1>
        <p className="text-muted-foreground">{t("pages.departments.subtitle")}</p>
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

      <Card>
        <CardHeader>
          <CardTitle>
            {editId ? t("common.edit") : t("pages.departments.createDepartment")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <DynamicForm
            fields={fields}
            values={values}
            onChange={(n, v) => setValues((p) => ({ ...p, [n]: v }))}
            onSubmit={() => (editId ? updateMutation.mutate() : createMutation.mutate(values))}
            loading={createMutation.isPending || updateMutation.isPending}
            submitLabel={editId ? t("common.save") : t("common.create")}
          />
          {editId && (
            <Button
              variant="ghost"
              className="mt-2"
              onClick={() => {
                setEditId(null);
                setValues({ name: "", description: "" });
              }}
            >
              {t("common.cancel")}
            </Button>
          )}
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          { key: "description", header: t("common.description") },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    setEditId(row.id as string);
                    setValues({
                      name: row.name,
                      description: row.description,
                    });
                  }}
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
