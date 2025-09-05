"use client";
import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const statusCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "status",
  header: "Status",
  cell: ({ row }) => {
    const status = row.original.status.currentState;
    const statusStyle =
      status === "online" ? "bg-green-500" : status === "offline" ? "bg-red-500 animate-ping" : "bg-gray-500";

    return (
      <div>
        <span className="relative flex size-3">
          <span
            className={`absolute inline-flex h-full w-full animate-ping rounded-full ${statusStyle} opacity-75`}
          ></span>
          <span className="relative inline-flex size-3 rounded-full bg-zinc-500"></span>
        </span>
        <p>{status ?? "None"}</p>
      </div>
    );
  },
});
