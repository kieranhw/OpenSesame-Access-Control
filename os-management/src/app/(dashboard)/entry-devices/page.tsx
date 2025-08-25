import { Button } from "@/components/ui/button";
import { SectionCard } from "@/components/section-card";
import { DeviceDiscoveryCard } from "@/components/device-discovery-card";
import { EntryDataTable } from "./entry-data-table";

async function EntryDevicesPage() {

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-4 p-4">
      <DeviceDiscoveryCard />
      <SectionCard
        title="Entry Devices"
        subheader="Manage your system's entry devices."
        className="pb-0"
        button={<Button>Create New</Button>}
        bodyPadding={false}
      >
        <div className="border-t">
          <EntryDataTable />
        </div>
      </SectionCard>
    </div>
  );
}

export default EntryDevicesPage;
