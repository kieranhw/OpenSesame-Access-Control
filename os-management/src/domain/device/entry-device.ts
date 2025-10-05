import { BaseDevice } from "./base-device";
import { EntryDeviceData } from "../common/device";
import { ApiRoute, hubApiClient } from "@/lib/api/api";

export class EntryDevice extends BaseDevice {
  declare data: EntryDeviceData;

  constructor(data: EntryDeviceData) {
    super(data);
  }

  async rename(newName: string): Promise<void> {
    const response = await hubApiClient.patch(ApiRoute.ADMIN_ENTRY_DEVICES + "/" + this.id, {
      name: newName,
    });

    const responseText = await response.data;
    if (!response) {
      throw new Error(`Failed to rename device, ${responseText ?? "unknown error"}`);
    }

    this.data.name = newName;
  }

  toggleLock() {
    console.log(`Toggling lock status for device ${this.id}`);
  }

  get isLocked(): boolean {
    return this.data.lockStatus === "LOCKED";
  }
}
