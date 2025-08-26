"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const ipAddressCol: ColumnDef<BaseDevice> = {
  accessorKey: "ipAddress",
  header: "IP Address",
  cell: ({ row }) => {
    const ip = row.getValue("ipAddress") as string;
    return <div className="font-mono">{ip}</div>;
  },
};
