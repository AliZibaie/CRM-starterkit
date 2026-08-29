"use client";

import { useMemo, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Pipeline {
  id: string;
  name: string;
  stages_json: string;
}

interface Contact {
  id: string;
  name: string;
}

interface Deal {
  id: string;
  title: string;
  contact_id: string;
  pipeline_id: string;
  stage: string;
  amount: number;
}

export default function PipelinePage() {
  const queryClient = useQueryClient();
  const [pipelineName, setPipelineName] = useState("");
  const [stagesJson, setStagesJson] = useState('["lead","qualified","proposal","won","lost"]');
  const [dealTitle, setDealTitle] = useState("");
  const [contactId, setContactId] = useState("");
  const [pipelineId, setPipelineId] = useState("");
  const [stage, setStage] = useState("");
  const [dragDealId, setDragDealId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(null);

  const { data: pipelines = [] } = useQuery({
    queryKey: queryKeys.pipelines,
    queryFn: async () => {
      const res = await api.get<Pipeline[]>("/crm/pipelines");
      return res.data ?? [];
    },
  });

  const { data: contacts = [] } = useQuery({
    queryKey: queryKeys.contacts,
    queryFn: async () => {
      const res = await api.get<Contact[]>("/crm/contacts");
      return res.data ?? [];
    },
  });

  const { data: deals = [] } = useQuery({
    queryKey: queryKeys.deals(pipelineId),
    queryFn: async () => {
      const res = await api.get<Deal[]>(`/crm/deals?pipeline_id=${pipelineId}`);
      return res.data ?? [];
    },
    enabled: !!pipelineId,
  });

  const selectedPipeline = pipelines.find((p) => p.id === pipelineId);
  const stages = useMemo(() => {
    if (!selectedPipeline) return [];
    try {
      return JSON.parse(selectedPipeline.stages_json) as string[];
    } catch {
      return [];
    }
  }, [selectedPipeline]);

  const dealsByStage = useMemo(() => {
    const map: Record<string, Deal[]> = {};
    stages.forEach((s) => {
      map[s] = [];
    });
    deals.forEach((d) => {
      if (map[d.stage]) map[d.stage].push(d);
      else map[d.stage] = [d];
    });
    return map;
  }, [deals, stages]);

  const upsertPipeline = useMutation({
    mutationFn: async () => {
      await api.post("/crm/pipelines", { name: pipelineName, stages_json: stagesJson });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.pipelines });
      setPipelineName("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const createDeal = useMutation({
    mutationFn: async () => {
      await api.post("/crm/deals", {
        title: dealTitle,
        contact_id: contactId,
        pipeline_id: pipelineId,
        stage,
        amount: 0,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.deals(pipelineId) });
      setDealTitle("");
      setContactId("");
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const moveDeal = useMutation({
    mutationFn: async ({ id, newStage }: { id: string; newStage: string }) => {
      await api.post(`/crm/deals/${id}/stage`, { stage: newStage });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.deals(pipelineId) });
    },
    onError: (err) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
  });

  const handleDrop = (targetStage: string) => {
    if (!dragDealId || !targetStage) return;
    const deal = deals.find((d) => d.id === dragDealId);
    if (!deal || deal.stage === targetStage) {
      setDragDealId(null);
      return;
    }
    moveDeal.mutate({ id: dragDealId, newStage: targetStage });
    setDragDealId(null);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("pages.pipeline.title")}</h1>
        <p className="text-muted-foreground">{t("pages.pipeline.subtitle")}</p>
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

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.pipeline.upsertPipeline")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.name")}</Label>
              <Input value={pipelineName} onChange={(e) => setPipelineName(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("pages.pipeline.stagesJson")}</Label>
              <Input value={stagesJson} onChange={(e) => setStagesJson(e.target.value)} />
            </div>
            <Button onClick={() => upsertPipeline.mutate()} disabled={upsertPipeline.isPending}>
              {t("pages.pipeline.savePipeline")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("pages.pipeline.createDeal")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Select
              value={pipelineId}
              onValueChange={(v) => {
                setPipelineId(v);
                setStage("");
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectPipeline")} />
              </SelectTrigger>
              <SelectContent>
                {pipelines.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Input placeholder={t("common.title")} value={dealTitle} onChange={(e) => setDealTitle(e.target.value)} />
            <Select value={contactId} onValueChange={setContactId}>
              <SelectTrigger>
                <SelectValue placeholder={t("pages.pipeline.contactId")} />
              </SelectTrigger>
              <SelectContent>
                {contacts.map((c) => (
                  <SelectItem key={c.id} value={c.id}>
                    {c.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={stage} onValueChange={setStage}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectStage")} />
              </SelectTrigger>
              <SelectContent>
                {stages.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              onClick={() => createDeal.mutate()}
              disabled={createDeal.isPending || !dealTitle || !pipelineId || !stage}
            >
              {t("pages.pipeline.createDeal")}
            </Button>
          </CardContent>
        </Card>
      </div>

      {pipelineId && stages.length > 0 && (
        <div className="flex gap-4 overflow-x-auto pb-4">
          {stages.map((s) => (
            <div
              key={s}
              className="min-w-[220px] flex-1"
              onDragOver={(e) => e.preventDefault()}
              onDrop={() => handleDrop(s)}
            >
              <Card className="h-full">
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm capitalize">{s}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  {(dealsByStage[s] ?? []).map((deal) => (
                    <div
                      key={deal.id}
                      draggable
                      onDragStart={() => setDragDealId(deal.id)}
                      onDragEnd={() => setDragDealId(null)}
                      className="cursor-grab rounded-md border border-border bg-background p-3 text-sm shadow-sm active:cursor-grabbing"
                    >
                      <p className="font-medium">{deal.title}</p>
                    </div>
                  ))}
                  {(dealsByStage[s] ?? []).length === 0 && (
                    <p className="text-xs text-muted-foreground">{t("pages.pipeline.dropHere")}</p>
                  )}
                </CardContent>
              </Card>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
