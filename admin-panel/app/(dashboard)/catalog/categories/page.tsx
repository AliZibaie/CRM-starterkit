"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface Category {
  id: string;
  name: string;
  parent_id: string;
}

const fields: FormField[] = [
  { name: "name", label: t("common.name"), type: "text", required: true },
  { name: "parent_id", label: t("common.parentId"), type: "text" },
];

export default function CategoriesPage() {
  const queryClient = useQueryClient();
  const [values, setValues] = useState<Record<string, unknown>>({
    name: "",
    parent_id: "",
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.categories,
    queryFn: async () => {
      const res = await api.get<Category[]>("/catalog/categories");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/catalog/categories", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories });
      setValues({ name: "", parent_id: "" });
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.categories.title")}</h1>
        <p className="text-muted-foreground">{t("pages.categories.subtitle")}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.categories.createCategory")}</CardTitle>
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
          { key: "name", header: t("common.name") },
          { key: "parent_id", header: t("common.parentId") },
        ]}
      />
    </div>
  );
}
