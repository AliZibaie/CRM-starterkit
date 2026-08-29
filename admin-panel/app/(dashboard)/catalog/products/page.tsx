"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface Product {
  id: string;
  name: string;
  type: string;
  base_price: number;
  currency: string;
}

interface Category {
  id: string;
  name: string;
}

export default function ProductsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [values, setValues] = useState<Record<string, unknown>>({
    type: "product",
    name: "",
    description: "",
    category_id: "",
    currency: "USD",
    base_price: 0,
  });

  const { data: categories = [] } = useQuery({
    queryKey: queryKeys.categories,
    queryFn: async () => {
      const res = await api.get<Category[]>("/catalog/categories");
      return res.data ?? [];
    },
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.products,
    queryFn: async () => {
      const res = await api.get<Product[]>("/catalog/items");
      return res.data ?? [];
    },
  });

  const fields: FormField[] = [
    { name: "name", label: t("common.name"), type: "text", required: true },
    { name: "type", label: t("common.type"), type: "text", required: true },
    { name: "description", label: t("common.description"), type: "textarea" },
    {
      name: "category_id",
      label: t("common.category"),
      type: "select",
      options: categories.map((c) => ({ label: c.name, value: c.id })),
    },
    { name: "currency", label: t("common.currency"), type: "text" },
    { name: "base_price", label: t("common.basePrice"), type: "number" },
  ];

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/catalog/items", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.products });
      setShowForm(false);
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.products.title")}</h1>
          <p className="text-muted-foreground">{t("pages.products.subtitle")}</p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? t("common.cancel") : t("pages.products.addProduct")}
        </Button>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.products.createProduct")}</CardTitle>
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
      )}

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          { key: "type", header: t("common.type") },
          {
            key: "base_price",
            header: t("common.price"),
            render: (row) =>
              `${row.currency} ${((row.base_price as number) / 100).toFixed(2)}`,
          },
        ]}
      />
    </div>
  );
}
