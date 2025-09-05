"use client";
import { useStatusContext } from "@/contexts/status-context";
import { DataTable } from "../data-table";
import { useMemo, useState } from "react";
import { EntryDevice } from "@/types/device";
import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { selectColumn } from "../common-columns/select-col";
import { deviceNameCol } from "../common-columns/device-name-col";
import { macAddressCol } from "../common-columns/mac-address-col";
import { portCol } from "../common-columns/port-col";
import { ipAddressCol } from "../common-columns/ip-address-col";
import { updatedAtCol } from "../common-columns/updated-at-col";
import { entryActionsCol } from "./columns/entry-actions-col";
import { statusCol } from "../common-columns/status-col";

export function EntryDataTable() {
  const { entryDevices } = useStatusContext();
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  const columns = useMemo<ColumnDef<EntryDevice>[]>(() => {
    const cols: ColumnDef<EntryDevice>[] = [];
    cols.push(selectColumn<EntryDevice>());
    cols.push(deviceNameCol<EntryDevice>());
    cols.push(macAddressCol<EntryDevice>());
    cols.push(ipAddressCol<EntryDevice>());
    cols.push(portCol<EntryDevice>());
    cols.push(statusCol<EntryDevice>());

    cols.push(entryActionsCol);
    return cols;
  }, [entryDevices]);

  return (
    <DataTable
      data={entryDevices}
      columns={columns}
      rowSelection={rowSelection}
      onRowSelectionChange={setRowSelection}
    />
  );
}
