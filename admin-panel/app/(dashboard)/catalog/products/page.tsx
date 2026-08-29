"use client";

import { useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api, apiUpload } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

const PRODUCT_TYPES = [
  { value: "physical_product", label: "فیزیکی" },
  { value: "service", label: "خدمت" },
  { value: "course_license", label: "لایسنس دوره" },
  { value: "digital", label: "دیجیتال" },
];

interface Category {
  id: string;
  name: string;
}

interface Product {
  id: string;
  name: string;
  type: string;
  base_price: number;
  currency: string;
  category_id?: string;
  description?: string;
  is_active?: boolean;
}

interface Attribute {
  id: string;
  key: string;
  value: string;
}

interface Variant {
  id: string;
  sku: string;
  name: string;
  price: number;
}

interface ProductDetail {
  item: Product;
  attributes: Attribute[];
  variants: Variant[];
}

interface AttrRow {
  key: string;
  value: string;
}

interface VariantRow {
  sku: string;
  name: string;
  price: number;
  stock: number;
}

interface MediaUpload {
  url: string;
}

function isImageKey(key: string) {
  return key === "image" || key.startsWith("image_");
}

export default function ProductsPage() {
  const queryClient = useQueryClient();
  const fileRef = useRef<HTMLInputElement>(null);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [type, setType] = useState("physical_product");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [currency, setCurrency] = useState("IRR");
  const [basePrice, setBasePrice] = useState(0);
  const [attrs, setAttrs] = useState<AttrRow[]>([{ key: "", value: "" }]);
  const [variants, setVariants] = useState<VariantRow[]>([
    { sku: "", name: "", price: 0, stock: 0 },
  ]);
  const [imageUrls, setImageUrls] = useState<string[]>([]);
  const [uploading, setUploading] = useState(false);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const isPhysical = type === "physical_product";

  const { data: categories = [] } = useQuery({
    queryKey: queryKeys.categories,
    queryFn: async () => {
      const res = await api.get<Category[]>("/catalog/categories");
      return res.data ?? [];
    },
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.products,
    queryFn: async () => {
      const res = await api.get<Product[]>("/catalog/items");
      return res.data ?? [];
    },
  });

  const resetForm = () => {
    setEditId(null);
    setType("physical_product");
    setName("");
    setDescription("");
    setCategoryId("");
    setCurrency("IRR");
    setBasePrice(0);
    setAttrs([{ key: "", value: "" }]);
    setVariants([{ sku: "", name: "", price: 0, stock: 0 }]);
    setImageUrls([]);
    setShowForm(false);
  };

  const loadForEdit = async (id: string) => {
    setFeedback(null);
    const res = await api.get<ProductDetail>(`/catalog/items/${id}`);
    const detail = res.data!;
    setEditId(id);
    setType(detail.item.type);
    setName(detail.item.name);
    setDescription(detail.item.description ?? "");
    setCategoryId(detail.item.category_id ?? "");
    setCurrency(detail.item.currency || "IRR");
    setBasePrice(detail.item.base_price);
    const customAttrs = detail.attributes.filter((a) => !isImageKey(a.key));
    setAttrs(customAttrs.length ? customAttrs.map((a) => ({ key: a.key, value: a.value })) : [{ key: "", value: "" }]);
    setImageUrls(detail.attributes.filter((a) => isImageKey(a.key)).map((a) => a.value));
    setVariants(
      detail.variants.length
        ? detail.variants.map((v) => ({ sku: v.sku, name: v.name, price: v.price, stock: 0 }))
        : [{ sku: "", name: "", price: 0, stock: 0 }]
    );
    setShowForm(true);
  };

  const saveMutation = useMutation({
    mutationFn: async () => {
      const payload = {
        type,
        name,
        description,
        category_id: categoryId,
        currency,
        base_price: basePrice,
      };
      let itemId = editId;
      if (editId) {
        await api.put(`/catalog/items/${editId}`, { ...payload, is_active: true });
      } else {
        const res = await api.post<Product>("/catalog/items", payload);
        itemId = res.data!.id;
      }
      if (!itemId) return;

      const attrRows = attrs.filter((a) => a.key.trim());
      for (let i = 0; i < imageUrls.length; i++) {
        await api.post(`/catalog/items/${itemId}/attributes`, {
          key: `image_${i}`,
          value: imageUrls[i],
        });
      }
      for (const row of attrRows) {
        await api.post(`/catalog/items/${itemId}/attributes`, {
          key: row.key.trim(),
          value: row.value,
        });
      }

      if (!editId) {
        for (const v of variants.filter((x) => x.sku.trim() && x.name.trim())) {
          const vr = await api.post<Variant>(`/catalog/items/${itemId}/variants`, {
            sku: v.sku.trim(),
            name: v.name.trim(),
            price: v.price,
            attrs_json: "{}",
          });
          if (isPhysical && v.stock > 0 && vr.data?.id) {
            await api.put(`/catalog/inventory/${vr.data.id}`, { quantity: v.stock });
          }
        }
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.products });
      resetForm();
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/catalog/items/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.products });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const handleUpload = async (files: FileList | null) => {
    if (!files?.length) return;
    setUploading(true);
    setFeedback(null);
    try {
      const urls: string[] = [];
      for (const file of Array.from(files)) {
        const fd = new FormData();
        fd.append("file", file);
        const res = await apiUpload<MediaUpload>("/media/upload", fd);
        if (res.data?.url) urls.push(res.data.url);
      }
      setImageUrls((prev) => [...prev, ...urls]);
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = "";
    }
  };

  const typeLabel = (v: string) => PRODUCT_TYPES.find((p) => p.value === v)?.label ?? v;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.products.title")}</h1>
          <p className="text-muted-foreground">{t("pages.products.subtitle")}</p>
        </div>
        <Button
          onClick={() => {
            if (showForm && !editId) resetForm();
            else {
              resetForm();
              setShowForm(true);
            }
          }}
        >
          {showForm ? t("common.cancel") : t("pages.products.addProduct")}
        </Button>
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

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>
              {editId ? t("common.edit") : t("pages.products.createProduct")}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label>{t("common.type")}</Label>
                <Select value={type} onValueChange={setType}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PRODUCT_TYPES.map((p) => (
                      <SelectItem key={p.value} value={p.value}>
                        {p.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>{t("common.name")}</Label>
                <Input value={name} onChange={(e) => setName(e.target.value)} />
              </div>
              <div className="space-y-2 md:col-span-2">
                <Label>{t("common.description")}</Label>
                <Input value={description} onChange={(e) => setDescription(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>{t("common.category")}</Label>
                <Select value={categoryId || "__none"} onValueChange={(v) => setCategoryId(v === "__none" ? "" : v)}>
                  <SelectTrigger>
                    <SelectValue placeholder={t("common.select")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none">{t("common.optional")}</SelectItem>
                    {categories.map((c) => (
                      <SelectItem key={c.id} value={c.id}>
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>{t("common.currency")}</Label>
                <Input value={currency} onChange={(e) => setCurrency(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>{t("common.basePrice")}</Label>
                <Input
                  type="number"
                  value={basePrice}
                  onChange={(e) => setBasePrice(Number(e.target.value))}
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label>{t("pages.products.images")}</Label>
              <Input
                ref={fileRef}
                type="file"
                accept="image/*"
                multiple
                disabled={uploading}
                onChange={(e) => handleUpload(e.target.files)}
              />
              <div className="flex flex-wrap gap-2">
                {imageUrls.map((url, i) => (
                  <div key={url} className="relative">
                    <img src={url} alt="" className="h-16 w-16 rounded border object-cover" />
                    <Button
                      type="button"
                      size="sm"
                      variant="destructive"
                      className="absolute -end-2 -top-2 h-6 w-6 rounded-full p-0 text-xs"
                      onClick={() => setImageUrls((p) => p.filter((_, j) => j !== i))}
                    >
                      ×
                    </Button>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <Label>{t("pages.products.attributes")}</Label>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  onClick={() => setAttrs((p) => [...p, { key: "", value: "" }])}
                >
                  +
                </Button>
              </div>
              {attrs.map((row, i) => (
                <div key={i} className="flex gap-2">
                  <Input
                    placeholder={t("common.key")}
                    value={row.key}
                    onChange={(e) =>
                      setAttrs((p) => p.map((r, j) => (j === i ? { ...r, key: e.target.value } : r)))
                    }
                  />
                  <Input
                    placeholder={t("common.value")}
                    value={row.value}
                    onChange={(e) =>
                      setAttrs((p) => p.map((r, j) => (j === i ? { ...r, value: e.target.value } : r)))
                    }
                  />
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    disabled={attrs.length <= 1}
                    onClick={() => setAttrs((p) => p.filter((_, j) => j !== i))}
                  >
                    −
                  </Button>
                </div>
              ))}
            </div>

            {!editId && (
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label>{t("pages.products.variants")}</Label>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() =>
                      setVariants((p) => [...p, { sku: "", name: "", price: 0, stock: 0 }])
                    }
                  >
                    +
                  </Button>
                </div>
                {variants.map((row, i) => (
                  <div key={i} className="grid gap-2 md:grid-cols-5">
                    <Input
                      placeholder="SKU"
                      value={row.sku}
                      onChange={(e) =>
                        setVariants((p) =>
                          p.map((r, j) => (j === i ? { ...r, sku: e.target.value } : r))
                        )
                      }
                    />
                    <Input
                      placeholder={t("common.name")}
                      value={row.name}
                      onChange={(e) =>
                        setVariants((p) =>
                          p.map((r, j) => (j === i ? { ...r, name: e.target.value } : r))
                        )
                      }
                    />
                    <Input
                      type="number"
                      placeholder={t("common.price")}
                      value={row.price}
                      onChange={(e) =>
                        setVariants((p) =>
                          p.map((r, j) => (j === i ? { ...r, price: Number(e.target.value) } : r))
                        )
                      }
                    />
                    {isPhysical && (
                      <Input
                        type="number"
                        placeholder={t("pages.products.initialStock")}
                        value={row.stock}
                        onChange={(e) =>
                          setVariants((p) =>
                            p.map((r, j) => (j === i ? { ...r, stock: Number(e.target.value) } : r))
                          )
                        }
                      />
                    )}
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      disabled={variants.length <= 1}
                      onClick={() => setVariants((p) => p.filter((_, j) => j !== i))}
                    >
                      −
                    </Button>
                  </div>
                ))}
              </div>
            )}

            <Button onClick={() => saveMutation.mutate()} disabled={saveMutation.isPending || !name.trim()}>
              {saveMutation.isPending ? t("common.saving") : t("common.save")}
            </Button>
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          {
            key: "type",
            header: t("common.type"),
            render: (row) => typeLabel(String(row.type)),
          },
          {
            key: "base_price",
            header: t("common.price"),
            render: (row) =>
              `${row.currency} ${((row.base_price as number) / 100).toFixed(2)}`,
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button size="sm" variant="outline" onClick={() => loadForEdit(row.id as string)}>
                  {t("common.edit")}
                </Button>
                <Button size="sm" variant="destructive" onClick={() => setDeleteId(row.id as string)}>
                  {t("common.delete")}
                </Button>
              </div>
            ),
          },
        ]}
      />

      <Dialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">{t("common.confirmDeleteMessage")}</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
              disabled={deleteMutation.isPending}
            >
              {t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
