export type DeviceType = "entry" | "other";

export interface BaseDeviceData {
  id: number;
  macAddress: string;
  ipAddress: string;
  port: number;
  deviceType: DeviceType;
  instanceType: string;
  instanceName?: string;
  name?: string;
  description?: string;
  isOnline: boolean;
  lastSeen: number;
  createdAt: number;
  updatedAt: number;
}

export interface EntryDeviceData extends BaseDeviceData {
  deviceType: "entry";
  lockStatus: "LOCKED" | "UNLOCKED" | "UNKNOWN";
}
