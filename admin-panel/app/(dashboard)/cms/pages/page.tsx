"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t, statusLabel } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface CmsPage {
  id: string;
  title: string;
  slug: string;
  status: string;
}

interface PageDetail {
  page: CmsPage;
  blocks: { id: string; type: string; content: string; sort_order: number }[];
}

export default function CmsPagesPage() {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [slug, setSlug] = useState("");
  const [editPage, setEditPage] = useState<CmsPage | null>(null);
  const [blocksJson, setBlocksJson] = useState("[]");
  const [seoTitle, setSeoTitle] = useState("");
  const [seoDescription, setSeoDescription] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.cmsPages,
    queryFn: async () => {
      const res = await api.get<CmsPage[]>("/cms/pages");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/cms/pages", { title, slug });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cmsPages });
      setTitle("");
      setSlug("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const loadPage = async (page: CmsPage) => {
    setEditPage(page);
    setSeoTitle("");
    setSeoDescription("");
    setBlocksJson("[]");
    try {
      const res = await api.get<PageDetail>(`/cms/pages/id/${page.id}`);
      if (res.data?.blocks) {
        setBlocksJson(JSON.stringify(res.data.blocks, null, 2));
      }
      const seo = await api.get<{ title: string; description: string }>(
        `/cms/seo?entity_type=page&entity_id=${page.id}`
      );
      if (seo.data) {
        setSeoTitle(seo.data.title ?? "");
        setSeoDescription(seo.data.description ?? "");
      }
    } catch {
      /* no seo yet */
    }
  };

  const saveMutation = useMutation({
    mutationFn: async () => {
      if (!editPage) return;
      await api.put(`/cms/pages/${editPage.id}`, {
        title: editPage.title,
        slug: editPage.slug,
        status: editPage.status,
      });
      await api.put(`/cms/pages/${editPage.id}/blocks`, {
        blocks: JSON.parse(blocksJson),
      });
      await api.post("/cms/seo", {
        entity_type: "page",
        entity_id: editPage.id,
        title: seoTitle,
        description: seoDescription,
        keywords: "",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.cmsPages });
      setEditPage(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.cmsPages.title")}</h1>
        <p className="text-muted-foreground">{t("pages.cmsPages.subtitle")}</p>
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
          <CardTitle>{t("pages.cmsPages.createPage")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Input placeholder={t("common.title")} value={title} onChange={(e) => setTitle(e.target.value)} />
          <Input placeholder={t("common.slug")} value={slug} onChange={(e) => setSlug(e.target.value)} />
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
            {t("pages.cmsPages.create")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: t("common.title") },
          { key: "slug", header: t("common.slug") },
          {
            key: "status",
            header: t("common.status"),
            render: (row) => statusLabel(String(row.status)),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <Button size="sm" variant="outline" onClick={() => loadPage(row as unknown as CmsPage)}>
                {t("common.edit")}
              </Button>
            ),
          },
        ]}
      />

      <Dialog open={!!editPage} onOpenChange={() => setEditPage(null)}>
        <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("pages.cmsPages.editPage")}</DialogTitle>
          </DialogHeader>
          {editPage && (
            <div className="space-y-3">
              <div className="space-y-2">
                <Label>{t("common.title")}</Label>
                <Input value={editPage.title} onChange={(e) => setEditPage({ ...editPage, title: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>{t("common.slug")}</Label>
                <Input value={editPage.slug} onChange={(e) => setEditPage({ ...editPage, slug: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>{t("pages.cmsPages.blocks")}</Label>
                <Textarea className="min-h-[120px] font-mono text-xs" value={blocksJson} onChange={(e) => setBlocksJson(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>{t("pages.cmsPages.seoTitle")}</Label>
                <Input value={seoTitle} onChange={(e) => setSeoTitle(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label>{t("pages.cmsPages.seoDescription")}</Label>
                <Textarea value={seoDescription} onChange={(e) => setSeoDescription(e.target.value)} />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditPage(null)}>{t("common.cancel")}</Button>
            <Button onClick={() => saveMutation.mutate()} disabled={saveMutation.isPending}>{t("common.save")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
