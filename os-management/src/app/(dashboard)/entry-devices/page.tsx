import { Button } from "@/components/ui/button";
import { SectionCard } from "@/components/section-card";
import { EntryDataTable } from "@/components/data-table/entry-table/entry-device-table";
import { PageLayout } from "@/components/page-layout";

async function EntryDevicesPage() {
  return (
    <PageLayout>
      <SectionCard
        title="Entry Devices"
        subheader="Manage your system's entry devices."
        className="pb-0"
        bodyPadding={false}
      >
        <div className="border-t">
          <EntryDataTable />
        </div>
      </SectionCard>
    </PageLayout>
  );
}

export default EntryDevicesPage;
