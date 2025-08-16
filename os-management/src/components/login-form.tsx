import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Lock } from "lucide-react";
import React, { useState } from "react";

type DivProps = Omit<React.HTMLAttributes<HTMLDivElement>, "onSubmit">;

interface LoginFormProps extends DivProps {
  onSubmit: (password: string) => void;
  loading?: boolean;
  error?: string;
}

export function LoginForm({
  onSubmit,
  loading = false,
  error,
  className,
  ...props
}: LoginFormProps) {
  const [password, setPassword] = useState("");

  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    onSubmit(password);
  };

  return (
    <div className={cn("flex max-w-96 flex-col gap-6", className)} {...props}>
      <Card className="space-y-8 py-8">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 w-fit rounded-full bg-zinc-800 p-3">
            <Lock className="h-6 w-6 text-zinc-300" />
          </div>
          <CardTitle>Admin Access</CardTitle>
          <CardDescription>
            Enter your admin password to access the management control panel
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="password">Admin Password</Label>
              <Input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button
              type="submit"
              className="w-full"
              variant={"secondary"}
              disabled={loading}
            >
              {loading ? "Logging in…" : "Login"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
