import { AppSidebar } from "@/components/sidebar/app-sidebar";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Breadcrumbs } from "./breadcrumbs";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <AppSidebar />
      {/* Make SidebarInset a full-height flex column */}
      <SidebarInset className="flex flex-col min-h-screen">
        <nav
          className="sticky top-0 z-50 h-12 border-b bg-background/25 backdrop-blur-md 
                     shrink-0 flex items-center gap-2 px-4 transition-[width,height] ease-linear 
                     group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
        >
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mr-2 data-[orientation=vertical]:h-4"
          />
          <Breadcrumbs />
        </nav>

        {/* This is the scrollable, growing content area */}
        <main className="flex-1 flex flex-col">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}