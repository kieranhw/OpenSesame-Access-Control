"use client";
import { ColumnDef } from "@tanstack/react-table";
import { EntryDevice } from "@/types/device";

export const deviceNameCol: ColumnDef<EntryDevice> = {
  accessorKey: "name",
  header: "Device Name",
  cell: ({ row }) => {
    const name = row.getValue("name") as string | undefined;
    return <div>{name ?? "—"}</div>;
  },
};
