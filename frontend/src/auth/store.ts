import type { User } from "../types";

// A tiny framework-agnostic auth store. It lives outside React so the fetch
// client and the router's `beforeLoad` guards can read the token synchronously,
// while React components subscribe via `useSyncExternalStore` (see useAuth).

const TOKEN_KEY = "ecom.token";
const USER_KEY = "ecom.user";

interface AuthState {
  token: string | null;
  user: User | null;
}

function readInitial(): AuthState {
  try {
    const token = localStorage.getItem(TOKEN_KEY);
    const rawUser = localStorage.getItem(USER_KEY);
    return {
      token: token,
      user: rawUser ? (JSON.parse(rawUser) as User) : null,
    };
  } catch {
    return { token: null, user: null };
  }
}

let state: AuthState = readInitial();
const listeners = new Set<() => void>();
let unauthorizedHandler: (() => void) | null = null;

function emit() {
  for (const l of listeners) l();
}

export const authStore = {
  subscribe(listener: () => void): () => void {
    listeners.add(listener);
    return () => {
      listeners.delete(listener);
    };
  },

  getSnapshot(): AuthState {
    return state;
  },

  getToken(): string | null {
    return state.token;
  },

  isAuthed(): boolean {
    return state.token != null;
  },

  setSession(user: User, token: string) {
    state = { user, token };
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    emit();
  },

  logout() {
    state = { token: null, user: null };
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    emit();
  },

  // Registered once at startup (main.tsx) so a 401 from the api client can
  // route the user back to /login through the real router.
  setUnauthorizedHandler(fn: () => void) {
    unauthorizedHandler = fn;
  },

  handleUnauthorized() {
    this.logout();
    unauthorizedHandler?.();
  },
};
