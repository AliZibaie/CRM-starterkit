"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import { t } from "@/lib/i18n";

export interface Column<T> {
  key: string;
  header: string;
  render?: (row: T) => React.ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyField: keyof T | ((row: T) => string);
  emptyMessage?: string;
  loading?: boolean;
}

export function DataTable<T extends Record<string, unknown>>({
  columns,
  data,
  keyField,
  emptyMessage,
  loading,
}: DataTableProps<T>) {
  const getKey = (row: T) =>
    typeof keyField === "function" ? keyField(row) : String(row[keyField]);

  if (loading) {
    return (
      <div className="rounded-md border border-border p-8 text-center text-muted-foreground">
        {t("common.loading")}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-2xl border border-border/70 bg-card shadow-sm">
      <table className="w-full text-sm">
        <thead className="border-b border-border bg-muted/50">
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className={cn(
                  "px-4 py-3 text-start font-medium text-muted-foreground",
                  col.className
                )}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="px-4 py-8 text-center text-muted-foreground"
              >
                {emptyMessage ?? t("common.noData")}
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr
                key={getKey(row)}
                className="border-b border-border last:border-0 hover:bg-muted/30"
              >
                {columns.map((col) => (
                  <td key={col.key} className={cn("px-4 py-3 text-start", col.className)}>
                    {col.render
                      ? col.render(row)
                      : String(row[col.key] ?? "")}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
