"use client";

import { useMemo, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
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

interface Category {
  id: string;
  name: string;
  parent_id?: string;
}

function categoryPath(categories: Category[], id: string): string {
  const cat = categories.find((c) => c.id === id);
  if (!cat) return "";
  if (cat.parent_id) {
    const parent = categoryPath(categories, cat.parent_id);
    return parent ? `${parent} / ${cat.name}` : cat.name;
  }
  return cat.name;
}

export default function CategoriesPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [parentId, setParentId] = useState("");
  const [editCategory, setEditCategory] = useState<Category | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.categories,
    queryFn: async () => {
      const res = await api.get<Category[]>("/catalog/categories");
      return res.data ?? [];
    },
  });

  const sorted = useMemo(() => {
    return [...data].sort((a, b) => categoryPath(data, a.id).localeCompare(categoryPath(data, b.id)));
  }, [data]);

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/catalog/categories", { name, parent_id: parentId });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories });
      setName("");
      setParentId("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async (cat: Category) => {
      await api.put(`/catalog/categories/${cat.id}`, {
        name: cat.name,
        parent_id: cat.parent_id ?? "",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories });
      setEditCategory(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/catalog/categories/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.categories });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.categories.title")}</h1>
        <p className="text-muted-foreground">{t("pages.categories.subtitle")}</p>
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
          <CardTitle>{t("pages.categories.createCategory")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>{t("common.name")}</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("pages.categories.parentCategory")}</Label>
            <Select value={parentId || "__none"} onValueChange={(v) => setParentId(v === "__none" ? "" : v)}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.select")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="__none">{t("pages.categories.noParent")}</SelectItem>
                {data.map((c) => (
                  <SelectItem key={c.id} value={c.id}>
                    {categoryPath(data, c.id)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending || !name.trim()}>
            {t("common.create")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={sorted as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          {
            key: "name",
            header: t("common.name"),
            render: (row) => categoryPath(data, row.id as string),
          },
          {
            key: "parent_id",
            header: t("common.parentId"),
            render: (row) => {
              const pid = row.parent_id as string;
              return pid ? categoryPath(data, pid) : "—";
            },
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setEditCategory(row as unknown as Category)}
                >
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

      <Dialog open={!!editCategory} onOpenChange={() => setEditCategory(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.edit")}</DialogTitle>
          </DialogHeader>
          {editCategory && (
            <div className="space-y-3">
              <div className="space-y-2">
                <Label>{t("common.name")}</Label>
                <Input
                  value={editCategory.name}
                  onChange={(e) => setEditCategory({ ...editCategory, name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label>{t("pages.categories.parentCategory")}</Label>
                <Select
                  value={editCategory.parent_id || "__none"}
                  onValueChange={(v) =>
                    setEditCategory({ ...editCategory, parent_id: v === "__none" ? "" : v })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none">{t("pages.categories.noParent")}</SelectItem>
                    {data
                      .filter((c) => c.id !== editCategory.id)
                      .map((c) => (
                        <SelectItem key={c.id} value={c.id}>
                          {categoryPath(data, c.id)}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditCategory(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              onClick={() => editCategory && updateMutation.mutate(editCategory)}
              disabled={updateMutation.isPending}
            >
              {t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

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
