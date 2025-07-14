import { NextRequest, NextResponse } from "next/server";
import { VALIDATE_URL } from "./lib/constants";
import { getConfig } from "./lib/api";

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
    }
  }

  const session = req.cookies.get("os_session")?.value;

  // check if we need to be configured
  if (!session) {
    try {
      const cfgRes = await getConfig();
      if (cfgRes.data) {
        const config = cfgRes.data;
        console.log("config: " + config.configured);
        if (config.configured === false) {
          const url = req.nextUrl.clone();
          console.log("Redirecting to setup");
          url.pathname = "/setup";
          return NextResponse.redirect(url);
        }
      }
    } catch (e) {
      console.error("Could not fetch config:", e);
    }
  }

  // if we are configured but no session then login
  if (!session && !pathname.startsWith("/login")) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // validate the session token with the hub
  const rawCookieHeader = req.headers.get("cookie") ?? "";
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

  return NextResponse.next();
}

// Run on all paths
export const config = {
  matcher: ["/:path*"],
};
