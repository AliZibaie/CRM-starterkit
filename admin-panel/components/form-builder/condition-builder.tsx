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
  buildError?: string;
  building?: boolean;
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

function isGroupNode(node: ConditionNode) {
  return (node.children?.length ?? 0) > 0 || Boolean(node.logic);
}

function LeafEditor({
  columns,
  node,
  onChange,
}: {
  columns: ColumnMeta[];
  node: ConditionNode;
  onChange: (node: ConditionNode) => void;
}) {
  return (
    <div className="grid gap-3 md:grid-cols-3">
      <div className="space-y-2">
        <Label>{t("common.field")}</Label>
        <Select
          value={node.field || ""}
          onValueChange={(v) => onChange({ field: v, operator: node.operator || "eq", value: node.value })}
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
          value={node.operator || "eq"}
          onValueChange={(v) => onChange({ ...node, operator: v })}
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
          value={String(node.value ?? "")}
          disabled={node.operator === "is_null" || node.operator === "is_not_null"}
          onChange={(e) => onChange({ ...node, value: e.target.value })}
        />
      </div>
    </div>
  );
}

function GroupEditor({
  columns,
  node,
  onChange,
  depth = 0,
}: {
  columns: ColumnMeta[];
  node: ConditionNode;
  onChange: (node: ConditionNode) => void;
  depth?: number;
}) {
  const children = node.children ?? [];

  const updateChild = (index: number, child: ConditionNode) => {
    const next = [...children];
    next[index] = child;
    onChange({ logic: node.logic || "and", children: next });
  };

  const removeChild = (index: number) => {
    onChange({ logic: node.logic || "and", children: children.filter((_, i) => i !== index) });
  };

  return (
    <div className="space-y-3">
      <div className="space-y-2">
        <Label>{t("common.logic")}</Label>
        <Select
          value={node.logic || "and"}
          onValueChange={(v) => onChange({ logic: v, children })}
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
      {children.map((child, index) => (
        <div key={index} className="space-y-2 rounded-md border border-border p-3">
          {isGroupNode(child) && child.children ? (
            <GroupEditor
              columns={columns}
              node={child}
              onChange={(updated) => updateChild(index, updated)}
              depth={depth + 1}
            />
          ) : (
            <LeafEditor
              columns={columns}
              node={child}
              onChange={(updated) => updateChild(index, updated)}
            />
          )}
          <div className="flex gap-2">
            <Button type="button" variant="ghost" size="sm" onClick={() => removeChild(index)}>
              {t("conditionBuilder.remove")}
            </Button>
          </div>
        </div>
      ))}
      <div className="flex flex-wrap gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() =>
            onChange({
              logic: node.logic || "and",
              children: [...children, { field: "", operator: "eq", value: "" }],
            })
          }
        >
          {t("conditionBuilder.addRule")}
        </Button>
        {depth < 2 && (
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              onChange({
                logic: node.logic || "and",
                children: [
                  ...children,
                  { logic: "and", children: [{ field: "", operator: "eq", value: "" }] },
                ],
              })
            }
          >
            {t("conditionBuilder.addNestedGroup")}
          </Button>
        )}
      </div>
    </div>
  );
}

export function ConditionBuilder({
  columns,
  value,
  onChange,
  onBuild,
  buildResult,
  buildError,
  building,
}: ConditionBuilderProps) {
  const [mode, setMode] = useState<"leaf" | "group">(
    value.children?.length ? "group" : "leaf"
  );

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
            onClick={() => {
              setMode("leaf");
              onChange({
                field: value.field || "",
                operator: value.operator || "eq",
                value: value.value ?? "",
              });
            }}
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
                onChange({
                  logic: "and",
                  children: [{ field: "", operator: "eq", value: "" }],
                });
              }
            }}
          >
            {t("conditionBuilder.group")}
          </Button>
        </div>

        {mode === "leaf" ? (
          <LeafEditor
            columns={columns}
            node={value}
            onChange={(node) => onChange({ ...node, children: undefined, logic: undefined })}
          />
        ) : (
          <GroupEditor columns={columns} node={value} onChange={onChange} />
        )}

        {onBuild && (
          <div className="space-y-2">
            <Button type="button" onClick={onBuild} disabled={building}>
              {building ? t("common.saving") : t("conditionBuilder.buildCondition")}
            </Button>
            {buildError && (
              <p className="text-sm text-destructive">{buildError}</p>
            )}
            {buildResult && (
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs ltr-field">
                {buildResult}
              </pre>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
