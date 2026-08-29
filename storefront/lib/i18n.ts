import fa from "@/messages/fa.json";

export function t(key: string, vars?: Record<string, string | number>): string {
  const parts = key.split(".");
  let val: unknown = fa;
  for (const part of parts) {
    if (val && typeof val === "object" && part in val) {
      val = (val as Record<string, unknown>)[part];
    } else {
      return key;
    }
  }
  if (typeof val !== "string") return key;
  if (!vars) return val;
  return Object.entries(vars).reduce(
    (text, [k, v]) => text.replace(`{${k}}`, String(v)),
    val
  );
}

export function statusLabel(status: string): string {
  const key = `status.${status}`;
  const label = t(key);
  return label === key ? status : label;
}
