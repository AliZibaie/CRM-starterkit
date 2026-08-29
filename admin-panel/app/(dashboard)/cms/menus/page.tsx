"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { t } from "@/lib/i18n";
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
        <h1 className="text-2xl font-bold">{t("pages.cmsMenus.title")}</h1>
        <p className="text-muted-foreground">{t("pages.cmsMenus.subtitle")}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.cmsMenus.upsertMenu")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>{t("common.menuName")}</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>{t("common.itemsJson")}</Label>
            <Input
              value={itemsJson}
              onChange={(e) => setItemsJson(e.target.value)}
              placeholder={t("pages.cmsMenus.itemsPlaceholder")}
            />
          </div>
          <Button onClick={() => upsertMutation.mutate()} disabled={upsertMutation.isPending}>
            {t("pages.cmsMenus.save")}
          </Button>
          {success && (
            <p className="text-sm text-emerald-600">{t("common.menuSavedSuccess")}</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
