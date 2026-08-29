import Link from "next/link";
import { t } from "@/lib/i18n";

export function SiteHeader() {
  return (
    <header className="border-b border-border bg-card">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-4 px-4">
        <div className="flex items-center gap-6">
          <Link href="/" className="text-lg font-bold text-primary">
            {t("brand")}
          </Link>
          <nav className="hidden items-center gap-4 text-sm sm:flex">
            <Link href="/categories/all" className="hover:text-primary">
              {t("nav.categories")}
            </Link>
            <Link href="/account/orders" className="hover:text-primary">
              {t("nav.orders")}
            </Link>
            <Link href="/account/tickets" className="hover:text-primary">
              {t("nav.support")}
            </Link>
          </nav>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <Link
            href="/categories/all"
            className="rounded-md px-2 py-1.5 hover:bg-accent"
            aria-label={t("nav.categories")}
          >
            🔍
          </Link>
          <Link
            href="/cart"
            className="rounded-md px-2 py-1.5 hover:bg-accent"
            aria-label={t("nav.cart")}
          >
            🛒
          </Link>
          <Link
            href="/account/login"
            className="flex items-center gap-1 rounded-md px-2 py-1.5 hover:bg-accent"
          >
            <span>👤</span>
            <span className="hidden sm:inline">{t("nav.account")}</span>
          </Link>
        </div>
      </div>
    </header>
  );
}
