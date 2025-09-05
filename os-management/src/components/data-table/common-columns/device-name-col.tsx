"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice, EntryDevice } from "@/types/device";

export const deviceNameCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "name",
  header: "Device Name",
  cell: ({ row }) => {
    const name = row.getValue("name") as string | undefined;
    return <div>{name ?? "—"}</div>;
  },
});
