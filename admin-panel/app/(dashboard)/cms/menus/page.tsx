"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function CmsMenusPage() {
  const [name, setName] = useState("");
  const [itemsJson, setItemsJson] = useState("[]");
  const [success, setSuccess] = useState(false);

  const upsertMutation = useMutation({
    mutationFn: async () => {
      await api.post("/cms/menus", { name, items_json: itemsJson });
    },
    onSuccess: () => {
      setSuccess(true);
      setName("");
      setItemsJson("[]");
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">CMS Menus</h1>
        <p className="text-muted-foreground">Upsert navigation menus</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Upsert Menu</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Menu Name</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Items JSON</Label>
            <Input
              value={itemsJson}
              onChange={(e) => setItemsJson(e.target.value)}
              placeholder='[{"label":"Home","url":"/"}]'
            />
          </div>
          <Button onClick={() => upsertMutation.mutate()} disabled={upsertMutation.isPending}>
            Save Menu
          </Button>
          {success && (
            <p className="text-sm text-emerald-600">Menu saved successfully.</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
