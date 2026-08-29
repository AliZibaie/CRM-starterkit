"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { publicApi } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { formatPrice } from "@/lib/utils";
import { t } from "@/lib/i18n";
import { Breadcrumb } from "@/components/layout/breadcrumb";
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

type SortOption = "newest" | "price-asc" | "price-desc";

export default function CategoryPage() {
  const params = useParams<{ slug: string }>();
  const isAll = params.slug === "all";
  const [sort, setSort] = useState<SortOption>("newest");
  const [minPrice, setMinPrice] = useState("");
  const [maxPrice, setMaxPrice] = useState("");

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
    return <p className="text-muted-foreground">{t("common.loading")}</p>;
  }

  const category = data?.category;
  let items = data?.items ?? [];

  const min = minPrice ? Number(minPrice) : null;
  const max = maxPrice ? Number(maxPrice) : null;
  if (min !== null) items = items.filter((p) => p.base_price >= min);
  if (max !== null) items = items.filter((p) => p.base_price <= max);

  items = [...items].sort((a, b) => {
    if (sort === "price-asc") return a.base_price - b.base_price;
    if (sort === "price-desc") return b.base_price - a.base_price;
    return 0;
  });

  const title = isAll ? t("home.allProducts") : category?.name ?? params.slug;

  return (
    <div className="space-y-6">
      <Breadcrumb
        items={[
          { label: t("brand"), href: "/" },
          { label: t("nav.categories"), href: "/categories/all" },
          ...(isAll ? [] : [{ label: title }]),
        ]}
      />
      <h1 className="text-3xl font-bold">{title}</h1>
      <div className="grid gap-6 lg:grid-cols-[1fr_240px]">
        <div>
          {items.length === 0 ? (
            <p className="text-muted-foreground">{t("category.empty")}</p>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              {items.map((product) => (
                <Link key={product.id} href={`/products/${product.slug}`}>
                  <Card className="h-full transition-shadow hover:shadow-md">
                    <CardHeader>
                      <CardTitle className="text-base">{product.name}</CardTitle>
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
        <aside className="order-first lg:order-last">
          <Card className="sticky top-4">
            <CardHeader>
              <CardTitle className="text-base">{t("category.filter")}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">{t("category.sort")}</label>
                <select
                  value={sort}
                  onChange={(e) => setSort(e.target.value as SortOption)}
                  className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
                >
                  <option value="newest">{t("category.sortNewest")}</option>
                  <option value="price-asc">{t("category.sortPriceAsc")}</option>
                  <option value="price-desc">{t("category.sortPriceDesc")}</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">حداقل قیمت (ریال)</label>
                <input
                  type="number"
                  value={minPrice}
                  onChange={(e) => setMinPrice(e.target.value)}
                  className="ltr-field flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
                  placeholder="0"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">حداکثر قیمت (ریال)</label>
                <input
                  type="number"
                  value={maxPrice}
                  onChange={(e) => setMaxPrice(e.target.value)}
                  className="ltr-field flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm"
                  placeholder="0"
                />
              </div>
            </CardContent>
          </Card>
        </aside>
      </div>
    </div>
  );
}
