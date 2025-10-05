import { EntryDeviceData } from "./device";

export interface StatusResponse {
  etag: number;
  systemName: string;
  entryDevices: EntryDeviceData[];
  //   accessDevices: AccessDevice[];
}
