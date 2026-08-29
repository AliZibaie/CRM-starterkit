"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { publicApi, api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { t } from "@/lib/i18n";
import { formatJalali } from "@/lib/date";
import { Breadcrumb } from "@/components/layout/breadcrumb";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

interface Product {
  id: string;
  name: string;
  slug: string;
  description?: string;
  base_price: number;
  currency: string;
  type?: string;
}

interface Variant {
  id: string;
  name: string;
  price: number;
  sku: string;
}

interface Attribute {
  key: string;
  value: string;
}

interface Review {
  id: string;
  rating: number;
  comment: string;
  created_at: string;
}

export default function ProductDetailPage() {
  const params = useParams<{ slug: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);
  const [activeTab, setActiveTab] = useState<"desc" | "reviews">("desc");

  const { data, isLoading, error } = useQuery({
    queryKey: queryKeys.products.detail(params.slug),
    queryFn: async () => {
      const res = await publicApi.get<{
        item: Product;
        attributes: Attribute[];
        variants: Variant[];
      }>(`/catalog/public/items/${params.slug}`);
      return res.data!;
    },
  });

  const { data: reviews = [] } = useQuery({
    queryKey: ["reviews", data?.item.id],
    enabled: !!data?.item.id,
    queryFn: async () => {
      const res = await publicApi.get<Review[]>(
        `/reviews/items/${data!.item.id}`
      );
      return res.data ?? [];
    },
  });

  const addToCart = useMutation({
    mutationFn: async () => {
      const variant = data?.variants?.[0];
      await api.post("/cart/items", {
        item_id: data!.item.id,
        variant_id: variant?.id ?? "",
        quantity: 1,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cart });
      router.push("/cart");
    },
  });

  if (isLoading) {
    return <p className="text-muted-foreground">{t("product.loading")}</p>;
  }

  if (error || !data) {
    return <p className="text-destructive">{t("product.notFound")}</p>;
  }

  const { item, attributes, variants } = data;
  const price = variants[0]?.price ?? item.base_price;

  function handleAddToCart() {
    if (!accessToken) {
      router.push(`/account/login?redirect=/products/${params.slug}`);
      return;
    }
    addToCart.mutate();
  }

  return (
    <div className="space-y-6">
      <Breadcrumb
        items={[
          { label: t("brand"), href: "/" },
          { label: t("nav.categories"), href: "/categories/all" },
          { label: item.name },
        ]}
      />
      <div className="grid gap-8 lg:grid-cols-2">
        <div className="order-1 flex aspect-square items-center justify-center rounded-xl border border-border bg-muted/30">
          <span className="text-6xl text-muted-foreground/40">📦</span>
        </div>
        <div className="order-2 space-y-4">
          <h1 className="text-3xl font-bold">{item.name}</h1>
          <p className="text-2xl font-bold text-primary">
            {formatPrice(price, item.currency)}
          </p>
          {variants.length > 0 && (
            <div className="space-y-1">
              <p className="text-sm font-medium">{t("product.variants")}</p>
              <ul className="space-y-1 text-sm text-muted-foreground">
                {variants.map((v) => (
                  <li key={v.id}>
                    {v.name} — {formatPrice(v.price, item.currency)}
                  </li>
                ))}
              </ul>
            </div>
          )}
          {attributes.length > 0 && (
            <div className="space-y-1">
              <p className="text-sm font-medium">{t("product.details")}</p>
              <dl className="space-y-1 text-sm">
                {attributes.map((a) => (
                  <div key={a.key} className="flex gap-2">
                    <dt className="font-medium">{a.key}:</dt>
                    <dd className="text-muted-foreground">{a.value}</dd>
                  </div>
                ))}
              </dl>
            </div>
          )}
          <Button
            onClick={handleAddToCart}
            disabled={addToCart.isPending}
            className="w-full sm:w-auto"
            size="lg"
          >
            {addToCart.isPending ? t("product.adding") : t("product.addToCart")}
          </Button>
          {addToCart.error && (
            <p className="text-sm text-destructive">{addToCart.error.message}</p>
          )}
        </div>
      </div>
      <div className="space-y-4">
        <div className="flex gap-4 border-b border-border">
          <button
            type="button"
            onClick={() => setActiveTab("desc")}
            className={`border-b-2 px-4 py-2 text-sm font-medium transition-colors ${
              activeTab === "desc"
                ? "border-primary text-primary"
                : "border-transparent text-muted-foreground"
            }`}
          >
            {t("product.description")}
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("reviews")}
            className={`border-b-2 px-4 py-2 text-sm font-medium transition-colors ${
              activeTab === "reviews"
                ? "border-primary text-primary"
                : "border-transparent text-muted-foreground"
            }`}
          >
            {t("product.reviews")} ({reviews.length})
          </button>
        </div>
        {activeTab === "desc" && (
          <Card>
            <CardContent className="pt-6 text-start leading-7">
              {item.description || (
                <p className="text-muted-foreground">{t("cms.empty")}</p>
              )}
            </CardContent>
          </Card>
        )}
        {activeTab === "reviews" && (
          <div className="space-y-3">
            {reviews.length === 0 ? (
              <p className="text-muted-foreground">{t("product.noReviews")}</p>
            ) : (
              reviews.map((review) => (
                <Card key={review.id}>
                  <CardContent className="pt-6">
                    <div className="mb-2 flex items-center gap-2">
                      <span className="text-amber-500">
                        {"★".repeat(review.rating)}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {formatJalali(review.created_at)}
                      </span>
                    </div>
                    <p className="text-sm leading-7">{review.comment}</p>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
}
