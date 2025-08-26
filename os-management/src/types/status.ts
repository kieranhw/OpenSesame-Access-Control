import { DiscoveredDevice, EntryDevice } from "./device";

export interface StatusResponse {
  etag: number;
  systemName: string;
  entryDevices: EntryDevice[];
  discoveredDevices: DiscoveredDevice[];
  //   accessDevices: AccessDevice[];
}
