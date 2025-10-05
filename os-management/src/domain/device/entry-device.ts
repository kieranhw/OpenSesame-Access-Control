import { BaseDevice } from "./base-device";
import { EntryDeviceData } from "../common/device";

export class EntryDevice extends BaseDevice {
  declare data: EntryDeviceData;

  constructor(data: EntryDeviceData) {
    super(data);
  }

  async rename(newName: string): Promise<void> {
    // TODO: implement PATCH device API call
    //const { data, error } = await api.device.update();
    this.data.name = newName;
  }

  toggleLock() {
    console.log(`Toggling lock status for device ${this.id}`);
  }

  get isLocked(): boolean {
    return this.data.lockStatus === "LOCKED";
  }
}
