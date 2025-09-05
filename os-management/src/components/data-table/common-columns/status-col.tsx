"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const statusCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "status",
  header: "Status",
  cell: ({ row }) => {
    const isOnline = row.original.isOnline;

    return (
      <div className="flex items-center gap-2">
        <span className="relative flex size-3">
          <span
            className={`absolute inline-flex h-full w-full rounded-full ${isOnline ? "bg-green-500" : "animate-ping bg-destructive"} opacity-75`}
          ></span>
          <span
            className={`relative inline-flex size-3 rounded-full ${isOnline ? "bg-green-500" : "bg-destructive"}`}
          ></span>
        </span>
        <p>{isOnline ? "Online" : "Offline"}</p>
      </div>
    );
  },
});
