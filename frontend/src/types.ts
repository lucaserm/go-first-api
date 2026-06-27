// Shapes mirror the Go API. Every successful response is wrapped in
// `{ data: <payload> }`; the api client unwraps `.data` before these are seen.

export interface User {
  id: string;
  username: string;
  email: string;
}

export interface LoginResponse {
  user: User;
  accessToken: string;
}

export interface Product {
  id: number;
  name: string;
  slug: string | null;
  description: string;
  status: string;
  category_id: number | null;
  created_at: string;
}

export interface Variant {
  id: number;
  product_id: number;
  sku: string;
  price_in_cents: number;
  stock: number;
  weight_grams: number;
  created_at: string;
}

export interface ProductOption {
  id: number;
  product_id: number;
  name: string;
  position: number;
}

export interface ProductImage {
  id: number;
  product_id: number;
  variant_id: number | null;
  url: string;
  position: number;
  created_at: string;
}

export interface ProductDetail {
  product: Product;
  variants: Variant[];
  options: ProductOption[];
  images: ProductImage[];
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id: number | null;
  created_at: string;
}

export interface CartLineItem {
  variantId: number;
  sku: string;
  productName: string;
  quantity: number;
  unitPriceInCents: number;
  lineTotalInCents: number;
}

export interface CartResponse {
  items: CartLineItem[];
  subtotalInCents: number;
  itemCount: number;
}

// The address model is serialized straight from the DB row (snake_case).
export interface Address {
  id: number;
  recipient_name: string;
  line1: string;
  line2: string;
  city: string;
  region: string;
  postal_code: string;
  country: string;
  phone: string;
  is_default: boolean;
  created_at?: string;
}

// POST /addresses uses a camelCase payload (see addresses/types.go).
export interface CreateAddressPayload {
  recipientName: string;
  line1: string;
  line2: string;
  city: string;
  region: string;
  postalCode: string;
  country: string;
  phone: string;
  isDefault: boolean;
}

export interface ShippingAddress {
  recipientName: string;
  line1: string;
  line2: string;
  city: string;
  region: string;
  postalCode: string;
  country: string;
  phone: string;
}

export interface OrderLineItem {
  variantId: number;
  sku: string;
  productName: string;
  quantity: number;
  unitPriceInCents: number;
  lineTotalInCents: number;
}

export interface OrderResponse {
  id: number;
  status: string;
  currency: string;
  subtotalCents: number;
  shippingCents: number;
  taxCents: number;
  totalCents: number;
  shipping: ShippingAddress;
  items: OrderLineItem[];
}
