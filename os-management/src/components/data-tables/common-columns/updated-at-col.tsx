"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const updatedAtCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "updatedAt",
  header: "Updated At",
  cell: ({ row }) => {
    const updatedAt = row.getValue("updatedAt") as string;
    const date = new Date(updatedAt);
    const formattedDate = date.toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "2-digit",
      year: "2-digit",
    });
    return <div>{formattedDate}</div>;
  },
});
