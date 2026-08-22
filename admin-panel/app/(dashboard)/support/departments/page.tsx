"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface Department {
  id: string;
  name: string;
  description: string;
}

const fields: FormField[] = [
  { name: "name", label: "Name", type: "text", required: true },
  { name: "description", label: "Description", type: "textarea" },
];

export default function SupportDepartmentsPage() {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, unknown>>({
    name: "",
    description: "",
  });

  const { data = [], isLoading } = useQuery({
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
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Support Departments</h1>
        <p className="text-muted-foreground">Manage support departments</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Department</CardTitle>
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
          { key: "description", header: "Description" },
        ]}
      />
    </div>
  );
}
