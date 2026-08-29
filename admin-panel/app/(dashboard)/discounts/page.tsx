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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Discount {
  id: string;
  code: string;
  type: string;
  value: number;
  min_order_amount: number;
  max_uses: number;
  used_count: number;
  starts_at?: string;
  ends_at?: string;
  is_active: boolean;
}

interface PreviewResult {
  code: string;
  discount_amount: number;
  final_amount: number;
}

function discountStatus(d: Discount): { label: string; variant: "success" | "secondary" | "destructive" } {
  const now = Date.now();
  if (!d.is_active) return { label: t("common.inactive"), variant: "secondary" };
  if (d.ends_at && new Date(d.ends_at).getTime() < now) {
    return { label: t("pages.discounts.expired"), variant: "destructive" };
  }
  if (d.max_uses > 0 && d.used_count >= d.max_uses) {
    return { label: t("pages.discounts.maxUsesReached"), variant: "destructive" };
  }
  return { label: t("common.active"), variant: "success" };
}

export default function DiscountsPage() {
  const queryClient = useQueryClient();
  const [code, setCode] = useState("");
  const [type, setType] = useState("percent");
  const [value, setValue] = useState(10);
  const [maxUses, setMaxUses] = useState(0);
  const [minOrder, setMinOrder] = useState(0);
  const [previewCode, setPreviewCode] = useState("");
  const [previewAmount, setPreviewAmount] = useState(10000);
  const [previewResult, setPreviewResult] = useState<PreviewResult | null>(null);
  const [editDiscount, setEditDiscount] = useState<Discount | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.discounts,
    queryFn: async () => {
      const res = await api.get<Discount[]>("/discounts");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/discounts", {
        code,
        type,
        value,
        max_uses: maxUses,
        min_order_amount: minOrder,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.discounts });
      setCode("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async (d: Discount) => {
      await api.put(`/discounts/${d.id}`, {
        type: d.type,
        value: d.value,
        max_uses: d.max_uses,
        min_order_amount: d.min_order_amount,
        is_active: d.is_active,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.discounts });
      setEditDiscount(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deactivateMutation = useMutation({
    mutationFn: async (d: Discount) => {
      await api.put(`/discounts/${d.id}`, {
        type: d.type,
        value: d.value,
        max_uses: d.max_uses,
        min_order_amount: d.min_order_amount,
        is_active: false,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.discounts });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/discounts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.discounts });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const previewMutation = useMutation({
    mutationFn: async () => {
      const res = await api.post<PreviewResult>("/discounts/preview", {
        code: previewCode,
        order_amount: previewAmount,
      });
      return res.data!;
    },
    onSuccess: (data) => {
      setPreviewResult(data);
      setFeedback(null);
    },
    onError: (err) => {
      setPreviewResult(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.discounts.title")}</h1>
        <p className="text-muted-foreground">{t("pages.discounts.subtitle")}</p>
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

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.discounts.createDiscount")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.code")}</Label>
              <Input value={code} onChange={(e) => setCode(e.target.value.toUpperCase())} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.type")}</Label>
              <Select value={type} onValueChange={setType}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="percent">{t("pages.discounts.percent")}</SelectItem>
                  <SelectItem value="fixed">{t("pages.discounts.fixed")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{type === "percent" ? t("pages.discounts.percentValue") : t("common.amountCents")}</Label>
              <Input
                type="number"
                value={value}
                onChange={(e) => setValue(Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("pages.discounts.maxUses")}</Label>
              <Input
                type="number"
                min={0}
                value={maxUses}
                onChange={(e) => setMaxUses(Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("pages.discounts.minOrder")}</Label>
              <Input
                type="number"
                min={0}
                value={minOrder}
                onChange={(e) => setMinOrder(Number(e.target.value))}
              />
            </div>
            <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending || !code.trim()}>
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
              <Input value={previewCode} onChange={(e) => setPreviewCode(e.target.value)} />
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
              <div className="rounded-md bg-muted p-3 text-sm space-y-1">
                <p>
                  {t("pages.discounts.discountAmount")}: {(previewResult.discount_amount / 100).toFixed(2)}
                </p>
                <p>
                  {t("pages.discounts.finalAmount")}: {(previewResult.final_amount / 100).toFixed(2)}
                </p>
              </div>
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
          {
            key: "type",
            header: t("common.type"),
            render: (row) =>
              row.type === "percent" ? t("pages.discounts.percent") : t("pages.discounts.fixed"),
          },
          { key: "value", header: t("common.value") },
          {
            key: "usage",
            header: t("pages.discounts.usage"),
            render: (row) => `${row.used_count} / ${row.max_uses || "∞"}`,
          },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => {
              const st = discountStatus(row as unknown as Discount);
              return <Badge variant={st.variant}>{st.label}</Badge>;
            },
          },
          {
            key: "ends_at",
            header: t("common.to"),
            render: (row) =>
              row.ends_at ? formatJalaliDateTime(row.ends_at as string) : "—",
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => {
              const d = row as unknown as Discount;
              return (
                <div className="flex gap-2">
                  <Button size="sm" variant="outline" onClick={() => setEditDiscount({ ...d })}>
                    {t("common.edit")}
                  </Button>
                  {d.is_active && (
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => deactivateMutation.mutate(d)}
                      disabled={deactivateMutation.isPending}
                    >
                      {t("common.inactive")}
                    </Button>
                  )}
                  <Button size="sm" variant="destructive" onClick={() => setDeleteId(d.id)}>
                    {t("common.delete")}
                  </Button>
                </div>
              );
            },
          },
        ]}
      />

      <Dialog open={!!editDiscount} onOpenChange={() => setEditDiscount(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.edit")} — {editDiscount?.code}</DialogTitle>
          </DialogHeader>
          {editDiscount && (
            <div className="space-y-3">
              <div className="space-y-2">
                <Label>{t("common.type")}</Label>
                <Select
                  value={editDiscount.type}
                  onValueChange={(v) => setEditDiscount({ ...editDiscount, type: v })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="percent">{t("pages.discounts.percent")}</SelectItem>
                    <SelectItem value="fixed">{t("pages.discounts.fixed")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>{t("common.value")}</Label>
                <Input
                  type="number"
                  value={editDiscount.value}
                  onChange={(e) =>
                    setEditDiscount({ ...editDiscount, value: Number(e.target.value) })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label>{t("pages.discounts.maxUses")}</Label>
                <Input
                  type="number"
                  value={editDiscount.max_uses}
                  onChange={(e) =>
                    setEditDiscount({ ...editDiscount, max_uses: Number(e.target.value) })
                  }
                />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDiscount(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              onClick={() => editDiscount && updateMutation.mutate(editDiscount)}
              disabled={updateMutation.isPending}
            >
              {t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">{t("common.confirmDeleteMessage")}</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
              disabled={deleteMutation.isPending}
            >
              {t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
