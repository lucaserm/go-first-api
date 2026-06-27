import { useState } from "react";
import { Link } from "@tanstack/react-router";
import {
  useCart,
  useUpdateItem,
  useRemoveItem,
} from "../hooks/useCart";
import {
  useAddresses,
  useCreateAddress,
  usePlaceOrder,
} from "../hooks/useOrders";
import { centsToUSD } from "../lib/money";
import { Button, Card, EmptyState, ErrorState, Field } from "../components/ui";
import type {
  Address,
  CartLineItem,
  CreateAddressPayload,
  OrderResponse,
} from "../types";

function Row({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between border-b border-line py-3 last:border-0">
      {children}
    </div>
  );
}

function LineItem({ item }: { item: CartLineItem }) {
  const update = useUpdateItem();
  const remove = useRemoveItem();

  return (
    <div className="grid grid-cols-[1fr_auto] gap-3 border-b border-line py-4 last:border-0">
      <div>
        <p className="font-display text-sm">{item.productName}</p>
        <p className="data mt-0.5 text-xs text-muted">{item.sku}</p>
        <div className="mt-2 inline-flex items-center rounded-[2px] border border-line">
          <button
            aria-label="Decrease quantity"
            className="px-2 py-0.5 text-base leading-none hover:bg-section disabled:opacity-30"
            disabled={item.quantity <= 1 || update.isPending}
            onClick={() =>
              update.mutate({
                variantId: item.variantId,
                quantity: item.quantity - 1,
              })
            }
          >
            −
          </button>
          <span className="data w-8 text-center text-xs">{item.quantity}</span>
          <button
            aria-label="Increase quantity"
            className="px-2 py-0.5 text-base leading-none hover:bg-section disabled:opacity-30"
            disabled={update.isPending}
            onClick={() =>
              update.mutate({
                variantId: item.variantId,
                quantity: item.quantity + 1,
              })
            }
          >
            +
          </button>
        </div>
      </div>

      <div className="text-right">
        <p className="data text-sm">{centsToUSD(item.lineTotalInCents)}</p>
        <p className="data mt-0.5 text-xs text-muted">
          {centsToUSD(item.unitPriceInCents)} ea
        </p>
        <button
          className="mt-2 text-xs text-danger hover:underline disabled:opacity-40"
          disabled={remove.isPending}
          onClick={() => remove.mutate(item.variantId)}
        >
          Remove
        </button>
      </div>
    </div>
  );
}

const emptyAddress: CreateAddressPayload = {
  recipientName: "",
  line1: "",
  line2: "",
  city: "",
  region: "",
  postalCode: "",
  country: "US",
  phone: "",
  isDefault: false,
};

function AddressForm({
  onCreated,
}: {
  onCreated: (a: Address) => void;
}) {
  const create = useCreateAddress();
  const [form, setForm] = useState<CreateAddressPayload>(emptyAddress);
  const [error, setError] = useState<string | null>(null);

  function set<K extends keyof CreateAddressPayload>(
    key: K,
    value: CreateAddressPayload[K],
  ) {
    setForm((f) => ({ ...f, [key]: value }));
  }

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const created = await create.mutateAsync(form);
      onCreated(created);
      setForm(emptyAddress);
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <Field
        label="Recipient name"
        required
        value={form.recipientName}
        onChange={(e) => set("recipientName", e.target.value)}
      />
      <Field
        label="Address line 1"
        required
        value={form.line1}
        onChange={(e) => set("line1", e.target.value)}
      />
      <Field
        label="Address line 2"
        value={form.line2}
        onChange={(e) => set("line2", e.target.value)}
      />
      <div className="grid grid-cols-2 gap-3">
        <Field
          label="City"
          required
          value={form.city}
          onChange={(e) => set("city", e.target.value)}
        />
        <Field
          label="Region / State"
          required
          value={form.region}
          onChange={(e) => set("region", e.target.value)}
        />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <Field
          label="Postal code"
          required
          value={form.postalCode}
          onChange={(e) => set("postalCode", e.target.value)}
        />
        <Field
          label="Country (2-letter)"
          required
          maxLength={2}
          value={form.country}
          onChange={(e) => set("country", e.target.value.toUpperCase())}
        />
      </div>
      <Field
        label="Phone"
        value={form.phone}
        onChange={(e) => set("phone", e.target.value)}
      />

      {error ? <p className="text-sm text-danger">{error}</p> : null}

      <Button type="submit" variant="secondary" disabled={create.isPending}>
        {create.isPending ? "Saving…" : "Save address"}
      </Button>
    </form>
  );
}

