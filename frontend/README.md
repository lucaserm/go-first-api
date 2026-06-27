# go-ecom — Frontend

A React storefront for the Go ecommerce API. Part of the monorepo; the API lives
in [`../backend`](../backend).

Design identity: **"the catalog as a ticket."** Data — prices, SKUs, stock,
order numbers, quantities — is the aesthetic, rendered in monospace with tabular
alignment and hairline dividers so cart and order views read like a precise
receipt / packing slip.

## Stack

- Vite + React 18 + TypeScript
- [TanStack Router](https://tanstack.com/router) (code-based routing)
- [TanStack Query](https://tanstack.com/query) for server state
- Tailwind CSS v4 via `@tailwindcss/vite` (CSS-first `@theme` config — no `tailwind.config.js`)
- Fonts: Space Grotesk, Inter, JetBrains Mono (Google Fonts)
- Package manager / runtime: **bun**

## Getting started

```bash
bun install
bun run dev      # starts Vite on http://localhost:5173
```

Other scripts:

```bash
bun run typecheck   # tsc --noEmit
bun run build       # tsc --noEmit && vite build
bun run preview     # preview the production build
```

## Configuration

The API base URL is read from `VITE_API_URL`. Copy the example and adjust if
needed:

```bash
cp .env.example .env
```

```
VITE_API_URL=http://localhost:8080
```

Every successful API response is wrapped in `{ "data": <payload> }` and errors
are `{ "error": "message" }` with a non-2xx status. The typed fetch client in
`src/lib/api.ts` unwraps `.data`, attaches `Authorization: Bearer <token>` when
logged in, and throws an `ApiError(message)` on failure.

## Running against the backend

Start the Go API from the sibling directory:

```bash
cd ../backend
# follow the backend README / Makefile to run migrations and start the server
```

The API listens on `:8080` by default, matching `VITE_API_URL`.

## Project structure

```
src/
  lib/        api.ts (fetch client), money.ts (centsToUSD)
  auth/       store.ts (token/user store), useAuth.ts
  hooks/      useProducts, useCart, useAuthMutations, useOrders
  components/ RootLayout, Nav, ui (Button/Field/Card/EmptyState)
  pages/      CatalogPage, ProductPage, LoginPage, RegisterPage, CartPage
  router.tsx  code-based TanStack Router route tree + /cart auth guard
  main.tsx    QueryClient + RouterProvider
```

## Routes

- `/` — catalog grid (product cards; price lives on the detail ticket)
- `/products/:id` — product detail with a selectable variant/SKU ticket panel
- `/login`, `/register` — receipt-styled auth forms
- `/cart` — auth-guarded receipt; checkout picks/adds an address then places an order
