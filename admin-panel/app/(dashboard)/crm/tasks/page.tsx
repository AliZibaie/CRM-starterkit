"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { t } from "@/lib/i18n";
import { datetimeLocalToRFC3339 } from "@/lib/datetime-input";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function TasksPage() {
  const [title, setTitle] = useState("");
  const [dueAtLocal, setDueAtLocal] = useState("");
  const [entityType, setEntityType] = useState("contact");
  const [entityId, setEntityId] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const createMutation = useMutation({
    mutationFn: async () => {
      const dueAt = datetimeLocalToRFC3339(dueAtLocal);
      await api.post("/crm/tasks/reminders", {
        title,
        due_at: dueAt,
        reminder_at: dueAt,
        entity_type: entityType,
        entity_id: entityId,
        notes,
      });
    },
    onSuccess: () => {
      setFeedback({ type: "success", text: t("common.taskCreatedSuccess") });
      setTitle("");
      setDueAtLocal("");
      setEntityId("");
      setNotes("");
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.tasks.title")}</h1>
        <p className="text-muted-foreground">{t("pages.tasks.subtitle")}</p>
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
          <CardTitle>{t("pages.tasks.createReminder")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>{t("common.title")}</Label>
              <Input value={title} onChange={(e) => setTitle(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("pages.tasks.dueAtLocal")}</Label>
              <Input type="datetime-local" value={dueAtLocal} onChange={(e) => setDueAtLocal(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.entityType")}</Label>
              <Select value={entityType} onValueChange={setEntityType}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="contact">{t("nav.contacts")}</SelectItem>
                  <SelectItem value="lead">{t("nav.leads")}</SelectItem>
                  <SelectItem value="deal">{t("pages.pipeline.pipeline")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.entityId")}</Label>
              <Input value={entityId} onChange={(e) => setEntityId(e.target.value)} />
            </div>
          </div>
          <div className="space-y-2">
            <Label>{t("common.notes")}</Label>
            <Textarea value={notes} onChange={(e) => setNotes(e.target.value)} />
          </div>
          <Button onClick={() => createMutation.mutate()} disabled={createMutation.isPending || !title.trim()}>
            {t("pages.tasks.createTask")}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
