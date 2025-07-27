"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { LoginForm } from "@/components/login-form";
import { login } from "@/lib/api/api";
import packageJson from "../../../package.json";
import { AppRoute } from "@/lib/app-routes";
import { LoginApiResponse } from "@/types/login-response";
import { ApiRoute } from "@/lib/api/api-routes";

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();
  const router = useRouter();

  const handleLogin = async (password: string) => {
    setError(undefined);
    setLoading(true);

    try {
      const response = await login(password);

      const loginData = response.data;
      if (loginData && loginData.authenticated && loginData.configured) {
        router.push(AppRoute.HOME);
      } else if (loginData && loginData.message) {
        setError(loginData.message);
      } else {
        setError("Login failed, please check your credentials.");
      }
    } catch (err) {
      console.error("Login API error:", err);
      if (axios.isAxiosError(err)) {
        const backendError = err.response?.data as
          | { message?: string; error?: string }
          | undefined;
        setError(
          backendError?.message ||
            backendError?.error ||
            err.message ||
            "Login failed",
        );
      } else {
        setError("An unexpected error occurred during login.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-muted flex h-full flex-col items-center justify-between overflow-scroll dark:bg-gradient-to-br dark:from-zinc-950 dark:via-black dark:to-zinc-900">
      <div />
      <div className="my-8 flex w-full max-w-96 flex-col gap-6">
        <LoginForm onSubmit={handleLogin} loading={loading} error={error} />
      </div>
      <footer className="border-divider bg-card w-full flex-none border-t px-4 py-2">
        <div className="flex h-8 items-center justify-center">
          <p className="text-muted-foreground text-sm">
            OpenSesame Management Version {packageJson.version}
          </p>
        </div>
      </footer>
    </div>
  );
}
