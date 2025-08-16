"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import api from "@/lib/api/api";
import { LoadState } from "@/types/load-state";
import { ConfigPatch, ConfigResponse } from "@/lib/api/config";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export const SessionTimeouts: Record<number, string> = {
  3600: "1 hour",
  43200: "12 hours",
  86400: "24 hours",
  259200: "3 days",
  604800: "7 days",
};

function getFormSchema(changingPassword: boolean) {
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
      password: changingPassword
        ? z.string().min(8, { message: "Must be at least 8 characters." })
        : z.string().optional(),
      newPassword: changingPassword
        ? z.string().min(8, { message: "Must be at least 8 characters." })
        : z.string().optional(),
      newPasswordConfirm: changingPassword
        ? z.string().min(8, { message: "Must be at least 8 characters." })
        : z.string().optional(),
    })
    .superRefine((data, ctx) => {
      if (changingPassword) {
        if (data.newPasswordConfirm !== data.newPassword) {
          ctx.addIssue({
            code: "custom",
            message: "New passwords do not match",
            path: ["newPasswordConfirm"],
          });
        }
      }
    });
}

type FormSchema = z.infer<ReturnType<typeof getFormSchema>>;
type UpdatableKeys = "systemName" | "sessionTimeoutSec";
function hasChanged<T extends UpdatableKeys>(key: T, data: FormSchema, config: ConfigResponse | undefined): boolean {
  if (!config) return false;
  return (data[key] as unknown) !== (config[key] as unknown);
}

export function SettingsForm() {
  // data states
  const [config, setConfig] = useState<ConfigResponse | undefined>();
  const [loadState, setLoadState] = useState<LoadState>(LoadState.LOADING);
  const [fetchError, setFetchError] = useState<string>();

  // ui states
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [authMethod, setAuthMethod] = useState<"password" | "backup-code">("password");

  const form = useForm<z.infer<ReturnType<typeof getFormSchema>>>({
    resolver: zodResolver(getFormSchema(showPasswordModal)),
    defaultValues: {
      systemName: "",
      password: "",
      newPassword: "",
      newPasswordConfirm: "",
      sessionTimeoutSec: 86400,
    },
  });

  const watchedSystemName = form.watch("systemName");
  const watchedPassword = form.watch("password");
  const watchedNewPassword = form.watch("newPassword");
  const watchedNewPasswordConfirm = form.watch("newPasswordConfirm");
  const watchedSessionTimeout = form.watch("sessionTimeoutSec");

  const watchedData = {
    systemName: watchedSystemName,
    password: watchedPassword,
    newPassword: watchedNewPassword,
    newPasswordConfirm: watchedNewPasswordConfirm,
    sessionTimeoutSec: watchedSessionTimeout,
  };

  useEffect(() => {
    async function fetchConfig() {
      setLoadState(LoadState.LOADING);
      const { data, error } = await api.config.GET();

      if (error) {
        setFetchError(error.message);
        setLoadState(LoadState.ERROR);
        return;
      }

      setConfig(data);
      setLoadState(LoadState.SUCCESS);
      form.reset({
        ...form.getValues(),
        systemName: data.systemName ?? "",
        sessionTimeoutSec: data.sessionTimeoutSec ?? 86400,
      });
    }

    fetchConfig();
  }, [form]);

  function onSubmit(data: z.infer<ReturnType<typeof getFormSchema>>) {
    setLoadState(LoadState.LOADING);
    setFetchError("");

    if (!config) {
      setFetchError("Error, configuration not loaded yet.");
      setLoadState(LoadState.ERROR);
      return;
    }

    const request: ConfigPatch = {};

    if (hasChanged("systemName", data, config)) {
      request.systemName = data.systemName.trim();
    }

    if (hasChanged("sessionTimeoutSec", data, config)) {
      request.sessionTimeoutSec = data.sessionTimeoutSec;
    }

    if (showPasswordModal && data.password && data.newPassword && data.newPasswordConfirm) {
      if (data.newPassword !== data.newPasswordConfirm) {
        setFetchError("New passwords do not match.");
        setLoadState(LoadState.ERROR);
        return;
      }

      if (authMethod === "password") {
        request.password = data.password;
      } else if (authMethod === "backup-code") {
        request.backupCode = data.password;
      }
      request.newPassword = data.newPassword;
    }

    if (Object.keys(request).length === 0) {
      setFetchError("Nothing to update.");
      setLoadState(LoadState.ERROR);
      return;
    }

    api.config.PATCH(request).then(({ data: resData, error }) => {
      if (error) {
        setFetchError(error.message ?? "Failed to update");
        setLoadState(LoadState.ERROR);
        return;
      }

      toast.success("Success!", {
        description: "Configuration updated successfully.",
      });

      setLoadState(LoadState.SUCCESS);

      if (resData) {
        setConfig((prev) => ({ ...prev, ...resData }));
      }
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
        <div className="mb-6 space-y-8 border-b pb-8">
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
                  value={field.value ? String(field.value) : undefined}
                >
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

          <div className="flex items-center justify-between">
            <div className="flex-col space-y-2">
              <FormLabel>Change Password</FormLabel>
              <FormDescription>Update your system&apos;s admin password</FormDescription>
            </div>
            <Button type="button" variant="outline" onClick={() => setShowPasswordModal((prev) => !prev)}>
              {showPasswordModal ? "Cancel" : "Change Password"}
            </Button>
          </div>

          {showPasswordModal && (
            <Card className="bg-secondary/25 dark:bg-secondary/50 p-4 shadow-none">
              <div className="flex items-center gap-4">
                <FormLabel htmlFor="auth-method">Authentication Method</FormLabel>
                <Switch
                  id="auth-method"
                  checked={authMethod === "backup-code"}
                  onCheckedChange={(checked) => setAuthMethod(checked ? "backup-code" : "password")}
                />
                <Badge variant={"default"}>{authMethod === "password" ? "Password" : "Backup Code"}</Badge>
              </div>

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Current {authMethod === "password" ? "Password" : "Backup Code"}</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        className="max-w-[400px]"
                        placeholder={`Enter your ${authMethod === "password" ? "password" : "backup code"}`}
                        disabled={loadState === LoadState.LOADING}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="newPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New Admin Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        className="max-w-[400px]"
                        placeholder="Enter your new admin password"
                        disabled={loadState === LoadState.LOADING}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>This password will be used to access admin settings.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="newPasswordConfirm"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Confirm New Admin Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        className="max-w-[400px]"
                        placeholder="Confirm your new admin password"
                        disabled={loadState === LoadState.LOADING}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </Card>
          )}
        </div>

        <div className="flex items-center justify-end gap-4">
          <p className="text-destructive text-sm">{fetchError}</p>
          <Button
            type="submit"
            disabled={
              loadState === LoadState.LOADING ||
              (!hasChanged("systemName", watchedData, config) &&
                !hasChanged("sessionTimeoutSec", watchedData, config) &&
                !showPasswordModal)
            }
          >
            {loadState === LoadState.LOADING ? "Saving..." : "Submit"}
          </Button>
        </div>
      </form>
    </Form>
  );
}
