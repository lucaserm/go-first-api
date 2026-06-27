import { useState } from "react";
import { Link, useNavigate, useSearch } from "@tanstack/react-router";
import { useLogin } from "../hooks/useAuthMutations";
import { Button, Card, Field } from "../components/ui";

export function LoginPage() {
  const navigate = useNavigate();
  const { redirect } = useSearch({ from: "/login" });
  const login = useLogin();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      await login.mutateAsync({ email, password });
      navigate({ to: redirect ?? "/" });
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card className="p-6">
        <div className="mb-6 border-b border-line pb-4">
          <h1 className="text-xl">Log in</h1>
          <p className="data mt-1 text-xs text-muted">SESSION · AUTH</p>
        </div>

        <form onSubmit={onSubmit} className="space-y-4">
          <Field
            label="Email"
            type="email"
            name="email"
            autoComplete="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <Field
            label="Password"
            type="password"
            name="password"
            autoComplete="current-password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          {error ? (
            <p className="rounded-[2px] border border-line bg-section px-3 py-2 text-sm text-danger">
              {error}
            </p>
          ) : null}

          <Button type="submit" className="w-full" disabled={login.isPending}>
            {login.isPending ? "Signing in…" : "Sign in"}
          </Button>
        </form>

        <p className="mt-4 text-center text-sm text-muted">
          No account?{" "}
          <Link to="/register" className="text-accent hover:text-accent-ink">
            Register
          </Link>
        </p>
      </Card>
    </div>
  );
}
