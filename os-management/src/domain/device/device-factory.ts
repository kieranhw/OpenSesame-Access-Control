import { BaseDevice } from "./base-device";
import { EntryDevice } from "./entry-device";
import { BaseDeviceData, EntryDeviceData } from "../common/device";
import { UnknownDeviceTypeError } from "./errors";

export class DeviceFactory {
  static create(data: BaseDeviceData): BaseDevice {
    switch (data.deviceType) {
      case "entry":
        return new EntryDevice(data as EntryDeviceData);
      default:
        throw new UnknownDeviceTypeError(data.deviceType, data.id);
    }
  }
}
