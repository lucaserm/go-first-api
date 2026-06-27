import { Link } from "@tanstack/react-router";
import { useAuth } from "../auth/useAuth";
import { useCart } from "../hooks/useCart";

export function Nav() {
  const { isAuthed, user, logout } = useAuth();
  const cart = useCart();
  const count = cart.data?.itemCount ?? 0;

  return (
    <header className="border-b border-line bg-paper">
      <nav className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
        <Link to="/" className="font-display text-lg font-bold tracking-tight">
          TICKET<span className="text-accent">.</span>
        </Link>

        <div className="flex items-center gap-5 text-sm">
          <Link
            to="/cart"
            className="flex items-center gap-1.5 hover:text-accent-ink"
          >
            <span>Cart</span>
            <span className="data rounded-[2px] bg-section px-1.5 py-0.5 text-xs">
              {count}
            </span>
          </Link>

          {isAuthed ? (
            <div className="flex items-center gap-4">
              <span className="hidden text-muted sm:inline">
                {user?.username}
              </span>
              <button
                onClick={() => logout()}
                className="hover:text-accent-ink"
              >
                Log out
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-4">
              <Link to="/login" className="hover:text-accent-ink">
                Log in
              </Link>
              <Link
                to="/register"
                className="rounded-[2px] bg-accent px-3 py-1.5 text-white hover:bg-accent-ink"
              >
                Register
              </Link>
            </div>
          )}
        </div>
      </nav>
    </header>
  );
}
