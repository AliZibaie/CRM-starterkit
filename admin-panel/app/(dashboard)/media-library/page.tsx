"use client";

import { useMemo, useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api, apiUpload } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

interface MediaItem {
  id: string;
  filename: string;
  mime_type: string;
  size: number;
  url: string;
  thumbnail_url?: string;
  created_at: string;
}

export default function MediaLibraryPage() {
  const queryClient = useQueryClient();
  const fileRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [search, setSearch] = useState("");
  const [gridView, setGridView] = useState(true);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.media,
    queryFn: async () => {
      const res = await api.get<MediaItem[]>("/media?scope=all");
      return res.data ?? [];
    },
  });

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return data;
    return data.filter((m) => m.filename?.toLowerCase().includes(q));
  }, [data, search]);

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/media/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.media });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const handleUpload = async () => {
    const files = fileRef.current?.files;
    if (!files?.length) return;
    setUploading(true);
    setFeedback(null);
    try {
      for (const file of Array.from(files)) {
        const formData = new FormData();
        formData.append("file", file);
        await apiUpload("/media/upload", formData);
      }
      queryClient.invalidateQueries({ queryKey: queryKeys.media });
      if (fileRef.current) fileRef.current.value = "";
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.media.title")}</h1>
          <p className="text-muted-foreground">{t("pages.media.subtitle")}</p>
        </div>
        <div className="flex gap-2">
          <Button variant={gridView ? "default" : "outline"} size="sm" onClick={() => setGridView(true)}>
            Grid
          </Button>
          <Button variant={!gridView ? "default" : "outline"} size="sm" onClick={() => setGridView(false)}>
            List
          </Button>
        </div>
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
          <CardTitle>{t("pages.media.uploadFile")}</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-wrap items-center gap-3">
          <input ref={fileRef} type="file" multiple className="text-sm" />
          <Button onClick={handleUpload} disabled={uploading}>
            {uploading ? t("common.uploading") : t("common.upload")}
          </Button>
          <Input
            className="max-w-xs"
            placeholder={t("common.search")}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </CardContent>
      </Card>

      {gridView ? (
        <div className="grid gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {filtered.map((item) => (
            <Card key={item.id}>
              <CardContent className="space-y-2 p-3">
                {(item.thumbnail_url || item.mime_type?.startsWith("image/")) && (
                  <img
                    src={item.thumbnail_url || item.url}
                    alt={item.filename}
                    className="h-32 w-full rounded object-cover"
                  />
                )}
                <p className="truncate text-sm font-medium">{item.filename}</p>
                <p className="text-xs text-muted-foreground">{Math.round(item.size / 1024)} KB</p>
                <Button size="sm" variant="destructive" onClick={() => deleteMutation.mutate(item.id)}>
                  {t("common.delete")}
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <DataTable
          loading={isLoading}
          data={filtered as unknown as Record<string, unknown>[]}
          keyField="id"
          columns={[
            {
              key: "preview",
              header: t("common.preview"),
              render: (row) => {
                const item = row as unknown as MediaItem;
                if (!item.mime_type?.startsWith("image/")) return "—";
                return (
                  <img
                    src={item.thumbnail_url || item.url}
                    alt=""
                    className="h-10 w-10 rounded object-cover"
                  />
                );
              },
            },
            { key: "filename", header: t("common.filename") },
            { key: "mime_type", header: t("common.type") },
            {
              key: "size",
              header: t("common.size"),
              render: (row) => `${Math.round((row.size as number) / 1024)} KB`,
            },
            {
              key: "created_at",
              header: t("common.createdAt"),
              render: (row) => formatJalaliDateTime(row.created_at as string),
            },
            {
              key: "actions",
              header: t("common.actions"),
              render: (row) => (
                <Button size="sm" variant="destructive" onClick={() => deleteMutation.mutate(row.id as string)}>
                  {t("common.delete")}
                </Button>
              ),
            },
          ]}
        />
      )}
    </div>
  );
}
