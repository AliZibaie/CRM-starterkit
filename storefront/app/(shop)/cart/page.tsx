"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { t } from "@/lib/i18n";
import { StepIndicator } from "@/components/checkout/step-indicator";
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

  if (!accessToken) return null;

  if (isLoading) {
    return <p className="text-muted-foreground">{t("cart.loading")}</p>;
  }

  if (error) {
    return <p className="text-destructive">{error.message}</p>;
  }

  const items = data?.items ?? [];
  const total = data?.total ?? 0;
  const currency = data?.cart?.currency ?? "IRR";

  const steps = [
    { id: "cart", label: t("checkout.stepCart") },
    { id: "checkout", label: t("checkout.stepCheckout") },
    { id: "payment", label: t("checkout.stepPayment") },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">{t("cart.title")}</h1>
      <StepIndicator steps={steps} current={0} />
      {items.length === 0 ? (
        <div className="space-y-4">
          <p className="text-muted-foreground">{t("cart.empty")}</p>
          <Button asChild variant="outline">
            <Link href="/categories/all">{t("common.continueShopping")}</Link>
          </Button>
        </div>
      ) : (
        <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
          <div className="space-y-3">
            {items.map((item) => (
              <Card key={item.id}>
                <CardContent className="flex flex-wrap items-center justify-between gap-4 py-4">
                  <div className="text-start">
                    <p className="font-medium">{item.title}</p>
                    <p className="text-sm text-muted-foreground">
                      {t("cart.quantity")}: {item.quantity} ×{" "}
                      {formatPrice(item.unit_price, currency)}
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
                      {t("common.remove")}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
          <Card className="h-fit">
            <CardHeader>
              <CardTitle>
                {t("cart.total")}: {formatPrice(total, currency)}
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-3">
              <Button asChild size="lg">
                <Link href="/checkout">{t("cart.checkout")}</Link>
              </Button>
              <Button asChild variant="outline">
                <Link href="/categories/all">{t("common.continueShopping")}</Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
