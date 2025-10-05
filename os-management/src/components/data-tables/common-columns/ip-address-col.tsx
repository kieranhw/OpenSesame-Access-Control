"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const ipAddressCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "ipAddress",
  header: "IP Address",
  accessorFn: (device) => device.data.ipAddress,
  cell: ({ getValue }) => {
    const ip = getValue<string>();
    return <div className="font-mono">{ip}</div>;
  },
});
