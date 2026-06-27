import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "../lib/api";
import type { Category, Product, ProductDetail } from "../types";

export function useProducts() {
  return useQuery({
    queryKey: ["products"],
    queryFn: () => apiFetch<{ products: Product[] }>("/products"),
    select: (d) => d.products ?? [],
  });
}

export function useProduct(id: number) {
  return useQuery({
    queryKey: ["product", id],
    queryFn: () => apiFetch<ProductDetail>(`/products/${id}`),
    enabled: Number.isFinite(id),
  });
}

export function useCategories() {
  return useQuery({
    queryKey: ["categories"],
    queryFn: () => apiFetch<{ categories: Category[] }>("/categories"),
    select: (d) => d.categories ?? [],
  });
}
