import { useState } from "react";
import { Button } from "@/components/ui/button";
import { DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { BaseDevice } from "@/domain/device/base-device";
import { useStatusContext } from "@/contexts/status-context";

interface DeviceDetailsDialogProps {
  device: BaseDevice;
  onClose?: () => void;
}

export function EditNameDialog({ device, onClose }: DeviceDetailsDialogProps) {
  const { bumpState } = useStatusContext();
  const [isLoading, setIsLoading] = useState(false);
  const [nameValue, setNameValue] = useState<string>(device.data.name ?? "");

  const handleSubmit = async () => {
    try {
      await device.rename("Test Name");
      bumpState();
      // setIsLoading(true);
      // const response = await updateDeviceName(device.id, nameValue);
      // if (response.success) {
      //   console.log("✅ Name updated successfully");
      // }

      onClose?.();
    } catch (err) {
      console.error("❌ Failed to update name", err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <DialogContent className="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>Edit Entry Name</DialogTitle>
      </DialogHeader>

      <div className="my-2 space-y-1">
        <Label className="text-muted-foreground">Device Name</Label>
        <Input value={nameValue} onChange={(e) => setNameValue(e.target.value)} disabled={isLoading} />
      </div>

      <DialogFooter>
        <DialogClose asChild>
          <Button variant="outline" disabled={isLoading}>
            Close
          </Button>
        </DialogClose>

        <Button onClick={handleSubmit} disabled={isLoading || !nameValue.trim()}>
          {isLoading ? "Saving..." : "Submit"}
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
