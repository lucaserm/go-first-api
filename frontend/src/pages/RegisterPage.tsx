import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { useLogin, useRegister } from "../hooks/useAuthMutations";
import { Button, Card, Field } from "../components/ui";

export function RegisterPage() {
  const navigate = useNavigate();
  const register = useRegister();
  const login = useLogin();

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      await register.mutateAsync({ username, email, password });
      // Register returns the user but no token; log in to establish a session.
      await login.mutateAsync({ email, password });
      navigate({ to: "/" });
    } catch (err) {
      setError((err as Error).message);
    }
  }

  const pending = register.isPending || login.isPending;

  return (
    <div className="mx-auto max-w-md">
      <Card className="p-6">
        <div className="mb-6 border-b border-line pb-4">
          <h1 className="text-xl">Create account</h1>
          <p className="data mt-1 text-xs text-muted">NEW · CUSTOMER</p>
        </div>

        <form onSubmit={onSubmit} className="space-y-4">
          <Field
            label="Username"
            name="username"
            autoComplete="username"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
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
            autoComplete="new-password"
            minLength={3}
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          {error ? (
            <p className="rounded-[2px] border border-line bg-section px-3 py-2 text-sm text-danger">
              {error}
            </p>
          ) : null}

          <Button type="submit" className="w-full" disabled={pending}>
            {pending ? "Creating…" : "Create account"}
          </Button>
        </form>

        <p className="mt-4 text-center text-sm text-muted">
          Already registered?{" "}
          <Link to="/login" className="text-accent hover:text-accent-ink">
            Log in
          </Link>
        </p>
      </Card>
    </div>
  );
}
