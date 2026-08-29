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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Lead {
  id: string;
  name: string;
  email: string;
  phone: string;
  source: string;
  status: string;
  owner_id?: string;
}

interface User {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
}

const fields: FormField[] = [
  { name: "name", label: t("common.name"), type: "text", required: true },
  { name: "email", label: t("common.email"), type: "email" },
  { name: "phone", label: t("common.phone"), type: "text" },
  { name: "source", label: t("common.source"), type: "text" },
];

export default function LeadsPage() {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, unknown>>({
    name: "",
    email: "",
    phone: "",
    source: "",
  });
  const [assignLeadId, setAssignLeadId] = useState<string | null>(null);
  const [ownerId, setOwnerId] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.leads,
    queryFn: async () => {
      const res = await api.get<Lead[]>("/crm/leads");
      return res.data ?? [];
    },
  });

  const { data: users = [] } = useQuery({
    queryKey: queryKeys.users(1),
    queryFn: async () => {
      const res = await api.get<User[]>("/users?page=1&per_page=100");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/crm/leads", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.leads });
      setValues({ name: "", email: "", phone: "", source: "" });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const convertMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/crm/leads/${id}/convert`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.leads });
      queryClient.invalidateQueries({ queryKey: queryKeys.contacts });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const assignMutation = useMutation({
    mutationFn: async () => {
      if (!assignLeadId || !ownerId) return;
      await api.post("/crm/assign-owner", {
        entity_type: "lead",
        entity_id: assignLeadId,
        owner_id: ownerId,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.leads });
      setAssignLeadId(null);
      setOwnerId("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const userLabel = (u: User) =>
    [u.first_name, u.last_name].filter(Boolean).join(" ") || u.email;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.leads.title")}</h1>
        <p className="text-muted-foreground">{t("pages.leads.subtitle")}</p>
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
          <CardTitle>{t("pages.leads.createLead")}</CardTitle>
        </CardHeader>
        <CardContent>
          <DynamicForm
            fields={fields}
            values={values}
            onChange={(n, v) => setValues((p) => ({ ...p, [n]: v }))}
            onSubmit={() => createMutation.mutate(values)}
            loading={createMutation.isPending}
          />
        </CardContent>
      </Card>

      {assignLeadId && (
        <Card>
          <CardHeader>
            <CardTitle>{t("common.assignAgent")}</CardTitle>
          </CardHeader>
          <CardContent className="flex gap-2">
            <Select value={ownerId} onValueChange={setOwnerId}>
              <SelectTrigger className="max-w-xs">
                <SelectValue placeholder={t("common.selectUser")} />
              </SelectTrigger>
              <SelectContent>
                {users.map((u) => (
                  <SelectItem key={u.id} value={u.id}>
                    {userLabel(u)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button onClick={() => assignMutation.mutate()} disabled={assignMutation.isPending || !ownerId}>
              {t("common.assign")}
            </Button>
            <Button variant="outline" onClick={() => setAssignLeadId(null)}>
              {t("common.cancel")}
            </Button>
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
          { key: "source", header: t("common.source") },
          { key: "status", header: t("common.status") },
          { key: "owner_id", header: t("common.userId"), className: "font-mono text-xs" },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button size="sm" variant="outline" onClick={() => setAssignLeadId(row.id as string)}>
                  {t("common.assign")}
                </Button>
                <Button
                  size="sm"
                  onClick={() => convertMutation.mutate(row.id as string)}
                  disabled={row.status === "converted" || convertMutation.isPending}
                >
                  {t("common.convert")}
                </Button>
              </div>
            ),
          },
        ]}
      />
    </div>
  );
}
