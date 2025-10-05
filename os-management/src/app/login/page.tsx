"use client";
import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { LoginForm } from "@/components/login-form";
import packageJson from "../../../package.json";
import { AppRoute } from "@/lib/app-routes";
import api from "@/lib/api/api";
import { LoadState } from "@/domain/common/load-state";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const urlError = searchParams.get("error") || undefined;

  const [error, setError] = useState<string | undefined>(urlError);
  const [loading, setLoading] = useState<LoadState>(LoadState.IDLE);

  useEffect(() => {
    if (urlError) setError(urlError);
  }, [urlError]);

  async function handleLogin(password: string) {
    setError(undefined);
    setLoading(LoadState.LOADING);

    const { data, error: loginError } = await api.auth.login(password);

    if (data && data.authenticated && data.configured) {
      router.push(AppRoute.HOME);
      return;
    }

    if (loginError) {
      setError(loginError.message);
      setLoading(LoadState.ERROR);
      return;
    }

    setError("Login failed");
    setLoading(LoadState.ERROR);
  }

  return (
    <div className="bg-muted flex h-full flex-col items-center justify-between overflow-scroll dark:bg-gradient-to-br dark:from-zinc-950 dark:via-black dark:to-zinc-900">
      <div />
      <div className="my-8 flex w-full max-w-96 flex-col gap-6">
        <LoginForm onSubmit={handleLogin} loading={loading === LoadState.LOADING} error={error} />
      </div>
      <footer className="border-divider bg-card w-full flex-none border-t px-4 py-2">
        <div className="flex h-8 items-center justify-center">
          <p className="text-muted-foreground text-sm">OpenSesame Management Version {packageJson.version}</p>
        </div>
      </footer>
    </div>
  );
}
