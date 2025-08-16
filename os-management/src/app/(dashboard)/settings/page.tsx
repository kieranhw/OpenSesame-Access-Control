import { SettingsForm } from "./settings-form";
import { MonitorCog } from "lucide-react";
import { SectionCard } from "@/components/section-card";

async function SettingsPage() {
  return (
    <div className="mx-auto my-6 flex w-full max-w-4xl flex-1 flex-col gap-4 px-4">
      <SectionCard title="System Settings" icon={MonitorCog}>
        <SettingsForm />
      </SectionCard>
    </div>
  );
}

export default SettingsPage;
