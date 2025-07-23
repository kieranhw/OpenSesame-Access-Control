"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { LoginForm } from "@/components/login-form";
import { api, ApiRoute } from "@/lib/api/api";
import packageJson from "../../../package.json";

interface LoginResponse {
  success: boolean;
}

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();
  const router = useRouter();

  const handleLogin = async (password: string) => {
    setError(undefined);
    setLoading(true);

    try {
      const { data } = await api.post<LoginResponse>(ApiRoute.LOGIN, {
        password,
      });

      if (data.success) {
        router.push("/");
      } else {
        
      }
    } catch (err) {
      console.error("Login error", err);
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || err.message || "Login failed");
      } else {
        setError("Login failed");
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
