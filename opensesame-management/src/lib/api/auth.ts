import { ApiResponse, hubApiClient } from "./api";

export interface AuthResponse {
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
async function login(password: string): Promise<ApiResponse<AuthResponse>> {
    try {
        const res = await hubApiClient.post<AuthResponse>("/login", { password });

        switch (res.status) {
            case 200:
                return { data: res.data };
            case 401:
                return { error: new Error("Invalid password") };
            case 428:
                return { error: new Error("System configuration required") };
            default:
                console.error(`Unexpected response from /login: ${res.status}`);
                return { error: new Error("Login failed, please try again") };
        }
    } catch (err) {
        return { error: new Error("Unknown error, check hub is online") };
    }
}

/**
 * Validate current session
 */
async function getSession(): Promise<ApiResponse<AuthResponse>> {
    try {
        const res = await hubApiClient.get<AuthResponse>("/session");

        if (res.status === 401) {
            return { error: new Error("Login required") };
        }

        if (res.status === 428) {
            return { error: new Error("System configuration required") };
        }

        if (res.status !== 200) {
            console.error(`Unexpected response from /session: ${res.status}`);
            return { error: new Error("Session validation failed") };
        }

        return { data: res.data };
    } catch (err) {
        console.error("Unexpected error from /session:", err);
        return { error: new Error("Unknown error during session validation") };
    }
}

/**
 * Logout current session
 */
async function logout(): Promise<ApiResponse<LogoutResponse>> {
    try {
        const res = await hubApiClient.delete<LogoutResponse>("/logout");

        if (res.status !== 200) {
            console.error(`Unexpected response from /logout: ${res.status}`);
            return { error: new Error("Logout failed") };
        }

        return { data: res.data };
    } catch (err) {
        console.error("Unexpected error from /logout:", err);
        return { error: new Error("Unknown error during logout") };
    }
}

export const auth = {
    login,
    getSession,
    logout,
};