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
import { useState } from "react";
import { Dialog } from "@/components/ui/dialog";
import { DeviceDetailsDialog } from "@/components/dialogs/device-details";
import { EditNameDialog } from "@/components/dialogs/edit-name-dialog";
import { EntryDevice } from "@/domain/device/entry-device";

export const entryActionsCol: ColumnDef<EntryDevice> = {
  id: "entryActions",
  enableHiding: false,
  cell: ({ row }) => {
    const device = row.original;
    const [dialogContent, setDialogContent] = useState<React.ReactNode>(null);
    const [isDialogOpen, setIsDialogOpen] = useState(false);

    const openDialog = (content: React.ReactNode) => {
      setDialogContent(content);
      setIsDialogOpen(true);
    };

    return (
      <>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => openDialog(<DeviceDetailsDialog device={device} />)}>
              Details
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem
              onClick={() => openDialog(<EditNameDialog device={device} onClose={() => setIsDialogOpen(false)} />)}
            >
              Edit Name
            </DropdownMenuItem>
            <DropdownMenuItem variant="destructive" onClick={() => console.log("Remove", device)}>
              Remove
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          {dialogContent}
        </Dialog>
      </>
    );
  },
};
