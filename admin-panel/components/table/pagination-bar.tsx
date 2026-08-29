"use client";

import { Button } from "@/components/ui/button";
import { t } from "@/lib/i18n";

interface PaginationBarProps {
  page: number;
  perPage: number;
  total: number;
  onPageChange: (page: number) => void;
}

export function PaginationBar({
  page,
  perPage,
  total,
  onPageChange,
}: PaginationBarProps) {
  const totalPages = Math.max(1, Math.ceil(total / perPage));
  if (totalPages <= 1) return null;

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 pt-4">
      <p className="text-sm text-muted-foreground">
        {t("common.pageOf", {
          page: String(page),
          total: String(totalPages),
        })}
      </p>
      <div className="flex gap-2">
        <Button
          size="sm"
          variant="outline"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          {t("common.previous")}
        </Button>
        <Button
          size="sm"
          variant="outline"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          {t("common.next")}
        </Button>
      </div>
    </div>
  );
}
