"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface CartItem {
  id: string;
  item_id: string;
  variant_id: string;
  title: string;
  unit_price: number;
  quantity: number;
}

interface CartData {
  cart: { id: string; currency: string };
  items: CartItem[];
  total: number;
}

export default function CartPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);

  useEffect(() => {
    if (!accessToken) {
      router.push("/account/login?redirect=/cart");
    }
  }, [accessToken, router]);

  const { data, isLoading, error } = useQuery({
    queryKey: queryKeys.cart,
    enabled: !!accessToken,
    queryFn: async () => {
      const res = await api.get<CartData>("/cart");
      return res.data!;
    },
  });

  const removeItem = useMutation({
    mutationFn: async (itemId: string) => {
      await api.delete(`/cart/items/${itemId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cart });
    },
  });

  if (!accessToken) {
    return null;
  }

  if (isLoading) {
    return <p className="text-muted-foreground">Loading cart...</p>;
  }

  if (error) {
    return <p className="text-destructive">{error.message}</p>;
  }

  const items = data?.items ?? [];
  const total = data?.total ?? 0;
  const currency = data?.cart?.currency ?? "IRR";

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Cart</h1>
      {items.length === 0 ? (
        <div className="space-y-4">
          <p className="text-muted-foreground">Your cart is empty.</p>
          <Button asChild variant="outline">
            <Link href="/">Continue Shopping</Link>
          </Button>
        </div>
      ) : (
        <>
          <div className="space-y-3">
            {items.map((item) => (
              <Card key={item.id}>
                <CardContent className="flex items-center justify-between py-4">
                  <div>
                    <p className="font-medium">{item.title}</p>
                    <p className="text-sm text-muted-foreground">
                      Qty: {item.quantity} × {formatPrice(item.unit_price, currency)}
                    </p>
                  </div>
                  <div className="flex items-center gap-4">
                    <p className="font-semibold">
                      {formatPrice(item.unit_price * item.quantity, currency)}
                    </p>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => removeItem.mutate(item.item_id)}
                      disabled={removeItem.isPending}
                    >
                      Remove
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
          <Card>
            <CardHeader>
              <CardTitle>Total: {formatPrice(total, currency)}</CardTitle>
            </CardHeader>
            <CardContent>
              <Button asChild>
                <Link href="/checkout">Proceed to Checkout</Link>
              </Button>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
