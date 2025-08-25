import { ApiResponse, ApiRoute, hubApiClient } from "./api";
import { ApiError } from "./api-error";

export interface SessionResponse {
  message?: string;
  authenticated: boolean;
  configured: boolean;
}

export interface LogoutResponse {
  success: boolean;
}

/**
 * Login with password
 */
async function login(password: string): Promise<ApiResponse<SessionResponse>> {
  try {
    const res = await hubApiClient.post<SessionResponse>(ApiRoute.SESSION, {
      password,
    });

    switch (res.status) {
      case 200:
        return { data: res.data };
      case 401:
        return { error: new ApiError("Invalid password", 401) };
      case 428:
        return { error: new ApiError("System configuration required", 428) };
      default:
        console.error(`Unexpected response from /login: ${res.status}`);
        return { error: new ApiError("Login failed, please try again", res.status) };
    }
  } catch {
    return { error: new ApiError("Unknown error, check hub is online", 0) };
  }
}

/**
 * Validate current session
 */
async function getSession(): Promise<ApiResponse<SessionResponse>> {
  try {
    const res = await hubApiClient.get<SessionResponse>(ApiRoute.SESSION);

    switch (res.status) {
      case 200:
        return { data: res.data };
      case 401:
        return { error: new ApiError("Login required", 401) };
      case 428:
        return { error: new ApiError("System configuration required", 428) };
      default:
        console.error(`Unexpected response from /session: ${res.status}`);
        return { error: new ApiError("Session validation failed", res.status) };
    }
  } catch (err) {
    console.error("Unexpected error from /session:", err);
    return { error: new ApiError("Unknown error during session validation", 0) };
  }
}

/**
 * Logout current session
 */
async function logout(): Promise<ApiResponse<LogoutResponse>> {
  try {
    const res = await hubApiClient.delete<LogoutResponse>(ApiRoute.SESSION);

    if (res.status !== 200) {
      console.error(`Unexpected response from /logout: ${res.status}`);
      return { error: new ApiError("Logout failed", res.status) };
    }

    return { data: res.data };
  } catch (err) {
    console.error("Unexpected error from /logout:", err);
    return { error: new ApiError("Unknown error during logout", 0) };
  }
}

export const auth = {
  login,
  getSession,
  logout,
};
