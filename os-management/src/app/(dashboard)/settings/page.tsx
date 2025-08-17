import { SettingsForm } from "./settings-form";
import { MonitorCog } from "lucide-react";
import { SectionCard } from "@/components/section-card";

async function SettingsPage() {
  return (
    <div className="mx-auto flex w-full flex-1 flex-col gap-4 p-4 sm:p-8 max-w-6xl">
      <SectionCard title="System Settings" icon={MonitorCog}>
        <SettingsForm />
      </SectionCard>
    </div>
  );
}

export default SettingsPage;
