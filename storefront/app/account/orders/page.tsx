"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { formatPrice } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

interface Order {
  id: string;
  status: string;
  total_amount: number;
  currency: string;
  created_at: string;
}

export default function OrdersPage() {
  const { data: orders = [], isLoading, error } = useQuery({
    queryKey: queryKeys.orders.all,
    queryFn: async () => {
      const res = await api.get<Order[]>("/orders");
      return res.data ?? [];
    },
  });

  if (isLoading) {
    return <p className="text-muted-foreground">Loading orders...</p>;
  }

  if (error) {
    return <p className="text-destructive">{error.message}</p>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">My Orders</h1>
          <p className="text-muted-foreground">View your order history</p>
        </div>
        <div className="flex gap-2">
          <Button asChild variant="outline">
            <Link href="/account/tickets">Support</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/account/licenses">Licenses</Link>
          </Button>
        </div>
      </div>
      {orders.length === 0 ? (
        <p className="text-muted-foreground">No orders yet.</p>
      ) : (
        <div className="space-y-3">
          {orders.map((order) => (
            <Card key={order.id}>
              <CardContent className="flex flex-wrap items-center justify-between gap-4 py-4">
                <div>
                  <p className="font-mono text-sm">{order.id}</p>
                  <p className="text-sm text-muted-foreground">
                    {new Date(order.created_at).toLocaleString()}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <span className="rounded-full border px-3 py-1 text-xs capitalize">
                    {order.status}
                  </span>
                  <span className="font-semibold">
                    {formatPrice(order.total_amount, order.currency)}
                  </span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
