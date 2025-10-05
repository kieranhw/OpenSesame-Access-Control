"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const portCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "port",
  header: "Port",
  accessorFn: (device) => device.data.port,
  cell: ({ getValue }) => {
    const port = getValue<number>();
    return <div>{port}</div>;
  },
});
