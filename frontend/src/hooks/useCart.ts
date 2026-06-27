import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "../lib/api";
import { useAuth } from "../auth/useAuth";
import type { CartResponse } from "../types";

const CART_KEY = ["cart"];

export function useCart() {
  const { isAuthed } = useAuth();
  return useQuery({
    queryKey: CART_KEY,
    queryFn: () => apiFetch<CartResponse>("/cart"),
    enabled: isAuthed,
  });
}

export function useAddItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { variantId: number; quantity: number }) =>
      apiFetch<CartResponse>("/cart/items", { method: "POST", body: vars }),
    onSuccess: (data) => {
      qc.setQueryData(CART_KEY, data);
      qc.invalidateQueries({ queryKey: CART_KEY });
    },
  });
}

export function useUpdateItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { variantId: number; quantity: number }) =>
      apiFetch<CartResponse>(`/cart/items/${vars.variantId}`, {
        method: "PATCH",
        body: { quantity: vars.quantity },
      }),
    onSuccess: (data) => {
      qc.setQueryData(CART_KEY, data);
      qc.invalidateQueries({ queryKey: CART_KEY });
    },
  });
}

export function useRemoveItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (variantId: number) =>
      apiFetch<void>(`/cart/items/${variantId}`, {
        method: "DELETE",
        expectNoContent: true,
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: CART_KEY }),
  });
}

export function useClearCart() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () =>
      apiFetch<void>("/cart", { method: "DELETE", expectNoContent: true }),
    onSuccess: () => qc.invalidateQueries({ queryKey: CART_KEY }),
  });
}
