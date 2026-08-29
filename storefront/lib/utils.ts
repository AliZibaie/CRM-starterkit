import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatPrice(amount: number, currency = "IRR") {
  const unit = currency === "IRR" ? "ریال" : currency;
  return `${new Intl.NumberFormat("fa-IR").format(amount)} ${unit}`;
}
