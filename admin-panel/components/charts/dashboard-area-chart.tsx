"use client";

import {
  Area,
  AreaChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { t } from "@/lib/i18n";

interface DashboardAreaChartProps {
  data: { name: string; value: number }[];
  color?: string;
  height?: number;
}

export function DashboardAreaChart({
  data,
  color = "#6366f1",
  height = 300,
}: DashboardAreaChartProps) {
  if (data.length === 0) {
    return (
      <div
        className="flex items-center justify-center text-sm text-muted-foreground"
        style={{ height }}
      >
        {t("common.noChartData")}
      </div>
    );
  }

  return (
    <div dir="rtl" style={{ width: "100%", height }}>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 8, right: 8, left: 8, bottom: 0 }}>
          <defs>
            <linearGradient id="areaFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor={color} stopOpacity={0.35} />
              <stop offset="95%" stopColor={color} stopOpacity={0.02} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" className="stroke-border/70" />
          <XAxis dataKey="name" tick={{ fontSize: 12 }} reversed />
          <YAxis tick={{ fontSize: 12 }} orientation="right" />
          <Tooltip />
          <Legend verticalAlign="top" align="right" />
          <Area
            type="monotone"
            dataKey="value"
            stroke={color}
            fill="url(#areaFill)"
            strokeWidth={2}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
