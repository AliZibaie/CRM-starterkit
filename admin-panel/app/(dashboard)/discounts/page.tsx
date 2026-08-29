"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface Discount {
  id: string;
  code: string;
  type: string;
  value: number;
  is_active: boolean;
}

export default function DiscountsPage() {
  const queryClient = useQueryClient();
  const [code, setCode] = useState("");
  const [type, setType] = useState("percent");
  const [value, setValue] = useState(10);
  const [previewCode, setPreviewCode] = useState("");
  const [previewAmount, setPreviewAmount] = useState(10000);
  const [previewResult, setPreviewResult] = useState<string | null>(null);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.discounts,
    queryFn: async () => {
      const res = await api.get<Discount[]>("/discounts");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/discounts", { code, type, value, is_active: true });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.discounts });
      setCode("");
    },
  });

  const previewMutation = useMutation({
    mutationFn: async () => {
      const res = await api.post<Record<string, unknown>>("/discounts/preview", {
        code: previewCode,
        amount: previewAmount,
      });
      return res.data;
    },
    onSuccess: (data) => setPreviewResult(JSON.stringify(data, null, 2)),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.discounts.title")}</h1>
        <p className="text-muted-foreground">{t("pages.discounts.subtitle")}</p>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.discounts.createDiscount")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.code")}</Label>
              <Input value={code} onChange={(e) => setCode(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.type")}</Label>
              <Input value={type} onChange={(e) => setType(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.value")}</Label>
              <Input
                type="number"
                value={value}
                onChange={(e) => setValue(Number(e.target.value))}
              />
            </div>
            <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
              {t("common.create")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("pages.discounts.previewDiscount")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.code")}</Label>
              <Input
                value={previewCode}
                onChange={(e) => setPreviewCode(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("common.amountCents")}</Label>
              <Input
                type="number"
                value={previewAmount}
                onChange={(e) => setPreviewAmount(Number(e.target.value))}
              />
            </div>
            <Button onClick={() => previewMutation.mutate()} disabled={previewMutation.isPending}>
              {t("common.preview")}
            </Button>
            {previewResult && (
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs">
                {previewResult}
              </pre>
            )}
          </CardContent>
        </Card>
      </div>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "code", header: t("common.code") },
          { key: "type", header: t("common.type") },
          { key: "value", header: t("common.value") },
          {
            key: "is_active",
            header: t("common.active"),
            render: (row) => (row.is_active ? t("common.yes") : t("common.no")),
          },
        ]}
      />
    </div>
  );
}
