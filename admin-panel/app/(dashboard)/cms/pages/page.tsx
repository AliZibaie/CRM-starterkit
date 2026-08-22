"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface CmsPage {
  id: string;
  title: string;
  slug: string;
  status: string;
}

export default function CmsPagesPage() {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [slug, setSlug] = useState("");
  const [blocksPageId, setBlocksPageId] = useState("");
  const [blocksJson, setBlocksJson] = useState("[]");

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
    },
  });

  const blocksMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/cms/pages/${blocksPageId}/blocks`, {
        blocks: JSON.parse(blocksJson),
      });
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">CMS Pages</h1>
        <p className="text-muted-foreground">Content page management</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Page</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>Title</Label>
            <Input value={title} onChange={(e) => setTitle(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Slug</Label>
            <Input value={slug} onChange={(e) => setSlug(e.target.value)} />
          </div>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
            Create Page
          </Button>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Set Page Blocks</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>Page ID</Label>
            <Input value={blocksPageId} onChange={(e) => setBlocksPageId(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Blocks JSON</Label>
            <Input value={blocksJson} onChange={(e) => setBlocksJson(e.target.value)} />
          </div>
          <Button onClick={() => blocksMutation.mutate()} disabled={blocksMutation.isPending}>
            Save Blocks
          </Button>
        </CardContent>
      </Card>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: "Title" },
          { key: "slug", header: "Slug" },
          { key: "status", header: "Status" },
        ]}
      />
    </div>
  );
}
