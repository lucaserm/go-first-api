import { Outlet } from "@tanstack/react-router";
import { Nav } from "./Nav";

export function RootLayout() {
  return (
    <div className="min-h-screen bg-paper">
      <Nav />
      <main className="mx-auto max-w-5xl px-4 py-8">
        <Outlet />
      </main>
      <footer className="mx-auto max-w-5xl px-4 py-10">
        <p className="data border-t border-line pt-4 text-xs text-muted">
          THE CATALOG AS A TICKET · {new Date().getFullYear()}
        </p>
      </footer>
    </div>
  );
}
