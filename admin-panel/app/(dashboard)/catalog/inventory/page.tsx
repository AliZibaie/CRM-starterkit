"use client";

import { useEffect, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Product {
  id: string;
  name: string;
  type: string;
}

interface Variant {
  id: string;
  sku: string;
  name: string;
  price: number;
}

interface ProductDetail {
  item: Product;
  variants: Variant[];
}

interface Inventory {
  variant_id: string;
  quantity: number;
}

export default function InventoryPage() {
  const queryClient = useQueryClient();
  const [selectedProduct, setSelectedProduct] = useState("");
  const [selectedVariant, setSelectedVariant] = useState("");
  const [quantity, setQuantity] = useState(0);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: products = [] } = useQuery({
    queryKey: queryKeys.products,
    queryFn: async () => {
      const res = await api.get<Product[]>("/catalog/items");
      return res.data ?? [];
    },
  });

  const { data: productDetail } = useQuery({
    queryKey: queryKeys.product(selectedProduct),
    queryFn: async () => {
      const res = await api.get<ProductDetail>(`/catalog/items/${selectedProduct}`);
      return res.data!;
    },
    enabled: !!selectedProduct,
  });

  const { data: inventory, refetch: refetchInventory } = useQuery({
    queryKey: queryKeys.inventory(selectedVariant),
    queryFn: async () => {
      const res = await api.get<Inventory>(`/catalog/inventory/${selectedVariant}`);
      return res.data;
    },
    enabled: !!selectedVariant,
  });

  useEffect(() => {
    if (inventory) {
      setQuantity(inventory.quantity);
    }
  }, [inventory]);

  const updateMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/catalog/inventory/${selectedVariant}`, { quantity });
    },
    onSuccess: async () => {
      await refetchInventory();
      queryClient.invalidateQueries({ queryKey: queryKeys.inventory(selectedVariant) });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const variants = productDetail?.variants ?? [];
  const physicalProducts = products.filter((p) => p.type === "physical_product");

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.inventory.title")}</h1>
        <p className="text-muted-foreground">{t("pages.inventory.subtitle")}</p>
      </div>

      {feedback && (
        <div
          className={`rounded-md border px-4 py-3 text-sm ${
            feedback.type === "error"
              ? "border-destructive/50 bg-destructive/10 text-destructive"
              : "border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400"
          }`}
        >
          {feedback.text}
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.inventory.updateInventory")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>{t("common.product")}</Label>
              <Select
                value={selectedProduct}
                onValueChange={(v) => {
                  setSelectedProduct(v);
                  setSelectedVariant("");
                  setFeedback(null);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t("common.selectProduct")} />
                </SelectTrigger>
                <SelectContent>
                  {physicalProducts.map((p) => (
                    <SelectItem key={p.id} value={p.id}>
                      {p.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.variant")}</Label>
              <Select
                value={selectedVariant}
                onValueChange={(v) => {
                  setSelectedVariant(v);
                  setFeedback(null);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t("common.selectVariant")} />
                </SelectTrigger>
                <SelectContent>
                  {variants.map((v) => (
                    <SelectItem key={v.id} value={v.id}>
                      {v.sku} — {v.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {inventory && (
            <p className="text-sm text-muted-foreground">
              {t("common.currentStock", { quantity: inventory.quantity, reserved: 0 })}
            </p>
          )}

          <div className="space-y-2">
            <Label>{t("common.quantity")}</Label>
            <Input
              type="number"
              min={0}
              value={quantity}
              onChange={(e) => setQuantity(Number(e.target.value))}
            />
          </div>

          <Button
            onClick={() => updateMutation.mutate()}
            disabled={updateMutation.isPending || !selectedVariant}
          >
            {updateMutation.isPending ? t("common.saving") : t("common.updateStock")}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
