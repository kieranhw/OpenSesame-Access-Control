import { NextRequest, NextResponse } from "next/server";
import { ApiRoute, HUB_BASE_URL } from "./lib/api/api";
import { AppRoute } from "./lib/routes";
import { AuthResponse } from "./app/types/auth";

// TODO: This needs some careful thought for the potential flow cases
// then some refactoring - will come back to it when decided how this
// needs to work.

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (isInternalPath(req)) return NextResponse.next();

  console.log("getting cookies");
  const cookies: string = req.headers.get("cookie") ?? "";
  let session: AuthResponse;
  let errorMsg: string | undefined;
  try {
    console.log("About to fetch");

    const response = await fetch(HUB_BASE_URL + ApiRoute.SESSION, {
      method: "GET",
      headers: {
        cookie: cookies,
      },
    });

    session = (await response.json()) as AuthResponse;
    console.log("Response: ", session);

    if (!response.ok) {
      errorMsg = "Error validating session, please log in.";
    }
  } catch {
    // TODO: parse error codes, but we should likely say here we're unable to get to the hub
    errorMsg = "Unable to reach the hub.";
    console.error("Unable to reach server");
    session = {
      configured: false,
      authenticated: false,
    } as AuthResponse;
  }

  if (!session.configured) {
    console.log("Session not configured");
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
      url.searchParams.set(
        "error",
        errorMsg ?? "Session expired, please log in.",
      );
      return NextResponse.redirect(url);
    }

    return NextResponse.next();
  }

  // if (pathname.includes(AppRoute.LOGIN)) {
  //   // If we're on setup, check if we're not configured
  // }

  // if (pathname.includes(AppRoute.Setup)) {
  //   checkSystemConfigured(req);
  // }

  // if (pathname.includes(AppRoute.SETUP)) {
  //   // if we're on setup and configured direct to login
  //   try {
  //     const cfgRes = await getConfig();
  //     if (cfgRes.data) {
  //       const config = cfgRes.data;
  //       if (config.configured === true) {
  //         const url = req.nextUrl.clone();
  //         url.pathname = "/login";
  //         return NextResponse.redirect(url);
  //       } else {
  //         return NextResponse.next();
  //       }
  //     }
  //   } catch (e) {
  //     console.error("Could not fetch config:", e);
  //     // To prevent continuing to potentially problematic logic when hub is unreachable,
  //     // proceed to next (allow access to /setup if config can't be fetched)
  //     return NextResponse.next();
  //   }
  // }

  // // check if we need to be configured
  // if (!sessionCookie) {
  //   try {
  //     const cfgRes = await getConfig();
  //     if (cfgRes.data) {
  //       const config = cfgRes.data;
  //       if (config.configured === false) {
  //         const url = req.nextUrl.clone();
  //         console.log("Redirecting to setup");
  //         url.pathname = "/setup";
  //         return NextResponse.redirect(url);
  //       }
  //     }
  //   } catch (e) {
  //     console.error("Could not fetch config:", e);
  //     // Proceed without redirecting to /setup if config can't be fetched
  //     // (assume configured or handle later in validation)
  //   }
  // }

  // // if we are configured but no session then login
  // if (!sessionCookie && !pathname.startsWith("/login")) {
  //   const url = req.nextUrl.clone();
  //   url.pathname = "/login";
  //   return NextResponse.redirect(url);
  // }

  // // validate the session token with the hub
  // const rawCookieHeader: string = req.headers.get("cookie") ?? "";
  // try {
  //   const apiRes = await fetch(ApiRoute.SESSION, {
  //     method: "GET",
  //     headers: {
  //       // send all cookies
  //       cookie: rawCookieHeader,
  //     },
  //   });

  //   console.log("API res " + apiRes.status);

  //   const isValidSession = apiRes.status === 200;

  //   if (isValidSession && pathname.startsWith("/login")) {
  //     // we're logged in so redirect to home
  //     const url = req.nextUrl.clone();
  //     url.pathname = "/";
  //     return NextResponse.redirect(url);
  //   }

  //   if (!isValidSession && !pathname.startsWith("/login")) {
  //     // if we don't have a valid session and we're not on login, redirect to login
  //     const url = req.nextUrl.clone();
  //     url.pathname = "/login";
  //     return NextResponse.redirect(url);
  //   }
  // } catch {
  //   const url = req.nextUrl.clone();
  //   const errorMsg = "Unable to validate session, please login.";

  //   if (pathname.startsWith("/login")) {
  //     if (url.searchParams.has("error")) {
  //       // Already on /login with error param, proceed to render the page
  //       return NextResponse.next();
  //     } else {
  //       // Add the error param and redirect to make it visible in the URL
  //       url.searchParams.set("error", errorMsg);
  //       return NextResponse.redirect(url);
  //     }
  //   } else {
  //     // Not on /login, redirect to /login with error param
  //     // Note: This preserves any existing query params from the original URL
  //     url.pathname = "/login";
  //     url.searchParams.set("error", errorMsg);
  //     return NextResponse.redirect(url);
  //   }
  // }

  return NextResponse.next();
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

// function validateSessionCookie(req: NextRequest): boolean {}

export const config = {
  matcher: ["/:path*"],
};
