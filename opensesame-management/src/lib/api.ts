import axios from "axios";
import camelcaseKeys from "camelcase-keys";
import { BASE_URL } from "./constants";
import { ConfigResponse } from "@/app/types/config";

export const api = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

api.interceptors.response.use((resp) => {
  resp.data = camelcaseKeys(resp.data, { deep: true });
  return resp;
});

export function getConfig() {
  return api.get<ConfigResponse>(`/management/config`);
}
