import { NextRequest, NextResponse } from "next/server";

const PUBLIC_PATHS = ["/login", "/api/login"];

const VALIDATE_URL = "http://localhost:11072/management/validate_session";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (
    pathname.startsWith("/_next/") ||
    pathname.startsWith("/static/") ||
    pathname === "/favicon.ico" ||
    PUBLIC_PATHS.includes(pathname) ||
    pathname.includes(".")
  ) {
    return NextResponse.next();
  }

  // get the session token
  const token = req.cookies.get("session_token")?.value;
  if (!token) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // validate the session token with the hub
  const apiRes = await fetch(VALIDATE_URL, {
    method: "GET",
    headers: {
      // forward the same cookies so the hub sees them
      cookie: req.headers.get("cookie") ?? "",
    },
  });

  if (apiRes.status !== 200) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Run on all paths
export const config = {
  matcher: ["/:path*"],
};
