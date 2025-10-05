export class UnknownDeviceTypeError extends Error {
  readonly code = "UNKNOWN_DEVICE_TYPE";
  constructor(
    public readonly deviceType: string,
    public readonly deviceId?: number,
  ) {
    super(`Unknown device type received "${deviceType}"` + (deviceId ? ` for device id "${deviceId}"` : ""));
    this.name = "UnknownDeviceTypeError";
  }
}
