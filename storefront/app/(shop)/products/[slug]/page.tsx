"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { publicApi, api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { formatPrice } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

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
    return <p className="text-muted-foreground">Loading product...</p>;
  }

  if (error || !data) {
    return <p className="text-destructive">Product not found.</p>;
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
    <div className="space-y-8">
      <div className="grid gap-8 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-2xl">{item.name}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {item.description && (
              <p className="text-muted-foreground">{item.description}</p>
            )}
            <p className="text-2xl font-bold text-primary">
              {formatPrice(price, item.currency)}
            </p>
            {variants.length > 0 && (
              <div className="space-y-1">
                <p className="text-sm font-medium">Variants</p>
                <ul className="text-sm text-muted-foreground">
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
                <p className="text-sm font-medium">Details</p>
                <dl className="text-sm">
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
              className="w-full"
            >
              {addToCart.isPending ? "Adding..." : "Add to Cart"}
            </Button>
            {addToCart.error && (
              <p className="text-sm text-destructive">{addToCart.error.message}</p>
            )}
          </CardContent>
        </Card>
      </div>

      {reviews.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-xl font-semibold">Reviews</h2>
          <div className="space-y-3">
            {reviews.map((review) => (
              <Card key={review.id}>
                <CardContent className="pt-6">
                  <div className="mb-1 flex items-center gap-2">
                    <span className="font-medium">{"★".repeat(review.rating)}</span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(review.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  <p className="text-sm">{review.comment}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
