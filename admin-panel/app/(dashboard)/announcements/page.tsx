"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface Announcement {
  id: string;
  title: string;
  body: string;
  created_at: string;
}

const fields: FormField[] = [
  { name: "title", label: "Title", type: "text", required: true },
  { name: "body", label: "Body", type: "textarea", required: true },
];

export default function AnnouncementsPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [values, setValues] = useState<Record<string, unknown>>({
    title: "",
    body: "",
  });

  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.announcements,
    queryFn: async () => {
      const res = await api.get<Announcement[]>("/announcements");
      return res.data ?? [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/announcements", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.announcements });
      setShowForm(false);
      setValues({ title: "", body: "" });
    },
  });

  const broadcastMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.post(`/announcements/${id}/broadcast`);
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Announcements</h1>
          <p className="text-muted-foreground">Create and broadcast announcements</p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "New Announcement"}
        </Button>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>Create Announcement</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={fields}
              values={values}
              onChange={(n, v) => setValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() => createMutation.mutate(values)}
              loading={createMutation.isPending}
            />
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "title", header: "Title" },
          { key: "body", header: "Body" },
          { key: "created_at", header: "Created" },
          {
            key: "actions",
            header: "Actions",
            render: (row) => (
              <Button
                size="sm"
                onClick={() => broadcastMutation.mutate(row.id as string)}
                disabled={broadcastMutation.isPending}
              >
                Broadcast
              </Button>
            ),
          },
        ]}
      />
    </div>
  );
}
