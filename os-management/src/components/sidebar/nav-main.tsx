"use client";

import { ChevronRight, type LucideIcon } from "lucide-react";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { JSX } from "react";
import Link from "next/link";

export interface NavSubItem {
  title: string;
  url: string;
  icon?: LucideIcon;
  isActive?: boolean;
}

export interface NavItem {
  title: string;
  url?: string; // optional when it’s a collapsible group
  icon?: LucideIcon;
  isActive?: boolean; // controls defaultOpen for groups
  items?: NavSubItem[];
}

export function NavMain({ items }: { items: NavItem[] }): JSX.Element {
  return (
    <SidebarGroup>
      <SidebarMenu>
        {items.map((item) =>
          item.items && item.items.length > 0 ? (
            <CollapsibleSidebarItem key={item.title} {...item} />
          ) : (
            <SidebarItem key={item.title} {...item} />
          )
        )}
      </SidebarMenu>
    </SidebarGroup>
  );
}

export function SidebarItem({
  title,
  url = "#",
  icon: Icon,
}: NavItem): JSX.Element {
  return (
    <SidebarMenuItem>
      <SidebarMenuButton asChild tooltip={title}>
        <Link href={url} draggable={false}>
          {Icon && <Icon />}
          <span>{title}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
}

export function CollapsibleSidebarItem({
  title,
  icon: Icon,
  items = [],
  isActive,
}: NavItem): JSX.Element {
  return (
    <Collapsible asChild defaultOpen={!!isActive} className="group/collapsible">
      <SidebarMenuItem>
        <CollapsibleTrigger asChild>
          <SidebarMenuButton tooltip={title}>
            {Icon && <Icon />}
            <span>{title}</span>
            <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
          </SidebarMenuButton>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <SidebarMenuSub>
            {items.map((subItem) => {
              const SubIcon = subItem.icon;
              return (
                <SidebarMenuSubItem key={subItem.title}>
                  <SidebarMenuSubButton asChild>
                    <Link href={subItem.url} draggable={false}>
                      {SubIcon && <SubIcon className="size-4" />}
                      <span>{subItem.title}</span>
                    </Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              );
            })}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  );
}