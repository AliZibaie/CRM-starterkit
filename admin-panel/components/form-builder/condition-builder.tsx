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
  { value: "eq", label: "Equals" },
  { value: "neq", label: "Not equals" },
  { value: "gt", label: "Greater than" },
  { value: "gte", label: "Greater or equal" },
  { value: "lt", label: "Less than" },
  { value: "lte", label: "Less or equal" },
  { value: "contains", label: "Contains" },
  { value: "is_null", label: "Is null" },
  { value: "is_not_null", label: "Is not null" },
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
        <CardTitle className="text-base">Condition Builder</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Button
            type="button"
            size="sm"
            variant={mode === "leaf" ? "default" : "outline"}
            onClick={() => setMode("leaf")}
          >
            Single Rule
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
            Group
          </Button>
        </div>

        {mode === "leaf" ? (
          <div className="grid gap-3 md:grid-cols-3">
            <div className="space-y-2">
              <Label>Field</Label>
              <Select
                value={value.field || ""}
                onValueChange={(v) => updateLeaf({ field: v })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Column" />
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
              <Label>Operator</Label>
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
                      {op.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Value</Label>
              <Input
                value={String(value.value ?? "")}
                onChange={(e) => updateLeaf({ value: e.target.value })}
              />
            </div>
          </div>
        ) : (
          <div className="space-y-3">
            <div className="space-y-2">
              <Label>Logic</Label>
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
                  <SelectItem value="and">AND</SelectItem>
                  <SelectItem value="or">OR</SelectItem>
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
                      <SelectValue placeholder="Field" />
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
                          {op.label}
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
                  Remove
                </Button>
              </div>
            ))}
            <Button type="button" variant="outline" size="sm" onClick={addChild}>
              Add Rule
            </Button>
          </div>
        )}

        {onBuild && (
          <div className="space-y-2">
            <Button type="button" onClick={onBuild}>
              Build Condition
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
