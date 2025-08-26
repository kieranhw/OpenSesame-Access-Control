export interface BaseDevice {
  id: number;
  macAddress: string;
  ipAddress: string;
  port: number;
}

export interface EntryDevice extends BaseDevice {
  name?: string;
  description?: string;
  createdAt: number;
  updatedAt: number;
}

export interface DiscoveredDevice extends BaseDevice {
  instance: string;
  type: string;
  lastSeen: number;
}
