export function formatJalali(
  value: string | Date,
  options?: Intl.DateTimeFormatOptions
): string {
  const date = typeof value === "string" ? new Date(value) : value;
  return new Intl.DateTimeFormat("fa-IR", {
    calendar: "persian",
    year: "numeric",
    month: "long",
    day: "numeric",
    ...options,
  }).format(date);
}

export function formatJalaliDateTime(value: string | Date): string {
  return formatJalali(value, {
    hour: "2-digit",
    minute: "2-digit",
  });
}
