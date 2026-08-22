"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";

export interface FormField {
  name: string;
  label: string;
  type: "text" | "email" | "password" | "number" | "textarea" | "select" | "checkbox";
  placeholder?: string;
  required?: boolean;
  options?: { label: string; value: string }[];
}

interface DynamicFormProps {
  fields: FormField[];
  values: Record<string, unknown>;
  onChange: (name: string, value: unknown) => void;
  onSubmit: () => void;
  submitLabel?: string;
  loading?: boolean;
}

export function DynamicForm({
  fields,
  values,
  onChange,
  onSubmit,
  submitLabel = "Submit",
  loading,
}: DynamicFormProps) {
  return (
    <form
      className="space-y-4"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
    >
      {fields.map((field) => (
        <div key={field.name} className="space-y-2">
          <Label htmlFor={field.name}>{field.label}</Label>
          {field.type === "textarea" ? (
            <Textarea
              id={field.name}
              value={String(values[field.name] ?? "")}
              onChange={(e) => onChange(field.name, e.target.value)}
              placeholder={field.placeholder}
              required={field.required}
            />
          ) : field.type === "select" ? (
            <Select
              value={String(values[field.name] ?? "")}
              onValueChange={(v) => onChange(field.name, v)}
            >
              <SelectTrigger id={field.name}>
                <SelectValue placeholder={field.placeholder || "Select"} />
              </SelectTrigger>
              <SelectContent>
                {field.options?.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          ) : field.type === "checkbox" ? (
            <input
              id={field.name}
              type="checkbox"
              checked={Boolean(values[field.name])}
              onChange={(e) => onChange(field.name, e.target.checked)}
              className="h-4 w-4"
            />
          ) : (
            <Input
              id={field.name}
              type={field.type}
              value={String(values[field.name] ?? "")}
              onChange={(e) =>
                onChange(
                  field.name,
                  field.type === "number" ? Number(e.target.value) : e.target.value
                )
              }
              placeholder={field.placeholder}
              required={field.required}
            />
          )}
        </div>
      ))}
      <Button type="submit" disabled={loading}>
        {loading ? "Saving..." : submitLabel}
      </Button>
    </form>
  );
}
