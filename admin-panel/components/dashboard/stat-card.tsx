import type { LucideIcon } from "lucide-react";
import { TrendingDown, TrendingUp } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface StatCardProps {
  title: string;
  value: string;
  icon: LucideIcon;
  tone?: "primary" | "success" | "warning" | "info";
  changePercent?: number;
}

const toneStyles = {
  primary: "from-indigo-500/15 to-indigo-500/5 text-indigo-600 dark:text-indigo-300",
  success: "from-emerald-500/15 to-emerald-500/5 text-emerald-600 dark:text-emerald-300",
  warning: "from-amber-500/15 to-amber-500/5 text-amber-600 dark:text-amber-300",
  info: "from-sky-500/15 to-sky-500/5 text-sky-600 dark:text-sky-300",
};

export function StatCard({
  title,
  value,
  icon: Icon,
  tone = "primary",
  changePercent,
}: StatCardProps) {
  const positive = changePercent !== undefined && changePercent >= 0;

  return (
    <Card className="overflow-hidden">
      <CardContent className="p-0">
        <div className="flex items-stretch">
          <div className={cn("flex w-1.5 shrink-0 bg-gradient-to-b", toneStyles[tone])} />
          <div className="flex flex-1 items-start justify-between gap-4 p-5">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">{title}</p>
              <p className="text-3xl font-bold tracking-tight">{value}</p>
              {changePercent !== undefined && (
                <div
                  className={cn(
                    "flex items-center gap-1 text-xs font-medium",
                    positive ? "text-emerald-600" : "text-red-500"
                  )}
                >
                  {positive ? (
                    <TrendingUp className="h-3.5 w-3.5" />
                  ) : (
                    <TrendingDown className="h-3.5 w-3.5" />
                  )}
                  <span>{Math.abs(changePercent).toLocaleString("fa-IR")}٪</span>
                </div>
              )}
            </div>
            <div
              className={cn(
                "rounded-2xl bg-gradient-to-br p-3",
                toneStyles[tone]
              )}
            >
              <Icon className="h-5 w-5" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
