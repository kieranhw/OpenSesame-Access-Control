import { Card } from "@/components/ui/card";
import { SettingsForm } from "./settings-form";

async function SettingsPage() {
  return (
    <div className="mx-auto my-6 flex w-full max-w-4xl flex-1 flex-col gap-4">
      <Card className="flex w-full flex-1 flex-col gap-4 p-6">
        <h1 className="tracking-narrow text-3xl">System Settings</h1>
        <SettingsForm />
      </Card>
    </div>
  );
}

export default SettingsPage;
