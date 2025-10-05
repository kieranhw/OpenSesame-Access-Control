"use client";
import { BaseDevice } from "@/domain/device/base-device";
import { ColumnDef } from "@tanstack/react-table";

export const instanceCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  id: "instance",
  header: "Instance",
  cell: ({ row }) => {
    const instanceType: string = row.original.data.instanceType;
    const instanceName: string | undefined = row.original.data.instanceName;

    return (
      <div className="flex flex-col">
        <span>{instanceName ?? "—"}</span>
        <span>{instanceType}</span>
      </div>
    );
  },
});
