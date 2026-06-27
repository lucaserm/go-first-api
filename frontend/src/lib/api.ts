import { authStore } from "../auth/store";

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  // Some auth'd routes return 204 No Content (cart delete, address delete).
  // The caller can opt out of expecting a JSON `{ data }` envelope.
  expectNoContent?: boolean;
}

/**
 * Typed fetch wrapper.
 * - Prefixes BASE_URL.
 * - Attaches `Authorization: Bearer <token>` when logged in.
 * - Unwraps the `{ data: <payload> }` envelope on success.
 * - Throws ApiError(message) on non-2xx, parsing `{ error }` for the message.
 * - On 401, clears auth and triggers the registered redirect.
 */
export async function apiFetch<T>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {};
  const token = authStore.getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  let body: string | undefined;
  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
    body = JSON.stringify(options.body);
  }

  let res: Response;
  try {
    res = await fetch(`${BASE_URL}${path}`, {
      method: options.method ?? "GET",
      headers,
      body,
    });
  } catch {
    throw new ApiError("Network error — is the API running?", 0);
  }

  if (res.status === 401) {
    authStore.handleUnauthorized();
    throw new ApiError("Your session has expired. Please sign in again.", 401);
  }

  if (!res.ok) {
    let message = `Request failed (${res.status})`;
    try {
      const parsed = (await res.json()) as { error?: string };
      if (parsed?.error) message = parsed.error;
    } catch {
      /* non-JSON error body; keep the default message */
    }
    throw new ApiError(message, res.status);
  }

  if (res.status === 204 || options.expectNoContent) {
    return undefined as T;
  }

  const json = (await res.json()) as { data: T };
  return json.data;
}
