import { NextRequest, NextResponse } from "next/server";
import { ApiRoute, HUB_BASE_URL } from "./lib/api/api";
import { AppRoute } from "./lib/routes";
import { AuthResponse } from "./app/types/auth";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (isInternalPath(req)) return NextResponse.next();

  const cookies: string = req.headers.get("cookie") ?? "";
  let session: AuthResponse;
  let errorMsg: string | undefined;
  try {
    const response = await fetch(HUB_BASE_URL + ApiRoute.SESSION, {
      method: "GET",
      headers: {
        cookie: cookies,
      },
    });

    session = (await response.json()) as AuthResponse;

    if (!response.ok) {
      errorMsg = "Please log in.";
    }
  } catch {
    // TODO: parse error codes, but we should likely say here we're unable to get to the hub
    errorMsg = "Unable to reach the hub, please try again later.";
    session = {
      configured: false,
      authenticated: false,
    } as AuthResponse;
  }

  if (!session.configured) {
    if (pathname.includes(AppRoute.SETUP)) {
      return NextResponse.next();
    }

    const url = req.nextUrl.clone();
    url.pathname = AppRoute.SETUP;
    return NextResponse.redirect(url);
  }

  if (session.authenticated) {
    if (
      pathname.includes(AppRoute.SETUP) ||
      pathname.includes(AppRoute.LOGIN)
    ) {
      const url = req.nextUrl.clone();
      url.pathname = AppRoute.HOME;
      return NextResponse.redirect(url);
    }
  } else {
    if (!pathname.includes(AppRoute.LOGIN)) {
      const url = req.nextUrl.clone();
      url.pathname = AppRoute.LOGIN;
      url.searchParams.set("error", errorMsg ?? "Please log in.");
      return NextResponse.redirect(url);
    }

    return NextResponse.next();
  }
}

function isInternalPath(req: NextRequest): boolean {
  const { pathname } = req.nextUrl;

  if (
    // skip internal, static, and paths with file extensions
    pathname.startsWith("/_next/") ||
    pathname.startsWith("/static/") ||
    pathname === "/favicon.ico" ||
    pathname.includes(".")
  ) {
    return true;
  }

  return false;
}

export const config = {
  matcher: ["/:path*"],
};
