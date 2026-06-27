import { useState } from "react";
import { Link, useNavigate, useParams } from "@tanstack/react-router";
import { useProduct } from "../hooks/useProducts";
import { useAddItem } from "../hooks/useCart";
import { useAuth } from "../auth/useAuth";
import { centsToUSD } from "../lib/money";
import { Button, Card, ErrorState } from "../components/ui";
import type { Variant } from "../types";

function QtyStepper({
  value,
  onChange,
}: {
  value: number;
  onChange: (n: number) => void;
}) {
  return (
    <div className="inline-flex items-center rounded-[2px] border border-line">
      <button
        type="button"
        aria-label="Decrease quantity"
        className="px-3 py-1 text-lg leading-none hover:bg-section disabled:opacity-30"
        disabled={value <= 1}
        onClick={() => onChange(Math.max(1, value - 1))}
      >
        −
      </button>
      <span className="data w-10 text-center text-sm">{value}</span>
      <button
        type="button"
        aria-label="Increase quantity"
        className="px-3 py-1 text-lg leading-none hover:bg-section"
        onClick={() => onChange(value + 1)}
      >
        +
      </button>
    </div>
  );
}

export function ProductPage() {
  const { id } = useParams({ from: "/products/$id" });
  const numericId = Number(id);
  const navigate = useNavigate();
  const { isAuthed } = useAuth();

  const detail = useProduct(numericId);
  const addItem = useAddItem();

  const [selected, setSelected] = useState<number | null>(null);
  const [qty, setQty] = useState(1);
  const [feedback, setFeedback] = useState<string | null>(null);

  if (detail.isLoading) {
    return <div className="h-64 animate-pulse rounded-[2px] bg-section" />;
  }
  if (detail.isError) {
    return <ErrorState message={(detail.error as Error).message} />;
  }
  if (!detail.data) return null;

  const { product, variants, options, images } = detail.data;
  const selectedVariant = variants.find((v) => v.id === selected) ?? null;

  async function handleAdd() {
    if (!selectedVariant) return;
    if (!isAuthed) {
      navigate({ to: "/login", search: { redirect: `/products/${id}` } });
      return;
    }
    setFeedback(null);
    try {
      await addItem.mutateAsync({
        variantId: selectedVariant.id,
        quantity: qty,
      });
      setFeedback("Added to cart.");
    } catch (e) {
      setFeedback((e as Error).message);
    }
  }

  return (
    <div className="grid grid-cols-1 gap-10 lg:grid-cols-2">
      <div>
        <Link to="/" className="text-sm text-accent hover:text-accent-ink">
          ← Catalog
        </Link>
        <h1 className="mt-4 text-3xl">{product.name}</h1>
        <p className="data mt-2 text-xs text-muted">
          {product.slug ?? `id ${product.id}`}
        </p>
        {product.description ? (
          <p className="mt-6 max-w-prose text-sm leading-relaxed text-ink">
            {product.description}
          </p>
        ) : null}

        {images.length > 0 ? (
          <div className="mt-6 grid grid-cols-2 gap-3">
            {images.map((img) => (
              <img
                key={img.id}
                src={img.url}
                alt={product.name}
                className="aspect-square w-full rounded-[2px] border border-line object-cover"
              />
            ))}
          </div>
        ) : null}

        {options.length > 0 ? (
          <div className="mt-8">
            <h3 className="text-xs uppercase tracking-wide text-muted">
              Options
            </h3>
            <ul className="mt-2 flex flex-wrap gap-2">
              {options.map((o) => (
                <li
                  key={o.id}
                  className="data rounded-[2px] border border-line px-2 py-1 text-xs"
                >
                  {o.name}
                </li>
              ))}
            </ul>
          </div>
        ) : null}
      </div>

      {/* The ticket panel: variants as a selectable SKU list. */}
      <Card className="self-start">
        <div className="border-b border-line px-5 py-4">
          <h2 className="font-display text-sm uppercase tracking-wide">
            Select a SKU
          </h2>
        </div>

        {variants.length === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-muted">
            No variants available for this product.
          </p>
        ) : (
          <ul>
            {variants.map((v: Variant) => {
              const out = v.stock <= 0;
              const active = v.id === selected;
              return (
                <li key={v.id} className="border-b border-line last:border-0">
                  <button
                    type="button"
                    disabled={out}
                    onClick={() => {
                      setSelected(v.id);
                      setQty(1);
                      setFeedback(null);
                    }}
                    className={`flex w-full items-center justify-between gap-3 px-5 py-3 text-left transition-colors disabled:opacity-50 ${
                      active ? "bg-section" : "hover:bg-section"
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <span
                        aria-hidden
                        className={`inline-block h-3 w-3 rounded-full border ${
                          active
                            ? "border-accent bg-accent"
                            : "border-line bg-surface"
                        }`}
                      />
                      <span className="data text-sm">{v.sku}</span>
                    </span>
                    <span className="flex items-center gap-4">
                      <span className="data text-sm">
                        {centsToUSD(v.price_in_cents)}
                      </span>
                      <span
                        className={`data text-xs ${
                          out ? "text-danger" : "text-stock"
                        }`}
                      >
                        {out ? "OUT" : `${v.stock} in stock`}
                      </span>
                    </span>
                  </button>
                </li>
              );
            })}
          </ul>
        )}

        <div className="border-t border-line px-5 py-4">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted">Quantity</span>
            <QtyStepper value={qty} onChange={setQty} />
          </div>

          {selectedVariant ? (
            <div className="mt-3 flex items-center justify-between text-sm">
              <span className="text-muted">Line total</span>
              <span className="data">
                {centsToUSD(selectedVariant.price_in_cents * qty)}
              </span>
            </div>
          ) : null}

          <Button
            className="mt-4 w-full"
            disabled={!selectedVariant || addItem.isPending}
            onClick={handleAdd}
          >
            {addItem.isPending
              ? "Adding…"
              : !isAuthed
                ? "Log in to add to cart"
                : "Add to cart"}
          </Button>

          {feedback ? (
            <p
              className={`mt-3 text-sm ${
                feedback === "Added to cart."
                  ? "text-stock"
                  : "text-danger"
              }`}
            >
              {feedback}
            </p>
          ) : null}
        </div>
      </Card>
    </div>
  );
}
