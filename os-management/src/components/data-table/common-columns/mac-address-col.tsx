import { ColumnDef } from "@tanstack/react-table";
import { BaseDevice } from "@/types/device";

export const macAddressCol = <T extends BaseDevice>(): ColumnDef<T> => ({
  accessorKey: "macAddress",
  header: "MAC Address",
  cell: ({ row }) => {
    const mac = row.getValue("macAddress") as string;
    return <div className="font-mono">{mac}</div>;
  },
});