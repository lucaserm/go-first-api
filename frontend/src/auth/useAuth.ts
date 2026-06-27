import { useSyncExternalStore } from "react";
import { authStore } from "./store";

/** Reactive view of the auth session for components. */
export function useAuth() {
  const state = useSyncExternalStore(
    authStore.subscribe,
    authStore.getSnapshot,
    authStore.getSnapshot,
  );

  return {
    user: state.user,
    token: state.token,
    isAuthed: state.token != null,
    logout: () => authStore.logout(),
  };
}
