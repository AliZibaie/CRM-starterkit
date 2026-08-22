import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const token = request.cookies.get("access_token")?.value;
  const path = request.nextUrl.pathname;
  if (path.startsWith("/account") && !path.startsWith("/account/login")) {
    if (!token) {
      return NextResponse.redirect(new URL("/account/login", request.url));
    }
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/account/:path*"],
};
