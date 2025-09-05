"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const portCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "port",
  header: "Port",
  cell: ({ row }) => {
    const port = row.getValue("port") as number;
    return <div>{port}</div>;
  },
});
