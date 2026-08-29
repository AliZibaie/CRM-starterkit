"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { t } from "@/lib/i18n";
import { StepIndicator } from "@/components/checkout/step-indicator";
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
  const [cardNumber, setCardNumber] = useState("");
  const [cardExpiry, setCardExpiry] = useState("");
  const [cardCvv, setCardCvv] = useState("");
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
        throw new Error(t("common.error"));
      }
    },
    onError: (err: Error) => setError(err.message),
  });

  if (!accessToken) return null;

  if (cartLoading) {
    return <p className="text-muted-foreground">{t("checkout.loading")}</p>;
  }

  const items = cart?.items ?? [];
  const total = cart?.total ?? 0;
  const currency = cart?.cart?.currency ?? "IRR";

  const steps = [
    { id: "cart", label: t("checkout.stepCart") },
    { id: "checkout", label: t("checkout.stepCheckout") },
    { id: "payment", label: t("checkout.stepPayment") },
  ];

  if (items.length === 0) {
    return (
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">{t("checkout.title")}</h1>
        <p className="text-muted-foreground">{t("checkout.empty")}</p>
        <Button onClick={() => router.push("/categories/all")}>
          {t("common.browseProducts")}
        </Button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="text-3xl font-bold">{t("checkout.title")}</h1>
      <StepIndicator steps={steps} current={1} />
      <Card>
        <CardHeader>
          <CardTitle>{t("checkout.summary")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {items.map((item, i) => (
            <div key={i} className="flex justify-between text-sm">
              <span>
                {item.title} × {item.quantity}
              </span>
              <span>{formatPrice(item.unit_price * item.quantity, currency)}</span>
            </div>
          ))}
          <div className="border-t border-border pt-3 font-semibold">
            {t("cart.total")}: {formatPrice(total, currency)}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="space-y-4 pt-6">
          <div className="space-y-2">
            <Label htmlFor="discount">{t("checkout.discountCode")}</Label>
            <Input
              id="discount"
              value={discountCode}
              onChange={(e) => setDiscountCode(e.target.value)}
              placeholder={t("common.optional")}
            />
          </div>
          {gateways.length > 0 && (
            <div className="space-y-2">
              <Label htmlFor="gateway">{t("checkout.gateway")}</Label>
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
          <div className="space-y-2">
            <Label htmlFor="card">{t("checkout.cardNumber")}</Label>
            <Input
              id="card"
              value={cardNumber}
              onChange={(e) => setCardNumber(e.target.value)}
              className="ltr-field"
              placeholder="0000 0000 0000 0000"
              inputMode="numeric"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="expiry">{t("checkout.cardExpiry")}</Label>
              <Input
                id="expiry"
                value={cardExpiry}
                onChange={(e) => setCardExpiry(e.target.value)}
                className="ltr-field"
                placeholder="MM/YY"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cvv">{t("checkout.cardCvv")}</Label>
              <Input
                id="cvv"
                value={cardCvv}
                onChange={(e) => setCardCvv(e.target.value)}
                className="ltr-field"
                placeholder="123"
                inputMode="numeric"
              />
            </div>
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-between">
            <Button variant="outline" onClick={() => router.push("/cart")}>
              {t("common.back")}
            </Button>
            <Button
              onClick={() => checkout.mutate()}
              disabled={checkout.isPending || !gateway}
            >
              {checkout.isPending ? t("checkout.processing") : t("checkout.payNow")}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
