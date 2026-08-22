"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";

export default function TasksPage() {
  const [title, setTitle] = useState("");
  const [dueAt, setDueAt] = useState("");
  const [entityType, setEntityType] = useState("contact");
  const [entityId, setEntityId] = useState("");
  const [notes, setNotes] = useState("");
  const [success, setSuccess] = useState(false);

  const createMutation = useMutation({
    mutationFn: async () => {
      await api.post("/crm/tasks/reminders", {
        title,
        due_at: dueAt,
        entity_type: entityType,
        entity_id: entityId,
        notes,
      });
    },
    onSuccess: () => {
      setSuccess(true);
      setTitle("");
      setDueAt("");
      setEntityId("");
      setNotes("");
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Tasks</h1>
        <p className="text-muted-foreground">Create task reminders</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create Task Reminder</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Title</Label>
              <Input value={title} onChange={(e) => setTitle(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Due At (RFC3339)</Label>
              <Input value={dueAt} onChange={(e) => setDueAt(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Entity Type</Label>
              <Input value={entityType} onChange={(e) => setEntityType(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>Entity ID</Label>
              <Input value={entityId} onChange={(e) => setEntityId(e.target.value)} />
            </div>
          </div>
          <div className="space-y-2">
            <Label>Notes</Label>
            <Textarea value={notes} onChange={(e) => setNotes(e.target.value)} />
          </div>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
            Create Task
          </Button>
          {success && (
            <p className="text-sm text-emerald-600">Task reminder created successfully.</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
