export interface EntryDevice {
  id: number;
  name?: string;
  description?: string;
  macAddress: string;
  ipAddress: string;
  port: number;
  createdAt: number;
  updatedAt: number;
}

export interface DiscoveredDevice {
  id: number;
  instance: string;
  ipAddress: string;
  macAddress: string;
  type: string;
  lastSeen: number;
}

export interface StatusResponse {
  etag: number;
  systemName: string;
  entryDevices: EntryDevice[];
  discoveredDevices: DiscoveredDevice[];
  //   accessDevices: AccessDevice[];
}
