"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { useAuthStore } from "@/stores/auth-store";
import { t } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: { id: string; email: string };
}

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const login = useAuthStore((s) => s.login);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const redirect = searchParams.get("redirect") || "/account/orders";

  const mutation = useMutation({
    mutationFn: async () => {
      const res = await api.post<LoginResponse>("/auth/login", {
        email,
        password,
      });
      return res.data!;
    },
    onSuccess: (data) => {
      login(data.access_token, data.refresh_token, data.user);
      router.push(redirect);
    },
    onError: (err: Error) => setError(err.message),
  });

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>{t("account.loginTitle")}</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            setError("");
            mutation.mutate();
          }}
        >
          <div className="space-y-2">
            <Label htmlFor="email">{t("account.email")}</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="ltr-field"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">{t("account.password")}</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="ltr-field"
              required
            />
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button
            type="submit"
            className="w-full"
            disabled={mutation.isPending}
          >
            {mutation.isPending ? t("account.signingIn") : t("account.signIn")}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

export default function LoginPage() {
  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <Suspense fallback={<p className="text-muted-foreground">{t("common.loading")}</p>}>
        <LoginForm />
      </Suspense>
    </div>
  );
}
