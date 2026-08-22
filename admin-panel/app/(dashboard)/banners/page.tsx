"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
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

interface Banner {
  id: string;
  title: string;
  image_url: string;
  position: string;
  is_active: boolean;
}

export default function BannersPage() {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [scheduleId, setScheduleId] = useState("");
  const [scheduleStart, setScheduleStart] = useState("");
  const [scheduleEnd, setScheduleEnd] = useState("");
  const [positionId, setPositionId] = useState("");
  const [position, setPosition] = useState("");

  const { data: banners = [], isLoading } = useQuery({
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
      await api.post("/banners", { title, image_url: imageUrl });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.banners });
      setTitle("");
      setImageUrl("");
    },
  });

  const scheduleMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/banners/${scheduleId}/schedule`, {
        starts_at: scheduleStart,
        ends_at: scheduleEnd,
      });
    },
  });

  const positionMutation = useMutation({
    mutationFn: async () => {
      await api.put(`/banners/${positionId}/position`, { position });
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.banners }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Banners</h1>
        <p className="text-muted-foreground">Manage promotional banners</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Banner</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <Label>Title</Label>
            <Input value={title} onChange={(e) => setTitle(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Image URL</Label>
            <Input value={imageUrl} onChange={(e) => setImageUrl(e.target.value)} />
          </div>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
            Create Banner
          </Button>
        </CardContent>
      </Card>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Schedule Banner</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>Banner ID</Label>
              <Input value={scheduleId} onChange={(e) => setScheduleId(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Starts At (RFC3339)</Label>
              <Input value={scheduleStart} onChange={(e) => setScheduleStart(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Ends At (RFC3339)</Label>
              <Input value={scheduleEnd} onChange={(e) => setScheduleEnd(e.target.value)} />
            </div>
            <Button onClick={() => scheduleMutation.mutate()} disabled={scheduleMutation.isPending}>
              Schedule
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Set Position</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>Banner ID</Label>
              <Input value={positionId} onChange={(e) => setPositionId(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Position</Label>
              <Select value={position} onValueChange={setPosition}>
                <SelectTrigger>
                  <SelectValue placeholder="Select position" />
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
            <Button onClick={() => positionMutation.mutate()} disabled={positionMutation.isPending}>
              Set Position
            </Button>
          </CardContent>
        </Card>
      </div>

      <DataTable
        loading={isLoading}
        data={banners as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: "Title" },
          { key: "position", header: "Position" },
          {
            key: "is_active",
            header: "Active",
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? "Yes" : "No"}
              </Badge>
            ),
          },
        ]}
      />
    </div>
  );
}
