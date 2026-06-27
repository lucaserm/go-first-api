import { useMutation } from "@tanstack/react-query";
import { apiFetch } from "../lib/api";
import { authStore } from "../auth/store";
import type { LoginResponse, User } from "../types";

export function useLogin() {
  return useMutation({
    mutationFn: (vars: { email: string; password: string }) =>
      apiFetch<LoginResponse>("/login", { method: "POST", body: vars }),
    onSuccess: (data) => {
      authStore.setSession(data.user, data.accessToken);
    },
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (vars: {
      username: string;
      email: string;
      password: string;
    }) => apiFetch<User>("/register", { method: "POST", body: vars }),
  });
}
