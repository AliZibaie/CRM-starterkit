"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import { t } from "@/lib/i18n";

interface SimpleBarChartProps {
  data: { name: string; value: number }[];
  color?: string;
  height?: number;
}

export function SimpleBarChart({
  data,
  color = "#2563eb",
  height = 280,
}: SimpleBarChartProps) {
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
        <BarChart data={data} margin={{ top: 8, right: 8, left: 8, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
          <XAxis dataKey="name" tick={{ fontSize: 12 }} reversed />
          <YAxis tick={{ fontSize: 12 }} orientation="right" />
          <Tooltip />
          <Legend verticalAlign="top" align="right" />
          <Bar dataKey="value" fill={color} radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
