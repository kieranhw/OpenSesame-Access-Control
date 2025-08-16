"use client";

import { Button } from "@/components/ui/button";
import { AppRoute } from "@/lib/app-routes";
import { Label } from "@radix-ui/react-label";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

interface SetupSuccessProps {
  backupCode: string;
}

export default function SetupSuccess({ backupCode }: SetupSuccessProps) {
  const router = useRouter();

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(backupCode);
      toast.success("Backup code copied to clipboard!");
    } catch {
      toast.error("Failed to copy backup code.");
    }
  }

  function navigateToLogin() {
    router.push(AppRoute.LOGIN);
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <Label>Backup Code</Label>
        <p className="text-muted-foreground text-sm">
          Please save the following backup code securely. You will need it to access your system if you lose your admin
          password.
        </p>
      </div>
      <div className="bg-muted rounded-md border p-4 font-mono text-lg">{backupCode}</div>
      <div className="flex items-center justify-end gap-2 border-t pt-6">
        <Button onClick={copyToClipboard} className="w-fit">
          Copy Backup Code
        </Button>
        <Button onClick={navigateToLogin} className="w-fit" variant="secondary">
          Continue
        </Button>
      </div>
    </div>
  );
}
