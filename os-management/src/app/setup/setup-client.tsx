"use client";

import { CircleAlert, CircleCheckBig, Settings } from "lucide-react";
import { SectionCard } from "@/components/section-card";
import SetupForm from "./setup-form";
import SetupSuccess from "./setup-success";
import { useState } from "react";

export default function SetupClient() {
  const [step, setStep] = useState<"form" | "success">("form");
  const [backupCode, setBackupCode] = useState<string | null>(null);

  return (
    <>
      {step === "form" && (
        <SectionCard
          title="System Setup Required"
          icon={Settings}
          subheader="Before we begin, please configure the required settings for your system."
        >
          <SetupForm
            onSuccess={(code: string) => {
              setBackupCode(code);
              setStep("success");
            }}
          />
        </SectionCard>
      )}

      {step === "success" && backupCode && (
        <SectionCard
          title="System Configured!"
          icon={CircleCheckBig}
          subheader={
            <div className="mt-6 flex items-center justify-start gap-2 rounded border bg-zinc-800 p-3 text-sm sm:text-md text-white">
              <CircleAlert /> Your system has been configured, but actions are required.
            </div>
          }
        >
          <SetupSuccess backupCode={backupCode} />
        </SectionCard>
      )}
    </>
  );
}
