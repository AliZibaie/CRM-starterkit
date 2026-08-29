import { t } from "@/lib/i18n";

export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-border bg-card py-6">
      <div className="mx-auto max-w-6xl px-4 text-center text-sm text-muted-foreground">
        © {new Date().getFullYear()} {t("brand")} — {t("footer.rights")}
      </div>
    </footer>
  );
}
