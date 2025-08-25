import { StatusResponse } from "@/types/status";
import { ApiResponse, hubApiClient } from "./api";
import { ApiError } from "./api-error";

async function GET(): Promise<ApiResponse<StatusResponse>> {
  try {
    const res = await hubApiClient.get<StatusResponse>("/status");
    switch (res.status) {
      case 200:
        return { data: res.data };
      default:
        console.error(`Unexpected response from /status: ${res.status}`);
        return { error: new ApiError("Failed to fetch status", res.status) };
    }
  } catch {
    return { error: new ApiError("Unknown error, check hub is online", 0) };
  }
}

async function LONG_POLL(timeout: number, etag: number): Promise<ApiResponse<StatusResponse>> {
  try {
    const res = await hubApiClient.get<StatusResponse>("/status", {
      params: { timeout, etag },
      timeout: (timeout + 5) * 1000, // axios request timeout slightly > server timeout
    });

    switch (res.status) {
      case 200:
        return { data: res.data };
      case 304:
        return { error: new ApiError("Not Modified", 304) };
      default:
        console.error(`Unexpected response from /status: ${res.status}`);
        return { error: new ApiError("Failed to poll status", res.status) };
    }
  } catch {
    return { error: new ApiError("Hub unreachable", 0) };
  }
}

export const status = {
  GET,
  LONG_POLL,
};
