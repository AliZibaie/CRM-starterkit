"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface Lead {
  id: string;
  name: string;
  email: string;
  phone: string;
  source: string;
  status: string;
}

const fields: FormField[] = [
  { name: "name", label: "Name", type: "text", required: true },
  { name: "email", label: "Email", type: "email" },
  { name: "phone", label: "Phone", type: "text" },
  { name: "source", label: "Source", type: "text" },
];

export default function LeadsPage() {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, unknown>>({
    name: "",
    email: "",
    phone: "",
    source: "",
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.leads,
    queryFn: async () => {
      const res = await api.get<Lead[]>("/crm/leads");
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
    },
  });

  const convertMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/crm/leads/${id}/convert`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.leads }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Leads</h1>
        <p className="text-muted-foreground">CRM lead management</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Lead</CardTitle>
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

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: "Name" },
          { key: "email", header: "Email" },
          { key: "phone", header: "Phone" },
          { key: "source", header: "Source" },
          { key: "status", header: "Status" },
          {
            key: "actions",
            header: "Actions",
            render: (row) => (
              <Button
                size="sm"
                onClick={() => convertMutation.mutate(row.id as string)}
              >
                Convert
              </Button>
            ),
          },
        ]}
      />
    </div>
  );
}
