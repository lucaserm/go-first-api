import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "../lib/api";
import { useAuth } from "../auth/useAuth";
import type {
  Address,
  CreateAddressPayload,
  OrderResponse,
} from "../types";

export function useAddresses() {
  const { isAuthed } = useAuth();
  return useQuery({
    queryKey: ["addresses"],
    queryFn: () => apiFetch<{ addresses: Address[] }>("/addresses"),
    select: (d) => d.addresses ?? [],
    enabled: isAuthed,
  });
}

export function useCreateAddress() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateAddressPayload) =>
      apiFetch<Address>("/addresses", { method: "POST", body: payload }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["addresses"] }),
  });
}

export function usePlaceOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { addressId: number }) =>
      apiFetch<OrderResponse>("/orders", { method: "POST", body: vars }),
    onSuccess: () => {
      // The order consumes the cart server-side; refresh both views.
      qc.invalidateQueries({ queryKey: ["cart"] });
      qc.invalidateQueries({ queryKey: ["orders"] });
    },
  });
}
