"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Order {
  id: string;
  status: string;
  created_at: string;
}

interface OrderItem {
  id: string;
  item_id: string;
  title: string;
  quantity: number;
}

interface OrderDetail {
  order: Order;
  items: OrderItem[];
}

interface License {
  orderId: string;
  orderDate: string;
  item: OrderItem;
}

export default function LicensesPage() {
  const { data: orders = [], isLoading: ordersLoading } = useQuery({
    queryKey: queryKeys.orders.all,
    queryFn: async () => {
      const res = await api.get<Order[]>("/orders");
      return res.data ?? [];
    },
  });

  const paidOrders = orders.filter((o) => o.status === "paid");

  const { data: licenses = [], isLoading: licensesLoading } = useQuery({
    queryKey: ["licenses", paidOrders.map((o) => o.id)],
    enabled: paidOrders.length > 0,
    queryFn: async () => {
      const results: License[] = [];
      for (const order of paidOrders) {
        const res = await api.get<OrderDetail>(`/orders/${order.id}`);
        const detail = res.data!;
        for (const item of detail.items ?? []) {
          results.push({
            orderId: order.id,
            orderDate: order.created_at,
            item,
          });
        }
      }
      return results;
    },
  });

  const isLoading = ordersLoading || licensesLoading;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Licenses</h1>
          <p className="text-muted-foreground">
            Products from your paid orders
          </p>
        </div>
        <Button asChild variant="outline">
          <Link href="/account/orders">Back to Orders</Link>
        </Button>
      </div>
      {isLoading ? (
        <p className="text-muted-foreground">Loading licenses...</p>
      ) : licenses.length === 0 ? (
        <p className="text-muted-foreground">
          No licenses yet. Complete a purchase to see your licenses here.
        </p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2">
          {licenses.map((license) => (
            <Card key={`${license.orderId}-${license.item.id}`}>
              <CardHeader>
                <CardTitle className="text-base">{license.item.title}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <p>
                  <span className="text-muted-foreground">License Key: </span>
                  <code className="rounded bg-muted px-2 py-1 font-mono text-xs">
                    {license.item.item_id}-{license.item.id}
                  </code>
                </p>
                <p>
                  <span className="text-muted-foreground">Quantity: </span>
                  {license.item.quantity}
                </p>
                <p>
                  <span className="text-muted-foreground">Order: </span>
                  <span className="font-mono text-xs">{license.orderId}</span>
                </p>
                <p className="text-muted-foreground">
                  Purchased {new Date(license.orderDate).toLocaleDateString()}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
