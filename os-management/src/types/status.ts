export interface EntryDevice {
  id: number;
  ip_address: string;
}

export interface DiscoveredDevice {
  id: number;
  instance: string;
  ip_address: string;
  mac_address: string;
  type?: string;
  last_seen?: number;
}

export interface AccessDevice {
  id: number;
  name: string;
}

export interface StatusResponse {
  etag: number;
  system_name: string;
  entry_devices: EntryDevice[];
  discovered_devices: DiscoveredDevice[];
  //   access_devices: AccessDevice[];
}
