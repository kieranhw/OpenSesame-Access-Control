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
import { ConfigResponse } from "@/types/config-response";
import { LoadState } from "@/types/load-state";
import { ConfigPatch } from "@/lib/api/config";

function getFormSchema(changingPassword: boolean) {
  return z
    .object({
      systemName: z
        .string()
        .min(1, { message: "Must be at least 1 character." })
        .max(50, { message: "Must be at most 50 characters." }),
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

const hasSystemNameChanged = (
  data: Pick<z.infer<ReturnType<typeof getFormSchema>>, "systemName">,
  config: ConfigResponse | undefined,
) => {
  if (!config) return false;
  return data.systemName.trim() !== config.systemName;
};

export function SettingsForm() {
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [authMethod, setAuthMethod] = useState<"password" | "backup-code">("password");
  const [fetchError, setFetchError] = useState<string>();
  const [config, setConfig] = useState<ConfigResponse | undefined>();
  const [loadState, setLoadState] = useState<LoadState>(LoadState.LOADING);

  const form = useForm<z.infer<ReturnType<typeof getFormSchema>>>({
    resolver: zodResolver(getFormSchema(showPasswordModal)),
    defaultValues: {
      systemName: "",
      password: "",
      newPassword: "",
      newPasswordConfirm: "",
    },
  });

  const watchedSystemName = form.watch("systemName");
  const watchedPassword = form.watch("password");
  const watchedNewPassword = form.watch("newPassword");
  const watchedNewPasswordConfirm = form.watch("newPasswordConfirm");

  const watchedData = {
    systemName: watchedSystemName,
    password: watchedPassword,
    newPassword: watchedNewPassword,
    newPasswordConfirm: watchedNewPasswordConfirm,
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
      });
    }

    fetchConfig();
  }, []);

  function onSubmit(data: z.infer<ReturnType<typeof getFormSchema>>) {
    if (!config) {
      toast.error("Error", { description: "Configuration not loaded yet." });
      return;
    }

    const request: ConfigPatch = {};

    if (hasSystemNameChanged(data, config)) {
      request.systemName = data.systemName.trim();
    }

    if (showPasswordModal && data.password && data.newPassword && data.newPasswordConfirm) {
      if (data.newPassword !== data.newPasswordConfirm) {
        setFetchError("New passwords do not match.");
        toast.error("Error", { description: "New passwords do not match." });
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
      toast.error("Error", { description: "Nothing to update." });
      return;
    }

    setLoadState(LoadState.LOADING);
    api.config.PATCH(request).then(({ data: resData, error }) => {
      if (error) {
        setLoadState(LoadState.ERROR);
        setFetchError(error.message);
        toast("Error", { description: error.message });
        return;
      }

      toast.success("Success!", {
        description: "Configuration updated successfully.",
      });

      if (resData) {
        setConfig((prev) => ({ ...prev, ...resData }));
      }

      setLoadState(LoadState.SUCCESS);
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
        <div className="mb-8 space-y-8 border-b pb-8">
          {/* System Name */}
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

          {/* Change Password Toggle */}
          <div className="flex items-center justify-between">
            <div className="flex-col space-y-2">
              <FormLabel>Change Password</FormLabel>
              <FormDescription>Update your system&apos;s admin password</FormDescription>
            </div>
            <Button type="button" variant="outline" onClick={() => setShowPasswordModal((prev) => !prev)}>
              {showPasswordModal ? "Cancel" : "Change Password"}
            </Button>
          </div>

          {/* Password Modal */}
          {showPasswordModal && (
            <Card className="bg-secondary/50 p-4 shadow-none">
              {/* Auth Method Switch */}
              <div className="flex items-center gap-4">
                <FormLabel htmlFor="auth-method">Authentication Method</FormLabel>
                <Switch
                  id="auth-method"
                  checked={authMethod === "backup-code"}
                  onCheckedChange={(checked) => setAuthMethod(checked ? "backup-code" : "password")}
                />
                <Badge variant={"default"}>{authMethod === "password" ? "Password" : "Backup Code"}</Badge>
              </div>

              {/* Current Password / Backup Code */}
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

              {/* New Admin Password */}
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

              {/* Confirm New Admin Password */}
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

        {/* Submit */}
        <div className="flex items-center justify-end gap-4">
          <p className="text-destructive text-sm">{fetchError}</p>
          <Button
            type="submit"
            disabled={
              loadState === LoadState.LOADING || (!hasSystemNameChanged(watchedData, config) && !showPasswordModal)
            }
          >
            {loadState === LoadState.LOADING ? "Saving..." : "Submit"}
          </Button>
        </div>
      </form>
    </Form>
  );
}
