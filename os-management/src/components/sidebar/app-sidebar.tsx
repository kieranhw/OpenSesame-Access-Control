"use client";

import * as React from "react";
import { DoorClosedLocked, FileLock2, House, KeyRound, Microchip, MonitorSmartphone, ServerCog, Settings, UserLock } from "lucide-react";

import { NavMain } from "@/components/sidebar/nav-main";
import { Sidebar, SidebarContent, SidebarFooter, SidebarRail } from "@/components/ui/sidebar";
import { AppRoute } from "@/lib/app-routes";

const sidebarItems = [
  {
    title: "Home",
    url: AppRoute.HOME,
    icon: House,
  },
  {
    title: "Device Management",
    icon: Microchip,
    items: [
      {
        title: "Entry Devices",
        url: AppRoute.ENTRY_DEVICES,
        icon: DoorClosedLocked,
      },
      {
        title: "Access Devices",
        url: AppRoute.ACCESS,
        icon: KeyRound,
      },
      {
        title: "Pairing",
        url: AppRoute.PAIRING,
        icon: ServerCog,
      }
    ],
  },
  {
    title: "Clients",
    url: AppRoute.CLIENTS,
    icon: MonitorSmartphone,
  },
  {
    title: "Credentials",
    url: AppRoute.CREDENTIALS,
    icon: UserLock,
  },
  {
    title: "Settings",
    url: AppRoute.SETTINGS,
    icon: Settings,
  },
];

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarContent className="pt-12">
        <NavMain items={sidebarItems} />
      </SidebarContent>
      <SidebarFooter></SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
