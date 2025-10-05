"use client";
import { useStatusContext } from "@/contexts/status-context";
import { DataTable } from "../data-table";
import { useState } from "react";
import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { deviceNameCol } from "../common-columns/device-name-col";
import { macAddressCol } from "../common-columns/mac-address-col";
import { portCol } from "../common-columns/port-col";
import { ipAddressCol } from "../common-columns/ip-address-col";
import { entryActionsCol } from "./columns/entry-actions-col";
import { statusCol } from "../common-columns/status-col";
import { EntryDevice } from "@/domain/device/entry-device";
import { instanceCol } from "../common-columns/instance-col";

const entryColumns: ColumnDef<EntryDevice>[] = [
  deviceNameCol<EntryDevice>(),
  instanceCol<EntryDevice>(),
  macAddressCol<EntryDevice>(),
  ipAddressCol<EntryDevice>(),
  portCol<EntryDevice>(),
  statusCol<EntryDevice>(),
  entryActionsCol,
];

export function EntryDataTable() {
  const { entryDevices } = useStatusContext();
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  return (
    <DataTable
      data={entryDevices}
      columns={entryColumns}
      rowSelection={rowSelection}
      onRowSelectionChange={setRowSelection}
    />
  );
}
