import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const publicPaths = ["/login"];

const protectedPrefixes = [
  "/dashboard",
  "/users",
  "/roles-permissions",
  "/settings",
  "/cron-jobs",
  "/announcements",
  "/banners",
  "/cms",
  "/media-library",
  "/support",
  "/catalog",
  "/orders",
  "/payments",
  "/discounts",
  "/crm",
  "/reviews",
  "/request-monitor",
  "/audit-log",
];

export function middleware(request: NextRequest) {
  const token = request.cookies.get("access_token")?.value;
  const { pathname } = request.nextUrl;
  const isPublic = publicPaths.some((p) => pathname.startsWith(p));
  const isProtected = protectedPrefixes.some(
    (p) => pathname === p || pathname.startsWith(`${p}/`)
  );

  if (isProtected && !token) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  if (isPublic && token) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
