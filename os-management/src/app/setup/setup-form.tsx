"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import api from "@/lib/api/api";
import { LoadState } from "@/domain/common/load-state";
import { ConfigPost } from "@/lib/api/config";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export const SessionTimeouts: Record<number, string> = {
  3600: "1 hour",
  43200: "12 hours",
  86400: "24 hours",
  259200: "3 days",
  604800: "7 days",
};

function getFormSchema() {
  return z
    .object({
      systemName: z
        .string()
        .min(1, { message: "Must be at least 1 character." })
        .max(50, { message: "Must be at most 50 characters." }),
      sessionTimeoutSec: z
        .number()
        .min(3600, { message: "Must be at least 1 hour." })
        .max(86400 * 30, { message: "Must be at most 30 days." }),
      adminPassword: z
        .string()
        .min(8, { message: "Must be at least 8 characters." })
        .max(100, { message: "Must be at most 100 characters." }),
      adminPasswordConfirm: z
        .string()
        .min(8, { message: "Must be at least 8 characters." })
        .max(100, { message: "Must be at most 100 characters." }),
    })
    .superRefine((data, ctx) => {
      if (data.adminPassword !== data.adminPasswordConfirm) {
        ctx.addIssue({
          code: "custom",
          message: "Passwords do not match",
          path: ["adminPasswordConfirm"],
        });
      }
    });
}

interface SetupFormProps {
  onSuccess: (backupCode: string) => void;
}

export function SetupForm({ onSuccess }: SetupFormProps) {
  const [fetchError, setFetchError] = useState<string>();
  const [loadState, setLoadState] = useState<LoadState>(LoadState.IDLE);

  const form = useForm<z.infer<ReturnType<typeof getFormSchema>>>({
    resolver: zodResolver(getFormSchema()),
    defaultValues: {
      systemName: "",
      adminPassword: "",
      adminPasswordConfirm: "",
      sessionTimeoutSec: 86400,
    },
    mode: "onChange",
  });

  async function onSubmit(data: z.infer<ReturnType<typeof getFormSchema>>) {
    const request: ConfigPost = {
      systemName: data.systemName,
      adminPassword: data.adminPassword,
      sessionTimeoutSec: data.sessionTimeoutSec,
    };

    setLoadState(LoadState.LOADING);
    api.config.POST(request).then(({ data, error }) => {
      if (error) {
        setLoadState(LoadState.ERROR);
        setFetchError(error.message);
        toast("Error", { description: error.message });
        return;
      }

      if (data.backupCode) {
        onSuccess(data.backupCode);
      } else {
        toast.error("Error", {
          description: "No backup code returned from server.",
        });
      }

      setLoadState(LoadState.SUCCESS);
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
        <div className="mb-6 space-y-8">
          <FormField
            control={form.control}
            name="systemName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>System Name</FormLabel>
                <FormControl>
                  <Input
                    className="max-w-[400px]"
                    disabled={loadState === LoadState.LOADING}
                    placeholder={loadState === LoadState.LOADING ? "Loading..." : "Enter a name for your system"}
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  This is your system&apos;s name, which will be displayed in client applications.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="sessionTimeoutSec"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Session Timeout</FormLabel>
                <Select
                  onValueChange={(val) => field.onChange(Number(val))}
                  value={field.value ? String(field.value) : ""}
                >
                  {" "}
                  <FormControl>
                    <SelectTrigger className="w-full max-w-[400px]">
                      <SelectValue placeholder="Select a timeout" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.entries(SessionTimeouts).map(([seconds, label]) => (
                      <SelectItem key={seconds} value={seconds}>
                        {label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  The amount of time before a user is automatically logged out due to inactivity.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="adminPassword"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Admin Password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    className="max-w-[400px]"
                    placeholder="Enter your password"
                    disabled={loadState === LoadState.LOADING}
                    {...field}
                  />
                </FormControl>
                <FormDescription>Your admin password, which must be at least 8 characters.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="adminPasswordConfirm"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Confirm Admin Password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    className="max-w-[400px]"
                    placeholder="Confirm your password"
                    disabled={loadState === LoadState.LOADING}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <div className="flex items-center justify-end gap-2 border-t pt-6">
          <p className="text-destructive text-sm">{fetchError}</p>
          <Button type="submit" disabled={loadState === LoadState.LOADING}>
            {loadState === LoadState.LOADING ? "Creating..." : "Create System"}
          </Button>
        </div>
      </form>
    </Form>
  );
}

export default SetupForm;
