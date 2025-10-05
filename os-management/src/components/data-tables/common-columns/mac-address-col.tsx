"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const macAddressCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "macAddress",
  header: "MAC Address",
  accessorFn: (device) => device.data.macAddress,
  cell: ({ getValue }) => {
    const mac = getValue<string>();
    return <div className="font-mono">{mac}</div>;
  },
});
