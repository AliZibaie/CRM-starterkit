"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/table/data-table";

interface Menu {
  id: string;
  name: string;
  location?: string;
  items_json: string;
}

interface MenuItem {
  label: string;
  url: string;
}

export default function CmsMenusPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [items, setItems] = useState<MenuItem[]>([{ label: "", url: "" }]);
  const [editMenu, setEditMenu] = useState<Menu | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data: menus = [], isLoading } = useQuery({
    queryKey: queryKeys.cmsMenus,
    queryFn: async () => {
      const res = await api.get<Menu[]>("/cms/menus/all");
      return res.data ?? [];
    },
  });

  const saveMutation = useMutation({
    mutationFn: async (payload: { id?: string; name: string; items: MenuItem[] }) => {
      await api.post("/cms/menus", {
        id: payload.id ?? "",
        name: payload.name,
        items_json: JSON.stringify(payload.items.filter((i) => i.label.trim())),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cmsMenus });
      setName("");
      setItems([{ label: "", url: "" }]);
      setEditMenu(null);
      setFeedback({ type: "success", text: t("common.menuSavedSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const startEdit = (menu: Menu) => {
    setEditMenu(menu);
    try {
      setItems(JSON.parse(menu.items_json) as MenuItem[]);
    } catch {
      setItems([{ label: "", url: "" }]);
    }
  };

  const moveItem = (index: number, dir: -1 | 1) => {
    const next = [...items];
    const target = index + dir;
    if (target < 0 || target >= next.length) return;
    [next[index], next[target]] = [next[target], next[index]];
    setItems(next);
  };

  const currentItems = editMenu ? items : items;
  const currentName = editMenu ? editMenu.name : name;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.cmsMenus.title")}</h1>
        <p className="text-muted-foreground">{t("pages.cmsMenus.subtitle")}</p>
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
          <CardTitle>{editMenu ? t("common.edit") : t("pages.cmsMenus.upsertMenu")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!editMenu && (
            <div className="space-y-2">
              <Label>{t("common.menuName")}</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} />
            </div>
          )}
          {currentItems.map((item, i) => (
            <div key={i} className="flex gap-2">
              <Input
                placeholder={t("common.title")}
                value={item.label}
                onChange={(e) =>
                  setItems((p) => p.map((x, j) => (j === i ? { ...x, label: e.target.value } : x)))
                }
              />
              <Input
                placeholder="URL"
                value={item.url}
                onChange={(e) =>
                  setItems((p) => p.map((x, j) => (j === i ? { ...x, url: e.target.value } : x)))
                }
              />
              <Button type="button" size="sm" variant="outline" onClick={() => moveItem(i, -1)}>
                ↑
              </Button>
              <Button type="button" size="sm" variant="outline" onClick={() => moveItem(i, 1)}>
                ↓
              </Button>
              <Button
                type="button"
                size="sm"
                variant="outline"
                onClick={() => setItems((p) => p.filter((_, j) => j !== i))}
              >
                −
              </Button>
            </div>
          ))}
          <div className="flex gap-2">
            <Button type="button" variant="outline" onClick={() => setItems((p) => [...p, { label: "", url: "" }])}>
              +
            </Button>
            <Button
              onClick={() =>
                saveMutation.mutate({
                  id: editMenu?.id,
                  name: editMenu?.name ?? name,
                  items: currentItems,
                })
              }
              disabled={saveMutation.isPending || !(editMenu?.name || name.trim())}
            >
              {t("pages.cmsMenus.save")}
            </Button>
            {editMenu && (
              <Button variant="outline" onClick={() => { setEditMenu(null); setItems([{ label: "", url: "" }]); }}>
                {t("common.cancel")}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={menus as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.menuName") },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <Button size="sm" variant="outline" onClick={() => startEdit(row as unknown as Menu)}>
                {t("common.edit")}
              </Button>
            ),
          },
        ]}
      />
    </div>
  );
}
