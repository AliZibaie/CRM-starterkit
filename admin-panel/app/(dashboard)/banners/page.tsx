"use client";

import { useRef, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api, apiUpload } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { formatJalaliDateTime } from "@/lib/date";
import { datetimeLocalToRFC3339, rfc3339ToDatetimeLocal } from "@/lib/datetime-input";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Banner {
  id: string;
  title: string;
  image_url: string;
  link_url?: string;
  position: string;
  is_active: boolean;
  starts_at?: string;
  ends_at?: string;
}

interface MediaUpload {
  url: string;
}

export default function BannersPage() {
  const queryClient = useQueryClient();
  const fileRef = useRef<HTMLInputElement>(null);
  const [title, setTitle] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [linkUrl, setLinkUrl] = useState("");
  const [position, setPosition] = useState("");
  const [editBanner, setEditBanner] = useState<Banner | null>(null);
  const [scheduleBanner, setScheduleBanner] = useState<Banner | null>(null);
  const [scheduleStart, setScheduleStart] = useState("");
  const [scheduleEnd, setScheduleEnd] = useState("");
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: banners = [], isLoading, isError, error, refetch } = useQuery({
    queryKey: queryKeys.banners,
    queryFn: async () => {
      const res = await api.get<Banner[]>("/banners");
      return res.data ?? [];
    },
  });

  const { data: positions = [] } = useQuery({
    queryKey: queryKeys.bannerPositions,
    queryFn: async () => {
      const res = await api.get<string[]>("/banners/positions");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/banners", {
        title,
        image_url: imageUrl,
        link_url: linkUrl,
        position,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.banners });
      setTitle("");
      setImageUrl("");
      setLinkUrl("");
      setPosition("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const updateMutation = useMutation({
    mutationFn: async (banner: Banner) => {
      await api.put(`/banners/${banner.id}`, {
        title: banner.title,
        image_url: banner.image_url,
        link_url: banner.link_url ?? "",
        position: banner.position,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.banners });
      setEditBanner(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const scheduleMutation = useMutation({
    mutationFn: async () => {
      if (!scheduleBanner) return;
      const starts = datetimeLocalToRFC3339(scheduleStart);
      const ends = datetimeLocalToRFC3339(scheduleEnd);
      await api.put(`/banners/${scheduleBanner.id}/schedule`, {
        starts_at: starts,
        ends_at: ends,
        is_active: true,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.banners });
      setScheduleBanner(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/banners/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.banners });
      setDeleteId(null);
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => {
      setDeleteId(null);
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    },
  });

  const handleUpload = async (file: File) => {
    setUploading(true);
    setFeedback(null);
    try {
      const form = new FormData();
      form.append("file", file);
      const res = await apiUpload<MediaUpload>("/media/upload", form);
      if (res.data?.url) {
        setImageUrl(res.data.url);
        if (editBanner) {
          setEditBanner({ ...editBanner, image_url: res.data.url });
        }
      }
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    } finally {
      setUploading(false);
    }
  };

  const openSchedule = (banner: Banner) => {
    setScheduleBanner(banner);
    setScheduleStart(rfc3339ToDatetimeLocal(banner.starts_at));
    setScheduleEnd(rfc3339ToDatetimeLocal(banner.ends_at));
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.banners.title")}</h1>
        <p className="text-muted-foreground">{t("pages.banners.subtitle")}</p>
      </div>

      {feedback && (
        <div
          className={`rounded-xl border px-4 py-3 text-sm ${
            feedback.type === "error"
              ? "border-destructive/40 bg-destructive/10 text-destructive"
              : "border-emerald-500/40 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
          }`}
        >
          {feedback.text}
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {getApiErrorMessage(error)}
          <Button size="sm" variant="outline" className="ms-3" onClick={() => refetch()}>
            {t("common.refresh")}
          </Button>
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.banners.createBanner")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>{t("common.title")}</Label>
            <Input value={title} onChange={(e) => setTitle(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.position")}</Label>
            <Select value={position} onValueChange={setPosition}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectPosition")} />
              </SelectTrigger>
              <SelectContent>
                {positions.map((p) => (
                  <SelectItem key={p} value={p}>
                    {p}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>{t("common.linkUrl")}</Label>
            <Input value={linkUrl} onChange={(e) => setLinkUrl(e.target.value)} className="ltr-field" />
          </div>
          <div className="space-y-2">
            <Label>{t("common.imageUrl")}</Label>
            <Input value={imageUrl} onChange={(e) => setImageUrl(e.target.value)} className="ltr-field" />
          </div>
          <div className="flex flex-wrap items-center gap-3">
            <input
              ref={fileRef}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (file) handleUpload(file);
              }}
            />
            <Button
              type="button"
              variant="outline"
              disabled={uploading}
              onClick={() => fileRef.current?.click()}
            >
              {uploading ? t("common.uploading") : t("common.uploadImage")}
            </Button>
          </div>
          {imageUrl && (
            <img
              src={imageUrl}
              alt={title || "preview"}
              className="max-h-40 rounded-lg border border-border object-contain"
            />
          )}
          <Button
            onClick={() => createMutation.mutate()}
            disabled={createMutation.isPending || !title || !imageUrl || !position}
          >
            {createMutation.isPending ? t("common.saving") : t("pages.banners.create")}
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={banners as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: t("common.title") },
          { key: "position", header: t("common.position") },
          {
            key: "preview",
            header: t("common.preview"),
            render: (row) =>
              row.image_url ? (
                <img
                  src={String(row.image_url)}
                  alt=""
                  className="h-12 w-20 rounded object-cover"
                />
              ) : (
                "—"
              ),
          },
          {
            key: "starts_at",
            header: t("pages.banners.startsAt"),
            render: (row) =>
              row.starts_at ? formatJalaliDateTime(row.starts_at as string) : "—",
          },
          {
            key: "is_active",
            header: t("common.active"),
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? t("common.yes") : t("common.no")}
              </Badge>
            ),
          },
          {
            key: "actions",
            header: t("common.actions"),
            render: (row) => {
              const banner = row as unknown as Banner;
              return (
                <div className="flex flex-wrap gap-2">
                  <Button size="sm" variant="outline" onClick={() => setEditBanner(banner)}>
                    {t("common.edit")}
                  </Button>
                  <Button size="sm" variant="secondary" onClick={() => openSchedule(banner)}>
                    {t("common.schedule")}
                  </Button>
                  <Button size="sm" variant="destructive" onClick={() => setDeleteId(banner.id)}>
                    {t("common.delete")}
                  </Button>
                </div>
              );
            },
          },
        ]}
      />

      {editBanner && (
        <Card>
          <CardHeader>
            <CardTitle>{t("common.edit")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.title")}</Label>
              <Input
                value={editBanner.title}
                onChange={(e) => setEditBanner({ ...editBanner, title: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("common.position")}</Label>
              <Select
                value={editBanner.position}
                onValueChange={(v) => setEditBanner({ ...editBanner, position: v })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {positions.map((p) => (
                    <SelectItem key={p} value={p}>
                      {p}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.imageUrl")}</Label>
              <Input
                className="ltr-field"
                value={editBanner.image_url}
                onChange={(e) => setEditBanner({ ...editBanner, image_url: e.target.value })}
              />
            </div>
            {editBanner.image_url && (
              <img
                src={editBanner.image_url}
                alt=""
                className="max-h-40 rounded-lg border border-border object-contain"
              />
            )}
            <div className="flex gap-2">
              <Button onClick={() => updateMutation.mutate(editBanner)} disabled={updateMutation.isPending}>
                {updateMutation.isPending ? t("common.saving") : t("common.save")}
              </Button>
              <Button variant="ghost" onClick={() => setEditBanner(null)}>
                {t("common.cancel")}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {scheduleBanner && (
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.banners.scheduleBanner")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <p className="text-sm text-muted-foreground">{scheduleBanner.title}</p>
            <div className="space-y-2">
              <Label>{t("pages.banners.startsAt")}</Label>
              <Input
                type="datetime-local"
                className="ltr-field"
                value={scheduleStart}
                onChange={(e) => setScheduleStart(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>{t("pages.banners.endsAt")}</Label>
              <Input
                type="datetime-local"
                className="ltr-field"
                value={scheduleEnd}
                onChange={(e) => setScheduleEnd(e.target.value)}
              />
            </div>
            <div className="flex gap-2">
              <Button onClick={() => scheduleMutation.mutate()} disabled={scheduleMutation.isPending}>
                {scheduleMutation.isPending ? t("common.saving") : t("common.schedule")}
              </Button>
              <Button variant="ghost" onClick={() => setScheduleBanner(null)}>
                {t("common.cancel")}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("common.confirmDelete")}</DialogTitle>
            <DialogDescription>{t("common.confirmDeleteMessage")}</DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
            >
              {deleteMutation.isPending ? t("common.saving") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
