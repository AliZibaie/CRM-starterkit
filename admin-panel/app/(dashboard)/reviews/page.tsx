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
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Review {
  id: string;
  item_id: string;
  user_id: string;
  rating: number;
  comment: string;
  is_visible: boolean;
  created_at: string;
}

export default function ReviewsPage() {
  const queryClient = useQueryClient();
  const [itemFilter, setItemFilter] = useState("");
  const [ratingFilter, setRatingFilter] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const filters = { item_id: itemFilter, rating: ratingFilter };

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.reviews(filters),
    queryFn: async () => {
      const params = new URLSearchParams();
      if (itemFilter) params.set("item_id", itemFilter);
      if (ratingFilter) params.set("rating", ratingFilter);
      const res = await api.get<Review[]>(`/reviews/admin?${params.toString()}`);
      return res.data ?? [];
    },
  });

  const visibilityMutation = useMutation({
    mutationFn: async ({ id, visible }: { id: string; visible: boolean }) => {
      await api.put(`/reviews/${id}/visibility`, { is_visible: visible });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reviews"] });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/reviews/${id}/admin`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reviews"] });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.reviews.title")}</h1>
        <p className="text-muted-foreground">{t("pages.reviews.subtitle")}</p>
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
          <CardTitle>{t("common.filters")}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>{t("common.product")} ID</Label>
            <Input value={itemFilter} onChange={(e) => setItemFilter(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("pages.reviews.rating")}</Label>
            <Input type="number" min={1} max={5} value={ratingFilter} onChange={(e) => setRatingFilter(e.target.value)} />
          </div>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "item_id", header: t("common.product"), className: "font-mono text-xs" },
          { key: "rating", header: t("pages.reviews.rating") },
          { key: "comment", header: t("common.message") },
          {
            key: "is_visible",
            header: t("common.status"),
            render: (row) => (
              <Badge variant={row.is_visible ? "success" : "secondary"}>
                {row.is_visible ? t("pages.reviews.approved") : t("pages.reviews.pending")}
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
            render: (row) => {
              const visible = Boolean(row.is_visible);
              return (
              <div className="flex gap-2">
                {!visible && (
                  <Button
                    size="sm"
                    onClick={() => visibilityMutation.mutate({ id: row.id as string, visible: true })}
                  >
                    {t("pages.reviews.approve")}
                  </Button>
                )}
                {visible && (
                  <Button
                    size="sm"
                    variant="secondary"
                    onClick={() => visibilityMutation.mutate({ id: row.id as string, visible: false })}
                  >
                    {t("pages.reviews.reject")}
                  </Button>
                )}
                <Button size="sm" variant="destructive" onClick={() => deleteMutation.mutate(row.id as string)}>
                  {t("common.delete")}
                </Button>
              </div>
            );
            },
          },
        ]}
      />
    </div>
  );
}
