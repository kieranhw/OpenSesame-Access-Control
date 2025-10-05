import api, { ApiResponse } from "@/lib/api/api";
import { BaseDeviceData, EntryDeviceData } from "../common/device";

export abstract class BaseDevice {
  constructor(public data: BaseDeviceData) {}

  get id() {
    return this.data.id;
  }

  get label() {
    return this.data.name ?? this.data.instanceName ?? "Unnamed Device";
  }

  get isEntryDevice(): boolean {
    return this.data.deviceType === "entry";
  }

  get isOnline(): boolean {
    return this.data.isOnline;
  }

  abstract rename(newName: string): Promise<void>;
}
