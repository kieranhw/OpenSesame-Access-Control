import { SectionCard } from "@/components/section-card";
import { DoorClosedLocked } from "lucide-react";
import { DataTableDemo } from "./entry-table";

async function SettingsPage() {
  return (
    <div className="mx-auto flex w-full flex-1 flex-col gap-4 p-4 sm:p-8 max-w-6xl">
      <h1 className="tracking-narrow flex items-center gap-2 text-xl font-semibold">
        <DoorClosedLocked className="h-6 w-6" />
        {"Entry Devices"}
      </h1>
      <DataTableDemo />
    </div>
  );
}

export default SettingsPage;
