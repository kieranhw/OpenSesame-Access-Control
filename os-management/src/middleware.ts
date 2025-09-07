import { NextRequest, NextResponse } from "next/server";
import { AppRoute } from "@/lib/app-routes";
import { ApiRoute, HUB_BASE_URI } from "@/lib/api/api";
import { SessionResponse } from "@/lib/api/session";

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  const currentPath = (route: AppRoute): boolean => {
    return pathname.includes(route);
  };

  if (isInternalPath(req)) return NextResponse.next();

  const cookies: string = req.headers.get("cookie") ?? "";
  let session: SessionResponse;
  let loginErrorMsg: string = "";
  const url = req.nextUrl.clone();

  try {
    const res = await fetch(HUB_BASE_URI + ApiRoute.ADMIN_SESSION, {
      method: "GET",
      headers: { cookie: cookies },
    });

    session = (await res.json()) as SessionResponse;
  } catch {
    loginErrorMsg = "Unable to reach OpenSesame, please ensure the service is running.";
    session = {
      /*
        If we can't communicate with the hub, we don't know if configuration is complete,
        so we assume that it is to prevent confusing the user with a redirect to /setup.
      */
      configured: true,
      authenticated: false,
    } as SessionResponse;
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
      if (loginErrorMsg) url.searchParams.set("error", loginErrorMsg);
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
