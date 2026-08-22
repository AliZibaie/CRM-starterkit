import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="border-b border-border bg-card">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
        <Link href="/" className="text-lg font-semibold">
          Store
        </Link>
        <nav className="flex items-center gap-4 text-sm">
          <Link href="/categories/all" className="hover:text-primary">
            Categories
          </Link>
          <Link href="/cart" className="hover:text-primary">
            Cart
          </Link>
          <Link href="/account/orders" className="hover:text-primary">
            Account
          </Link>
        </nav>
      </div>
    </header>
  );
}
