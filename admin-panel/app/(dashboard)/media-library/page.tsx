"use client";

import { useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api, apiUpload } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface MediaItem {
  id: string;
  filename: string;
  mime_type: string;
  size: number;
  url: string;
  created_at: string;
}

export default function MediaLibraryPage() {
  const queryClient = useQueryClient();
  const fileRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.media,
    queryFn: async () => {
      const res = await api.get<MediaItem[]>("/media");
      return res.data ?? [];
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/media/${id}`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.media }),
  });

  const handleUpload = async () => {
    const file = fileRef.current?.files?.[0];
    if (!file) return;
    setUploading(true);
    try {
      const formData = new FormData();
      formData.append("file", file);
      await apiUpload("/media/upload", formData);
      queryClient.invalidateQueries({ queryKey: queryKeys.media });
      if (fileRef.current) fileRef.current.value = "";
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.media.title")}</h1>
        <p className="text-muted-foreground">{t("pages.media.subtitle")}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.media.uploadFile")}</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center gap-3">
          <input ref={fileRef} type="file" className="text-sm" />
          <Button onClick={handleUpload} disabled={uploading}>
            {uploading ? t("common.uploading") : t("common.upload")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "filename", header: t("common.filename") },
          { key: "mime_type", header: t("common.type") },
          {
            key: "size",
            header: t("common.size"),
            render: (row) => `${Math.round((row.size as number) / 1024)} KB`,
          },
          { key: "url", header: t("common.url") },
          {
            key: "created_at",
            header: t("common.createdAt"),
            render: (row) => formatJalaliDateTime(row.created_at as string),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => (
              <Button
                size="sm"
                variant="destructive"
                onClick={() => deleteMutation.mutate(row.id as string)}
              >
                {t("common.delete")}
              </Button>
            ),
          },
        ]}
      />
    </div>
  );
}
