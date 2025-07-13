"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { LoginForm } from "@/components/login-form";

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();
  const router = useRouter();

  const handleLogin = async (password: string) => {
    setError(undefined);
    setLoading(true);
    try {
      const res = await fetch("/api/login", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ password }),
      });

      console.log("Login response:", res);
      if (!res.ok) {
        const text = await res.text();
        // TODO: read the error message and set it
        console.log("Login failed:", text);
        throw new Error(text || res.statusText);
      }
      // redirect on success
      //
      router.push("/");
    } catch {
      // TODO: handle specific error cases
      console.log("Error");
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
