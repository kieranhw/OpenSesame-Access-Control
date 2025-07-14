// app/login/page.tsx
"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { LoginForm } from "@/components/login-form";
import { api } from "@/lib/api";
import { LOGIN_URL } from "@/lib/constants";

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
      const { data } = await api.post<LoginResponse>(LOGIN_URL, {
        password,
      });
      console.log("Login succeeded", data);
      router.push("/");
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
    <div className="bg-muted flex h-full flex-col items-center justify-center overflow-scroll">
      <div className="flex w-full max-w-96 flex-col gap-6">
        <LoginForm onSubmit={handleLogin} loading={loading} error={error} />
      </div>
    </div>
  );
}
