"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const deviceNameCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  id: "name",
  header: "Device Name",
  accessorFn: (device) => device.data.name,
  cell: ({ getValue }) => {
    const name = getValue<string | undefined>();
    return <div>{name ?? "—"}</div>;
  },
});
