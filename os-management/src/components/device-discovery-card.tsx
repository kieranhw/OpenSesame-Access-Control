"use client";

import { useState } from "react";
import { Card } from "./ui/card";
import { ChevronDown, ChevronUp, Radar } from "lucide-react";
import { Label } from "./ui/label";
import { Badge } from "./ui/badge";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useStatusContext } from "@/contexts/status-context";
import { timeAgo } from "@/lib/utils";
import { Button } from "./ui/button";
import pluralize from "pluralize";

export function DeviceDiscoveryCard() {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const { discoveredDevices } = useStatusContext();

  const hasDevices = discoveredDevices.length > 0;

  return (
    <Card className="gap-0 p-0">
      <div
        className="hover:bg-muted/50 dark:hover:bg-muted/60 flex justify-between px-6 py-4 transition hover:cursor-pointer"
        onClick={() => setIsOpen(!isOpen)}
      >
        <div className="flex items-center justify-start gap-6">
          <div className="relative flex h-10 w-10 items-center justify-center">
            <span className="bg-accent absolute inline-flex h-full w-full animate-ping rounded-full opacity-75"></span>
            <span className="bg-accent absolute inline-flex h-full w-full rounded-full"></span>
            <Radar className="text-muted-foreground relative h-5 w-5" />
          </div>

          <div>
            <Label>Device Discovery</Label>
            <Label className="text-muted-foreground text-sm">
              {hasDevices ? pluralize("device", discoveredDevices.length, true) + " found" : "Scanning for devices..."}
            </Label>
          </div>
        </div>

        <div className="flex items-center justify-start gap-6">
          <Badge variant="secondary" className="rounded-full px-3 py-1">
            {discoveredDevices.length}
          </Badge>
          {isOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
        </div>
      </div>

      {isOpen ? (
        hasDevices ? (
          <div className="border-t">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[160px] pl-6">MAC Address</TableHead>
                  <TableHead>Instance</TableHead>
                  <TableHead>IP Address</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Last Seen</TableHead>
                  <TableHead className="pr-6 text-right"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {discoveredDevices.map((d) => (
                  <TableRow key={d.id}>
                    <TableCell className="pl-6 font-mono">{d.mac_address}</TableCell>
                    <TableCell>{d.instance}</TableCell>
                    <TableCell>{d.ip_address}</TableCell>
                    <TableCell>{d.type ?? "-"}</TableCell>
                    <TableCell>{d.last_seen ? timeAgo(d.last_seen) : "-"}</TableCell>
                    <TableCell className="w-32 pr-6 text-right">
                      <Button variant={"secondary"} size={"sm"}>
                        Register
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        ) : (
          <div className="flex h-64 flex-col items-center justify-center gap-2 border-t px-4 text-center">
            <h1 className="text-xl font-semibold">No Devices Found</h1>
            <p className="text-muted-foreground">
              Make sure your devices are powered on, connected to Wi‑Fi, and have the software installed correctly.
            </p>
            <p className="text-muted-foreground">
              You can also add a device manually below, and OpenSesame will attempt to connect.
            </p>
          </div>
        )
      ) : null}
    </Card>
  );
}
