"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const portCol: ColumnDef<BaseDevice> = {
  accessorKey: "port",
  header: "Port",
  cell: ({ row }) => {
    const port = row.getValue("port") as number;
    return <div>{port}</div>;
  },
};
