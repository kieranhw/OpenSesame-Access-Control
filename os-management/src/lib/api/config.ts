import { ApiResponse, hubApiClient } from "./api";

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
    const res = await hubApiClient.get<ConfigResponse>("/config");
    switch (res.status) {
      case 200:
        return { data: res.data };
      case 428:
        return { error: new Error("System configuration required") };
      default:
        console.error(`Unexpected response from /config: ${res.status}`);
        return { error: new Error("Config retrieval failed, please try again") };
    }
  } catch {
    return { error: new Error("Unknown error, check hub is online") };
  }
}

async function POST(request: ConfigPost): Promise<ApiResponse<ConfigResponse>> {
  try {
    const res = await hubApiClient.post<ConfigResponse>("/config", request);
    switch (res.status) {
      case 201:
        return { data: res.data };
      case 409:
        return { error: new Error("System already configured") };
      case 428:
        return { error: new Error("System configuration required") };
      default:
        console.error(`Unexpected response from /config: ${res.status}`);
        return { error: new Error("Failed to create configuration, please try again") };
    }
  } catch {
    return { error: new Error("Unknown error, check hub is online") };
  }
}

async function PATCH(request: ConfigPatch): Promise<ApiResponse<ConfigResponse>> {
  try {
    const res = await hubApiClient.patch<ConfigResponse>("/config", request);
    switch (res.status) {
      case 200:
        return { data: res.data };
      case 400:
        return { error: new Error("Error, " + res.data)}
      case 401:
        return { error: new Error("Unauthorized, please log in") };
      case 428:
        return { error: new Error("System configuration required") };
      default:
        console.error("Unexpected response from /config", res);
        return { error: new Error("Failed to update configuration, please try again") };
    }
  } catch {
    return { error: new Error("Unknown error, check hub is online") };
  }
}

export const config = {
  GET,
  POST,
  PATCH,
};
