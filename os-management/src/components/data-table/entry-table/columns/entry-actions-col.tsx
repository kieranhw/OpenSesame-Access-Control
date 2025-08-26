"use client";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal } from "lucide-react";
import { EntryDevice } from "@/types/device";

export const actionsCol: ColumnDef<EntryDevice> = {
  id: "actions",
  enableHiding: false,
  cell: ({ row }) => {
    const device = row.original;

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <span className="sr-only">Open menu</span>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={() => console.log("Details", device)}>Details</DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={() => console.log("Edit Name", device)}>Edit Name</DropdownMenuItem>
          <DropdownMenuItem variant="destructive" onClick={() => console.log("Remove", device)}>
            Remove
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  },
};
