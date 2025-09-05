interface DeviceStatus {
  currentState: "online" | "offline" | "unknown";
  lastSeen: number;
}

export interface BaseDevice {
  id: number;
  macAddress: string;
  ipAddress: string;
  port: number;
  status: DeviceStatus;
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
