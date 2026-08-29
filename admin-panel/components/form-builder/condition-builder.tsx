"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { t } from "@/lib/i18n";

export interface ConditionNode {
  logic?: string;
  field?: string;
  operator?: string;
  value?: unknown;
  children?: ConditionNode[];
}

interface ColumnMeta {
  name: string;
  type: string;
}

interface ConditionBuilderProps {
  columns: ColumnMeta[];
  value: ConditionNode;
  onChange: (node: ConditionNode) => void;
  onBuild?: () => void;
  buildResult?: string;
}

const OPERATORS = [
  { value: "eq", labelKey: "conditionBuilder.eq" },
  { value: "neq", labelKey: "conditionBuilder.neq" },
  { value: "gt", labelKey: "conditionBuilder.gt" },
  { value: "gte", labelKey: "conditionBuilder.gte" },
  { value: "lt", labelKey: "conditionBuilder.lt" },
  { value: "lte", labelKey: "conditionBuilder.lte" },
  { value: "contains", labelKey: "conditionBuilder.contains" },
  { value: "is_null", labelKey: "conditionBuilder.is_null" },
  { value: "is_not_null", labelKey: "conditionBuilder.is_not_null" },
];

export function ConditionBuilder({
  columns,
  value,
  onChange,
  onBuild,
  buildResult,
}: ConditionBuilderProps) {
  const [mode, setMode] = useState<"leaf" | "group">(
    value.children?.length ? "group" : "leaf"
  );

  const updateLeaf = (patch: Partial<ConditionNode>) => {
    onChange({ ...value, ...patch, children: undefined });
  };

  const updateGroup = (logic: string, children: ConditionNode[]) => {
    onChange({ logic, children });
  };

  const addChild = () => {
    const children = [...(value.children || []), { field: "", operator: "eq", value: "" }];
    updateGroup(value.logic || "and", children);
  };

  const updateChild = (index: number, child: ConditionNode) => {
    const children = [...(value.children || [])];
    children[index] = child;
    updateGroup(value.logic || "and", children);
  };

  const removeChild = (index: number) => {
    const children = (value.children || []).filter((_, i) => i !== index);
    updateGroup(value.logic || "and", children);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{t("conditionBuilder.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Button
            type="button"
            size="sm"
            variant={mode === "leaf" ? "default" : "outline"}
            onClick={() => setMode("leaf")}
          >
            {t("conditionBuilder.singleRule")}
          </Button>
          <Button
            type="button"
            size="sm"
            variant={mode === "group" ? "default" : "outline"}
            onClick={() => {
              setMode("group");
              if (!value.children?.length) {
                updateGroup("and", [{ field: "", operator: "eq", value: "" }]);
              }
            }}
          >
            {t("conditionBuilder.group")}
          </Button>
        </div>

        {mode === "leaf" ? (
          <div className="grid gap-3 md:grid-cols-3">
            <div className="space-y-2">
              <Label>{t("common.field")}</Label>
              <Select
                value={value.field || ""}
                onValueChange={(v) => updateLeaf({ field: v })}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t("conditionBuilder.column")} />
                </SelectTrigger>
                <SelectContent>
                  {columns.map((col) => (
                    <SelectItem key={col.name} value={col.name}>
                      {col.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.operator")}</Label>
              <Select
                value={value.operator || "eq"}
                onValueChange={(v) => updateLeaf({ operator: v })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {OPERATORS.map((op) => (
                    <SelectItem key={op.value} value={op.value}>
                      {t(op.labelKey)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>{t("common.value")}</Label>
              <Input
                value={String(value.value ?? "")}
                onChange={(e) => updateLeaf({ value: e.target.value })}
              />
            </div>
          </div>
        ) : (
          <div className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.logic")}</Label>
              <Select
                value={value.logic || "and"}
                onValueChange={(v) =>
                  updateGroup(v, value.children || [])
                }
              >
                <SelectTrigger className="w-32">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="and">{t("conditionBuilder.and")}</SelectItem>
                  <SelectItem value="or">{t("conditionBuilder.or")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            {(value.children || []).map((child, index) => (
              <div key={index} className="flex gap-2 rounded-md border border-border p-3">
                <div className="grid flex-1 gap-2 md:grid-cols-3">
                  <Select
                    value={child.field || ""}
                    onValueChange={(v) =>
                      updateChild(index, { ...child, field: v })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("common.field")} />
                    </SelectTrigger>
                    <SelectContent>
                      {columns.map((col) => (
                        <SelectItem key={col.name} value={col.name}>
                          {col.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Select
                    value={child.operator || "eq"}
                    onValueChange={(v) =>
                      updateChild(index, { ...child, operator: v })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {OPERATORS.map((op) => (
                        <SelectItem key={op.value} value={op.value}>
                          {t(op.labelKey)}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Input
                    value={String(child.value ?? "")}
                    onChange={(e) =>
                      updateChild(index, { ...child, value: e.target.value })
                    }
                  />
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => removeChild(index)}
                >
                  {t("conditionBuilder.remove")}
                </Button>
              </div>
            ))}
            <Button type="button" variant="outline" size="sm" onClick={addChild}>
              {t("conditionBuilder.addRule")}
            </Button>
          </div>
        )}

        {onBuild && (
          <div className="space-y-2">
            <Button type="button" onClick={onBuild}>
              {t("conditionBuilder.buildCondition")}
            </Button>
            {buildResult && (
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs">
                {buildResult}
              </pre>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
