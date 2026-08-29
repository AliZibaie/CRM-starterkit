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
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.settings,
    queryFn: async () => {
      const res = await api.get<Setting[]>("/settings");
      return res.data ?? [];
    },
  });

  const saveMutation = useMutation({
    mutationFn: async () => {
      await api.put("/settings", { key: key.trim(), value });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.settings });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const selectSetting = (row: Setting) => {
    setKey(row.key);
    setValue(row.value);
    setFeedback(null);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.settings.title")}</h1>
        <p className="text-muted-foreground">{t("pages.settings.subtitle")}</p>
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
          <CardTitle>{t("pages.settings.updateSetting")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <p className="text-sm text-muted-foreground">{t("pages.settings.editHint")}</p>
          <div className="space-y-2">
            <Label>{t("common.key")}</Label>
            <Input
              value={key}
              onChange={(e) => setKey(e.target.value)}
              className="ltr-field"
              placeholder="site.name"
            />
          </div>
          <div className="space-y-2">
            <Label>{t("common.value")}</Label>
            <Input value={value} onChange={(e) => setValue(e.target.value)} />
          </div>
          <Button
            onClick={() => saveMutation.mutate()}
            disabled={saveMutation.isPending || !key.trim()}
          >
            {saveMutation.isPending ? t("common.saving") : t("common.save")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="key"
        columns={[
          {
            key: "key",
            header: t("common.key"),
            render: (row) => (
              <button
                type="button"
                className="ltr-field text-start font-mono text-xs text-primary hover:underline"
                onClick={() => selectSetting(row as unknown as Setting)}
              >
                {String(row.key)}
              </button>
            ),
          },
          {
            key: "value",
            header: t("common.value"),
            render: (row) => (
              <button
                type="button"
                className="text-start hover:underline"
                onClick={() => selectSetting(row as unknown as Setting)}
              >
                {String(row.value ?? "")}
              </button>
            ),
          },
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