function OrderReceipt({ order }: { order: OrderResponse }) {
  return (
    <Card className="p-6">
      <div className="flex items-center justify-between border-b border-line pb-4">
        <h1 className="text-xl">Order placed</h1>
        <span className="data text-sm">#{order.id}</span>
      </div>

      <div className="mt-4 flex items-center justify-between text-sm">
        <span className="text-muted">Status</span>
        <span className="data uppercase">{order.status}</span>
      </div>

      <div className="mt-6">
        <h3 className="text-xs uppercase tracking-wide text-muted">Items</h3>
        <div className="mt-2">
          {order.items.map((it) => (
            <Row key={it.variantId}>
              <div>
                <p className="font-display text-sm">{it.productName}</p>
                <p className="data text-xs text-muted">
                  {it.sku} · {it.quantity} × {centsToUSD(it.unitPriceInCents)}
                </p>
              </div>
              <span className="data text-sm">
                {centsToUSD(it.lineTotalInCents)}
              </span>
            </Row>
          ))}
        </div>
      </div>

      <div className="mt-6 space-y-1 border-t border-line pt-4 text-sm">
        <div className="flex justify-between">
          <span className="text-muted">Subtotal</span>
          <span className="data">{centsToUSD(order.subtotalCents)}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-muted">Shipping</span>
          <span className="data">{centsToUSD(order.shippingCents)}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-muted">Tax</span>
          <span className="data">{centsToUSD(order.taxCents)}</span>
        </div>
        <div className="flex justify-between border-t border-line pt-2 font-medium">
          <span>Total ({order.currency})</span>
          <span className="data">{centsToUSD(order.totalCents)}</span>
        </div>
      </div>

      <div className="mt-6 border-t border-line pt-4">
        <h3 className="text-xs uppercase tracking-wide text-muted">
          Ship to
        </h3>
        <div className="data mt-2 text-xs leading-relaxed">
          <p>{order.shipping.recipientName}</p>
          <p>{order.shipping.line1}</p>
          {order.shipping.line2 ? <p>{order.shipping.line2}</p> : null}
          <p>
            {order.shipping.city}, {order.shipping.region}{" "}
            {order.shipping.postalCode}
          </p>
          <p>{order.shipping.country}</p>
        </div>
      </div>

      <Link
        to="/"
        className="mt-6 inline-block text-sm text-accent hover:text-accent-ink"
      >
        ← Continue shopping
      </Link>
    </Card>
  );
}

export function CartPage() {
  const cart = useCart();
  const addresses = useAddresses();
  const placeOrder = usePlaceOrder();

  const [checkingOut, setCheckingOut] = useState(false);
  const [selectedAddress, setSelectedAddress] = useState<number | null>(null);
  const [showAddrForm, setShowAddrForm] = useState(false);
  const [placedOrder, setPlacedOrder] = useState<OrderResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (placedOrder) {
    return (
      <div className="mx-auto max-w-xl">
        <OrderReceipt order={placedOrder} />
      </div>
    );
  }

  if (cart.isLoading) {
    return <div className="h-64 animate-pulse rounded-[2px] bg-section" />;
  }
  if (cart.isError) {
    return <ErrorState message={(cart.error as Error).message} />;
  }

  const items = cart.data?.items ?? [];

  if (items.length === 0) {
    return (
      <div className="mx-auto max-w-xl">
        <EmptyState title="Your cart is empty.">
          <Link to="/" className="text-accent hover:text-accent-ink">
            Browse the catalog.
          </Link>
        </EmptyState>
      </div>
    );
  }

  async function handlePlaceOrder() {
    if (selectedAddress == null) {
      setError("Select a shipping address.");
      return;
    }
    setError(null);
    try {
      const order = await placeOrder.mutateAsync({ addressId: selectedAddress });
      setPlacedOrder(order);
    } catch (err) {
      setError((err as Error).message);
    }
  }

  const addressList = addresses.data ?? [];

  return (
    <div className="mx-auto grid max-w-4xl grid-cols-1 gap-8 lg:grid-cols-[1fr_22rem]">
      <div>
        <h1 className="mb-4 text-2xl">Cart</h1>
        <Card className="px-5">
          {items.map((it) => (
            <LineItem key={it.variantId} item={it} />
          ))}
        </Card>
      </div>

      <div>
        <Card className="p-5">
          <h2 className="font-display text-sm uppercase tracking-wide">
            Summary
          </h2>
          <div className="mt-4 flex justify-between border-t border-line pt-4 text-sm">
            <span className="text-muted">Subtotal</span>
            <span className="data">
              {centsToUSD(cart.data?.subtotalInCents ?? 0)}
            </span>
          </div>
          <div className="mt-1 flex justify-between text-sm">
            <span className="text-muted">Items</span>
            <span className="data">{cart.data?.itemCount ?? 0}</span>
          </div>

          {!checkingOut ? (
            <Button
              className="mt-5 w-full"
              onClick={() => setCheckingOut(true)}
            >
              Checkout
            </Button>
          ) : (
            <div className="mt-5 border-t border-line pt-4">
              <h3 className="text-xs uppercase tracking-wide text-muted">
                Shipping address
              </h3>

              {addresses.isLoading ? (
                <p className="mt-2 text-sm text-muted">Loading addresses…</p>
              ) : addressList.length === 0 && !showAddrForm ? (
                <p className="mt-2 text-sm text-muted">
                  No saved addresses.
                </p>
              ) : (
                <ul className="mt-2 space-y-2">
                  {addressList.map((a) => (
                    <li key={a.id}>
                      <label className="flex cursor-pointer items-start gap-2 rounded-[2px] border border-line p-2 text-sm has-[:checked]:border-accent">
                        <input
                          type="radio"
                          name="address"
                          className="mt-1"
                          checked={selectedAddress === a.id}
                          onChange={() => setSelectedAddress(a.id)}
                        />
                        <span className="data text-xs leading-relaxed">
                          {a.recipient_name}
                          <br />
                          {a.line1}, {a.city} {a.postal_code}
                          <br />
                          {a.region}, {a.country}
                        </span>
                      </label>
                    </li>
                  ))}
                </ul>
              )}

              {showAddrForm ? (
                <div className="mt-4 border-t border-line pt-4">
                  <AddressForm
                    onCreated={(a) => {
                      setSelectedAddress(a.id);
                      setShowAddrForm(false);
                    }}
                  />
                </div>
              ) : (
                <button
                  className="mt-3 text-sm text-accent hover:text-accent-ink"
                  onClick={() => setShowAddrForm(true)}
                >
                  + Add a new address
                </button>
              )}

              {error ? (
                <p className="mt-3 text-sm text-danger">{error}</p>
              ) : null}

              <Button
                className="mt-4 w-full"
                disabled={placeOrder.isPending || selectedAddress == null}
                onClick={handlePlaceOrder}
              >
                {placeOrder.isPending ? "Placing order…" : "Place order"}
              </Button>
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}
