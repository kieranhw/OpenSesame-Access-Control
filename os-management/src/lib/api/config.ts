import { ApiResponse, ApiRoute, hubApiClient } from "./api";
import { ApiError } from "./api-error";

export type ConfigResponse = {
  configured: boolean;
  systemName: string;
  backupCode?: string;
  sessionTimeoutSec: number;
};

export type ConfigPost = {
  systemName: string;
  adminPassword: string;
  sessionTimeoutSec: number;
};

export type ConfigPatch = {
  systemName?: string;
  backupCode?: string;
  password?: string;
  newPassword?: string;
  sessionTimeoutSec?: number;
};

async function GET(): Promise<ApiResponse<ConfigResponse>> {
  try {
    const res = await hubApiClient.get<ConfigResponse>(ApiRoute.PUBLIC_CONFIG);
    switch (res.status) {
      case 200:
        return { data: res.data };
      case 428:
        return { error: new ApiError("System configuration required", 428) };
      default:
        console.error(`Unexpected response from /config: ${res.status}`);
        return {
          error: new ApiError("Config retrieval failed, please try again", res.status),
        };
    }
  } catch (error) {
    console.error("Error fetching config:", error);
    return { error: new ApiError("Unknown error, check hub is online", 0) };
  }
}

async function POST(request: ConfigPost): Promise<ApiResponse<ConfigResponse>> {
  try {
    const res = await hubApiClient.post<ConfigResponse>(ApiRoute.PUBLIC_CONFIG, request);
    switch (res.status) {
      case 201:
        return { data: res.data };
      case 409:
        return { error: new ApiError("System already configured", 409) };
      case 428:
        return { error: new ApiError("System configuration required", 428) };
      default:
        console.error(`Unexpected response from /config: ${res.status}`);
        return {
          error: new ApiError("Failed to create configuration, please try again", res.status),
        };
    }
  } catch (error) {
    console.error("Error creating config:", error);
    return { error: new ApiError("Unknown error, check hub is online", 0) };
  }
}

async function PATCH(request: ConfigPatch): Promise<ApiResponse<ConfigResponse>> {
  try {
    const res = await hubApiClient.patch<ConfigResponse>(ApiRoute.ADMIN_CONFIG, request);
    switch (res.status) {
      case 200:
        return { data: res.data };
      case 400:
        return { error: new ApiError("Error, " + res.data, 400) };
      case 401:
        return { error: new ApiError("Unauthorized, please log in", 401) };
      case 428:
        return { error: new ApiError("System configuration required", 428) };
      default:
        console.error("Unexpected response from /config", res);
        return {
          error: new ApiError("Failed to update configuration, please try again", res.status),
        };
    }
  } catch (error) {
    console.error("Error updating config:", error);
    return { error: new ApiError("Unknown error, check hub is online", 0) };
  }
}

export const config = {
  GET,
  POST,
  PATCH,
};
