import { Link } from "@tanstack/react-router";
import { useProducts, useCategories } from "../hooks/useProducts";
import { Card, EmptyState, ErrorState } from "../components/ui";
import type { Product } from "../types";

function CategoryName(categoryId: number | null, lookup: Map<number, string>) {
  if (categoryId == null) return "Uncategorized";
  return lookup.get(categoryId) ?? `Category #${categoryId}`;
}

function ProductCard({
  product,
  categoryName,
}: {
  product: Product;
  categoryName: string;
}) {
  return (
    <Card className="flex flex-col justify-between p-5 transition-colors hover:border-ink">
      <div>
        <h2 className="font-display text-lg leading-snug">{product.name}</h2>
        <p className="data mt-2 text-xs text-muted">
          {categoryName}
          {product.slug ? `  ·  ${product.slug}` : ""}
        </p>
      </div>
      <div className="mt-6 flex items-center justify-between border-t border-line pt-4">
        <span className="data text-xs uppercase tracking-wide text-stock">
          {product.status}
        </span>
        <Link
          to="/products/$id"
          params={{ id: String(product.id) }}
          className="text-sm font-medium text-accent hover:text-accent-ink"
        >
          View →
        </Link>
      </div>
    </Card>
  );
}

function Skeletons() {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <div
          key={i}
          className="h-40 animate-pulse rounded-[2px] border border-line bg-section"
        />
      ))}
    </div>
  );
}

export function CatalogPage() {
  const products = useProducts();
  const categories = useCategories();

  const lookup = new Map<number, string>(
    (categories.data ?? []).map((c) => [c.id, c.name]),
  );

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl">Catalog</h1>
        <p className="mt-1 text-sm text-muted">
          Browse the catalog. Prices and SKUs live on each product's ticket.
        </p>
      </div>

      {products.isLoading ? (
        <Skeletons />
      ) : products.isError ? (
        <ErrorState message={(products.error as Error).message} />
      ) : (products.data ?? []).length === 0 ? (
        <EmptyState title="No products yet.">
          The catalog is empty. Check back soon.
        </EmptyState>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {products.data!.map((p) => (
            <ProductCard
              key={p.id}
              product={p}
              categoryName={CategoryName(p.category_id, lookup)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
