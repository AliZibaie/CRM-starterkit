"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface Setting {
  key: string;
  value: string;
  updated_at: string;
}

export default function SettingsPage() {
  const queryClient = useQueryClient();
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.settings,
    queryFn: async () => {
      const res = await api.get<Setting[]>("/settings");
      return res.data ?? [];
    },
  });

  const saveMutation = useMutation({
    mutationFn: async () => {
      await api.put("/settings", { key, value });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.settings });
      setKey("");
      setValue("");
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.settings.title")}</h1>
        <p className="text-muted-foreground">{t("pages.settings.subtitle")}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.settings.updateSetting")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>{t("common.key")}</Label>
            <Input value={key} onChange={(e) => setKey(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.value")}</Label>
            <Input value={value} onChange={(e) => setValue(e.target.value)} />
          </div>
          <Button onClick={() => saveMutation.mutate()} disabled={saveMutation.isPending}>
            {t("common.save")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="key"
        columns={[
          { key: "key", header: t("common.key") },
          { key: "value", header: t("common.value") },
          {
            key: "updated_at",
            header: t("common.updatedAt"),
            render: (row) => formatJalaliDateTime(row.updated_at as string),
          },
        ]}
      />
    </div>
  );
}
