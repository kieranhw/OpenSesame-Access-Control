import { Button } from "@/components/ui/button";
import { DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Separator } from "../ui/separator";
import { BaseDevice } from "@/domain/device/base-device";

interface DeviceDetailsDialogProps {
  device: BaseDevice;
}

export function DeviceDetailsDialog({ device }: DeviceDetailsDialogProps) {
  function renderRowItem(label: string, detail?: string | number | null) {
    if (!detail) return null;

    return (
      <div className="flex items-start gap-4">
        <Label className="text-muted-foreground w-32 shrink-0">{label}</Label>
        <p className="flex-1 text-sm">{detail}</p>
      </div>
    );
  }

  return (
    <DialogContent className="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{device ? "Entry Device" : "Device"} Details</DialogTitle>
      </DialogHeader>

      <div className="space-y-2">
        {renderRowItem("Name", device.data.name)}
        {renderRowItem("Description", device.data.description)}
        {renderRowItem("Instance", `${device.data.instanceType}`)}
        {renderRowItem("Status", device.data.isOnline ? "Online" : "Offline")}
      </div>
      <Separator />
      <div className="space-y-2">
        {renderRowItem("MAC Address", device.data.macAddress)}
        {renderRowItem("IP Address", device.data.ipAddress)}
        {renderRowItem("Port", device.data.port)}
      </div>
      <DialogFooter>
        <DialogClose asChild>
          <Button variant="outline">Close</Button>
        </DialogClose>
      </DialogFooter>
    </DialogContent>
  );
}
