"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { publicApi } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { formatPrice } from "@/lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Product {
  id: string;
  name: string;
  slug: string;
  description?: string;
  base_price: number;
  currency: string;
}

interface Category {
  id: string;
  name: string;
  slug: string;
}

export default function CategoryPage() {
  const params = useParams<{ slug: string }>();
  const isAll = params.slug === "all";

  const { data, isLoading } = useQuery({
    queryKey: isAll
      ? queryKeys.products.all
      : queryKeys.categories.detail(params.slug),
    queryFn: async () => {
      if (isAll) {
        const res = await publicApi.get<Product[]>("/catalog/public/items");
        return { category: null as Category | null, items: res.data ?? [] };
      }
      const res = await publicApi.get<{
        category: Category;
        items: Product[];
      }>(`/catalog/public/categories/${params.slug}`);
      return {
        category: res.data!.category,
        items: res.data!.items ?? [],
      };
    },
  });

  if (isLoading) {
    return <p className="text-muted-foreground">Loading...</p>;
  }

  const category = data?.category;
  const items = data?.items ?? [];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">
          {isAll ? "All Products" : category?.name ?? params.slug}
        </h1>
        {!isAll && category && (
          <p className="text-muted-foreground">Category: {category.slug}</p>
        )}
      </div>
      {items.length === 0 ? (
        <p className="text-muted-foreground">No products in this category.</p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((product) => (
            <Link key={product.id} href={`/products/${product.slug}`}>
              <Card className="h-full transition-shadow hover:shadow-md">
                <CardHeader>
                  <CardTitle>{product.name}</CardTitle>
                </CardHeader>
                <CardContent>
                  {product.description && (
                    <p className="mb-2 line-clamp-2 text-sm text-muted-foreground">
                      {product.description}
                    </p>
                  )}
                  <p className="font-semibold text-primary">
                    {formatPrice(product.base_price, product.currency)}
                  </p>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
