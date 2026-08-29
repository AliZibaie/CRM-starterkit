"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { datetimeLocalToRFC3339 } from "@/lib/datetime-input";
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
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Order {
  id: string;
  user_id: string;
  status: string;
  subtotal: number;
  discount_code?: string;
  discount_amount: number;
  total_amount: number;
  currency: string;
  created_at: string;
}

interface OrderItem {
  id: string;
  title: string;
  unit_price: number;
  quantity: number;
  line_total: number;
}

interface Invoice {
  id: string;
  number: string;
  amount: number;
  status: string;
}

interface OrderDetail {
  order: Order;
  items: OrderItem[];
  invoice?: Invoice;
}

const ORDER_STATUSES = [
  "pending",
  "paid",
  "processing",
  "shipped",
  "delivered",
  "refunded",
  "cancelled",
];

function buildOrdersQuery(filters: {
  status: string;
  userId: string;
  from: string;
  to: string;
}) {
  const params = new URLSearchParams({ scope: "all" });
  if (filters.status) params.set("status", filters.status);
  if (filters.userId) params.set("user_id", filters.userId);
  const from = datetimeLocalToRFC3339(filters.from);
  const to = datetimeLocalToRFC3339(filters.to);
  if (from) params.set("from", from);
  if (to) params.set("to", to);
  return `/orders?${params.toString()}`;
}

export default function OrdersPage() {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState("");
  const [userFilter, setUserFilter] = useState("");
  const [fromFilter, setFromFilter] = useState("");
  const [toFilter, setToFilter] = useState("");
  const [detailId, setDetailId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const filters = { status: statusFilter, userId: userFilter, from: fromFilter, to: toFilter };

  const { data = [], isLoading, refetch } = useQuery({
    queryKey: queryKeys.ordersFiltered(filters),
    queryFn: async () => {
      const res = await api.get<Order[]>(buildOrdersQuery(filters));
      return res.data ?? [];
    },
  });

  const { data: detail } = useQuery({
    queryKey: queryKeys.order(detailId ?? ""),
    queryFn: async () => {
      const res = await api.get<OrderDetail>(`/orders/${detailId}`);
      return res.data!;
    },
    enabled: !!detailId,
  });

  const statusMutation = useMutation({
    mutationFn: async ({ id, status }: { id: string; status: string }) => {
      await api.put(`/orders/${id}/status`, { status });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const refundMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/orders/${id}/refund`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const formatMoney = (amount: number, currency: string) =>
    `${currency} ${(amount / 100).toFixed(2)}`;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.orders.title")}</h1>
        <p className="text-muted-foreground">{t("pages.orders.subtitle")}</p>
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
        <CardContent className="grid gap-4 md:grid-cols-4">
          <div className="space-y-2">
            <Label>{t("common.status")}</Label>
            <Select value={statusFilter || "__all"} onValueChange={(v) => setStatusFilter(v === "__all" ? "" : v)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="__all">{t("common.allStatuses")}</SelectItem>
                {ORDER_STATUSES.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>{t("common.userId")}</Label>
            <Input value={userFilter} onChange={(e) => setUserFilter(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.from")}</Label>
            <Input type="datetime-local" value={fromFilter} onChange={(e) => setFromFilter(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.to")}</Label>
            <Input type="datetime-local" value={toFilter} onChange={(e) => setToFilter(e.target.value)} />
          </div>
          <Button className="md:col-span-4 w-fit" variant="outline" onClick={() => refetch()}>
            {t("common.refresh")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "id", header: t("common.id"), className: "font-mono text-xs" },
          { key: "user_id", header: t("common.userId"), className: "font-mono text-xs" },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => (
              <Select
                value={String(row.status)}
                onValueChange={(status) =>
                  statusMutation.mutate({ id: row.id as string, status })
                }
              >
                <SelectTrigger className="h-8 w-36">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {ORDER_STATUSES.map((s) => (
                    <SelectItem key={s} value={s}>
                      {s}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            ),
          },
          {
            key: "total_amount",
            header: t("common.total"),
            render: (row) => formatMoney(row.total_amount as number, String(row.currency)),
          },
          {
            key: "created_at",
            header: t("common.createdAt"),
            render: (row) => formatJalaliDateTime(row.created_at as string),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button size="sm" variant="outline" onClick={() => setDetailId(row.id as string)}>
                  {t("common.view")}
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => refundMutation.mutate(row.id as string)}
                  disabled={refundMutation.isPending}
                >
                  {t("common.refund")}
                </Button>
              </div>
            ),
          },
        ]}
      />

      <Dialog open={!!detailId} onOpenChange={() => setDetailId(null)}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t("pages.orders.orderDetails")}</DialogTitle>
          </DialogHeader>
          {detail && (
            <div className="space-y-4 text-sm">
              <div className="grid grid-cols-2 gap-2">
                <span className="text-muted-foreground">{t("common.id")}</span>
                <span className="font-mono">{detail.order.id}</span>
                <span className="text-muted-foreground">{t("common.status")}</span>
                <Badge variant="outline">{detail.order.status}</Badge>
                <span className="text-muted-foreground">{t("common.userId")}</span>
                <span className="font-mono">{detail.order.user_id}</span>
                <span className="text-muted-foreground">{t("pages.orders.subtotal")}</span>
                <span>{formatMoney(detail.order.subtotal, detail.order.currency)}</span>
                {detail.order.discount_code && (
                  <>
                    <span className="text-muted-foreground">{t("pages.orders.discount")}</span>
                    <span>
                      {detail.order.discount_code} (−
                      {formatMoney(detail.order.discount_amount, detail.order.currency)})
                    </span>
                  </>
                )}
                <span className="text-muted-foreground">{t("common.total")}</span>
                <span className="font-semibold">
                  {formatMoney(detail.order.total_amount, detail.order.currency)}
                </span>
              </div>

              <div>
                <h4 className="mb-2 font-medium">{t("pages.orders.items")}</h4>
                <div className="rounded border divide-y">
                  {detail.items.map((item) => (
                    <div key={item.id} className="flex justify-between px-3 py-2">
                      <span>
                        {item.title} × {item.quantity}
                      </span>
                      <span>{formatMoney(item.line_total, detail.order.currency)}</span>
                    </div>
                  ))}
                </div>
              </div>

              {detail.invoice && (
                <div className="grid grid-cols-2 gap-2">
                  <span className="text-muted-foreground">{t("pages.orders.invoice")}</span>
                  <span>
                    {detail.invoice.number} — {detail.invoice.status} (
                    {formatMoney(detail.invoice.amount, detail.order.currency)})
                  </span>
                </div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
