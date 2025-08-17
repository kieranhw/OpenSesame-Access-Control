import { auth } from "./session";
import axios, { AxiosInstance } from "axios";
import camelcaseKeys from "camelcase-keys";
import snakecaseKeys from "snakecase-keys";
import { config } from "./config";

export type ApiResponse<T> = { data: T; error?: never } | { data?: never; error: Error };

export const HUB_BASE_URI = "http://localhost:11072/management";

export enum ApiRoute {
  SESSION = "/session",
  CONFIG = "/config",
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
};

export default api;
