// app/login/page.tsx
"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { LoginForm } from "@/components/login-form";
import { api } from "@/lib/api";

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
      // POST to http://localhost:11072/management/login
      const { data } = await api.post<LoginResponse>("/management/login", {
        password,
      });
      console.log("Login succeeded", data);
      // the browser has now stored os_session
      router.push("/");
    } catch (err: unknown) {
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
    <div className="bg-muted flex h-full flex-col items-center justify-center overflow-scroll">
      <div className="flex w-full max-w-96 flex-col gap-6">
        <LoginForm onSubmit={handleLogin} loading={loading} error={error} />
      </div>
    </div>
  );
}
