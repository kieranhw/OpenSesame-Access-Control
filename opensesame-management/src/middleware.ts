import { NextRequest, NextResponse } from "next/server";
import { VALIDATE_URL } from "./lib/constants";
import { getConfig } from "./lib/api";

// TODO: This needs some careful thought for the potential flow cases
// then some refactoring - will come back to it when decided how this
// needs to work.

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (
    pathname.startsWith("/_next/") ||
    pathname.startsWith("/static/") ||
    pathname === "/favicon.ico" ||
    pathname.includes(".")
  ) {
    // skip internal, static, and paths with file extensions
    return NextResponse.next();
  }

  if (pathname.includes("/setup")) {
    // if we're on setup and configured direct to login
    try {
      const cfgRes = await getConfig();
      if (cfgRes.data) {
        const config = cfgRes.data;
        if (config.configured === true) {
          const url = req.nextUrl.clone();
          url.pathname = "/login";
          return NextResponse.redirect(url);
        } else {
          return NextResponse.next();
        }
      }
    } catch (e) {
      console.error("Could not fetch config:", e);
      // To prevent continuing to potentially problematic logic when hub is unreachable,
      // proceed to next (allow access to /setup if config can't be fetched)
      return NextResponse.next();
    }
  }

  const sessionCookie: string | undefined =
    req.cookies.get("os_session")?.value;

  // check if we need to be configured
  if (!sessionCookie) {
    try {
      const cfgRes = await getConfig();
      if (cfgRes.data) {
        const config = cfgRes.data;
        if (config.configured === false) {
          const url = req.nextUrl.clone();
          console.log("Redirecting to setup");
          url.pathname = "/setup";
          return NextResponse.redirect(url);
        }
      }
    } catch (e) {
      console.error("Could not fetch config:", e);
      // Proceed without redirecting to /setup if config can't be fetched
      // (assume configured or handle later in validation)
    }
  }

  // if we are configured but no session then login
  if (!sessionCookie && !pathname.startsWith("/login")) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // validate the session token with the hub
  const rawCookieHeader: string = req.headers.get("cookie") ?? "";
  try {
    const apiRes = await fetch(VALIDATE_URL, {
      method: "GET",
      headers: {
        // send all cookies
        cookie: rawCookieHeader,
      },
    });

    console.log("API res " + apiRes.status);

    const isValidSession = apiRes.status === 200;

    if (isValidSession && pathname.startsWith("/login")) {
      // we're logged in so redirect to home
      const url = req.nextUrl.clone();
      url.pathname = "/";
      return NextResponse.redirect(url);
    }

    if (!isValidSession && !pathname.startsWith("/login")) {
      // if we don't have a valid session and we're not on login, redirect to login
      const url = req.nextUrl.clone();
      url.pathname = "/login";
      return NextResponse.redirect(url);
    }
  } catch {
    const url = req.nextUrl.clone();
    const errorMsg = "Unable to validate session, please login.";

    if (pathname.startsWith("/login")) {
      if (url.searchParams.has("error")) {
        // Already on /login with error param, proceed to render the page
        return NextResponse.next();
      } else {
        // Add the error param and redirect to make it visible in the URL
        url.searchParams.set("error", errorMsg);
        return NextResponse.redirect(url);
      }
    } else {
      // Not on /login, redirect to /login with error param
      // Note: This preserves any existing query params from the original URL
      url.pathname = "/login";
      url.searchParams.set("error", errorMsg);
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}

// Run on all paths
export const config = {
  matcher: ["/:path*"],
};