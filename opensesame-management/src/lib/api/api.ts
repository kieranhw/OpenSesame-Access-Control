import { auth } from "./auth";

import axios, { AxiosInstance } from "axios";

export type ApiResponse<T> =
    | { data: T; error?: never }
    | { data?: never; error: Error };

export const HUB_BASE_URI = "http://localhost:11072/management";

export enum ApiRoute {
    LOGIN = "/login",
    SESSION = "/session",
    LOGOUT = "/logout",
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

const api = {
    auth,
};

export default api;