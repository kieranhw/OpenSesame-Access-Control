import { NextRequest, NextResponse } from "next/server";
import { AppRoute } from "@/lib/app-routes";
import { ApiRoute, HUB_BASE_URI } from "@/lib/api/api";
import { AuthResponse } from "@/lib/api/auth";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  const currentPath = (route: AppRoute): boolean => {
    return pathname.includes(route);
  };

  if (isInternalPath(req)) return NextResponse.next();

  const cookies: string = req.headers.get("cookie") ?? "";
  let session: AuthResponse;
  let loginErrorMsg: string | undefined = "Session expired, please log in.";
  const url = req.nextUrl.clone();

  try {
    const res = await fetch(HUB_BASE_URI + ApiRoute.SESSION, {
      method: "GET",
      headers: { cookie: cookies },
    });

    session = (await res.json()) as AuthResponse;
  } catch {
    loginErrorMsg = "Unable to reach the hub, please try again later.";
    session = {
      // Set configured true here to prevent user access to /setup until we get a valid response
      configured: true,
      authenticated: false,
    } as AuthResponse;
  }

  if (!session.configured) {
    if (currentPath(AppRoute.SETUP)) {
      return NextResponse.next();
    } else {
      url.pathname = AppRoute.SETUP;
      return NextResponse.redirect(url);
    }
  }

  if (session.authenticated) {
    if (currentPath(AppRoute.SETUP) || currentPath(AppRoute.LOGIN)) {
      url.pathname = AppRoute.HOME;
      return NextResponse.redirect(url);
    }
  } else {
    if (!currentPath(AppRoute.LOGIN)) {
      url.pathname = AppRoute.LOGIN;
      url.searchParams.set("error", loginErrorMsg);
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
