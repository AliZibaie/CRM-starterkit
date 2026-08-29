import { ApiError } from "@/lib/api-client";
import { t } from "@/lib/i18n";

export function getApiErrorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    if (error.message === "network_error") {
      return t("common.networkError");
    }
    return error.message || t("common.error");
  }
  if (error instanceof Error) {
    return error.message;
  }
  return t("common.error");
}
