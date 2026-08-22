"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface CartData {
  cart: { currency: string };
  items: { title: string; unit_price: number; quantity: number }[];
  total: number;
}

interface OrderResponse {
  order: { id: string; total_amount: number; currency: string };
}

interface PaymentInitResponse {
  payment_url: string;
}

export default function CheckoutPage() {
  const router = useRouter();
  const accessToken = useAuthStore((s) => s.accessToken);
  const [discountCode, setDiscountCode] = useState("");
  const [gateway, setGateway] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    if (!accessToken) {
      router.push("/account/login?redirect=/checkout");
    }
  }, [accessToken, router]);

  const { data: cart, isLoading: cartLoading } = useQuery({
    queryKey: queryKeys.cart,
    enabled: !!accessToken,
    queryFn: async () => {
      const res = await api.get<CartData>("/cart");
      return res.data!;
    },
  });

  const { data: gateways = [] } = useQuery({
    queryKey: queryKeys.payments.gateways,
    enabled: !!accessToken,
    queryFn: async () => {
      const res = await api.get<string[]>("/payments/gateways");
      return res.data ?? [];
    },
  });

  useEffect(() => {
    if (gateways.length > 0 && !gateway) {
      setGateway(gateways[0]);
    }
  }, [gateways, gateway]);

  const checkout = useMutation({
    mutationFn: async () => {
      setError("");
      const orderRes = await api.post<OrderResponse>("/orders", {
        discount_code: discountCode || undefined,
      });
      const orderId = orderRes.data!.order.id;
      await api.post(`/orders/${orderId}/checkout`);
      const payRes = await api.post<PaymentInitResponse>("/payments/init", {
        order_id: orderId,
        gateway,
      });
      const paymentUrl = payRes.data!.payment_url;
      if (paymentUrl) {
        window.location.href = paymentUrl;
      } else {
        throw new Error("No payment URL returned");
      }
    },
    onError: (err: Error) => setError(err.message),
  });

  if (!accessToken) {
    return null;
  }

  if (cartLoading) {
    return <p className="text-muted-foreground">Loading checkout...</p>;
  }

  const items = cart?.items ?? [];
  const total = cart?.total ?? 0;
  const currency = cart?.cart?.currency ?? "IRR";

  if (items.length === 0) {
    return (
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">Checkout</h1>
        <p className="text-muted-foreground">Your cart is empty.</p>
        <Button onClick={() => router.push("/")}>Browse Products</Button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <h1 className="text-3xl font-bold">Checkout</h1>
      <Card>
        <CardHeader>
          <CardTitle>Order Summary</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {items.map((item, i) => (
            <div key={i} className="flex justify-between text-sm">
              <span>
                {item.title} × {item.quantity}
              </span>
              <span>
                {formatPrice(item.unit_price * item.quantity, currency)}
              </span>
            </div>
          ))}
          <div className="border-t pt-3 font-semibold">
            Total: {formatPrice(total, currency)}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="space-y-4 pt-6">
          <div className="space-y-2">
            <Label htmlFor="discount">Discount Code</Label>
            <Input
              id="discount"
              value={discountCode}
              onChange={(e) => setDiscountCode(e.target.value)}
              placeholder="Optional"
            />
          </div>
          {gateways.length > 0 && (
            <div className="space-y-2">
              <Label htmlFor="gateway">Payment Gateway</Label>
              <select
                id="gateway"
                value={gateway}
                onChange={(e) => setGateway(e.target.value)}
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
              >
                {gateways.map((g) => (
                  <option key={g} value={g}>
                    {g}
                  </option>
                ))}
              </select>
            </div>
          )}
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button
            className="w-full"
            onClick={() => checkout.mutate()}
            disabled={checkout.isPending || !gateway}
          >
            {checkout.isPending ? "Processing..." : "Pay Now"}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
