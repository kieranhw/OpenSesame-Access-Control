import { auth } from "./session";
import axios, { AxiosInstance } from "axios";
import camelcaseKeys from "camelcase-keys";
import snakecaseKeys from "snakecase-keys";
import { config } from "./config";
import { status } from "./status";
import { ApiError } from "./api-error";

export type ApiResponse<T = unknown> = { data?: T; error?: never } | { data?: never; error: ApiError };

export const HUB_BASE_URI = "http://localhost:11072";

export enum ApiRoute {
  PUBLIC_CONFIG = "/config",
  ADMIN_SESSION = "/admin/session",
  ADMIN_CONFIG = "/admin/config",
  ADMIN_STATUS = "/admin/status",
}

export const hubApiClient: AxiosInstance = axios.create({
  baseURL: HUB_BASE_URI,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
  validateStatus: () => true,
});

// convert outgoing requests to snake_case
hubApiClient.interceptors.request.use(
  (config) => {
    if (config.data && typeof config.data === "object") {
      config.data = snakecaseKeys(config.data, { deep: true });
    }
    if (config.params && typeof config.params === "object") {
      config.params = snakecaseKeys(config.params, { deep: true });
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// convert incoming responses to camelCase
hubApiClient.interceptors.response.use(
  (response) => {
    if (response.data && typeof response.data === "object") {
      response.data = camelcaseKeys(response.data, { deep: true });
    }
    return response;
  },
  (error) => Promise.reject(error),
);

const api = {
  auth,
  config,
  status,
};

export default api;
