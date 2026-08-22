import type { Metadata } from "next";
import "./globals.css";
import { Providers } from "@/components/providers";
import { SiteHeader } from "@/components/layout/site-header";

export const metadata: Metadata = {
  title: "Store",
  description: "CRM StarterKit Storefront",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <SiteHeader />
          <main className="mx-auto min-h-screen max-w-6xl px-4 py-6">{children}</main>
        </Providers>
      </body>
    </html>
  );
}
