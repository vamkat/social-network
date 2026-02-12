# Server Actions vs API Routes

> **Prerequisites:** [React Fundamentals](./01-react-fundamentals.md), [Next.js Fundamentals](./02-nextjs-fundamentals.md), [Architecture](./03-architecture.md).
>
> This guide explains why SocialSphere uses server actions for data fetching and walks through the full login flow.

---

## Two Ways to Talk to a Backend in Next.js

### API Routes (traditional)

API routes are HTTP endpoints defined in `app/api/*/route.js`. Client code calls them with `fetch()`:

```jsx
// app/api/posts/route.js (server)
export async function GET(request) {
    const data = await fetchFromDatabase();
    return Response.json(data);
}

// Component (client) — you manage the URL and fetch yourself
const res = await fetch("/api/posts");
const data = await res.json();
```

### Server Actions (what this project uses)

Server actions are functions marked `"use server"` that you import and call like regular functions:

```jsx
// actions/posts/get-public-posts.js (server)
"use server";
export async function getPublicPosts({ limit, offset }) {
    const result = await serverApiRequest(`/posts?limit=${limit}&offset=${offset}`);
    return { success: result.ok, data: result.data };
}

// Component (client) — just import and call
import { getPublicPosts } from "@/actions/posts/get-public-posts";
const result = await getPublicPosts({ limit: 10, offset: 0 });
```

Behind the scenes, Next.js handles the network request. You never write a URL, parse a response, or manage headers.

---

## This Project's Data Flow

Every data request follows the same path:

```
UI Component → Server Action → serverApiRequest() → Go Backend
```

Let's trace the complete login flow to see how this works in practice.

### Step 1: The Login Form (client component)

<!-- src/components/forms/LoginForm.js -->
```jsx
"use client";

import { login } from "@/actions/auth/login";  // Import the server action
import { useStore } from "@/store/store";

export default function LoginForm() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const setUser = useStore((state) => state.setUser);

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);
        const email = formData.get("email");
        const password = formData.get("password");

        try {
            // Call the server action — looks like a normal function call
            // But it actually sends a POST request to the server
            const resp = await login({ email, password });

            if (!resp.success || resp.error) {
                setError(resp.error || "Invalid credentials");
                setIsLoading(false);
                return;
            }

            // Store user data in Zustand (persisted to localStorage)
            setUser({
                id: resp.user_id,
                username: resp.username,
                avatar_url: resp.avatar_url || ""
            });

            // Redirect to feed
            window.location.href = "/feed/public";
        } catch (error) {
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    return <form onSubmit={handleSubmit}>{/* ... */}</form>;
}
```

### Step 2: The Server Action

<!-- src/actions/auth/login.js -->
```jsx
"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function login(credentials) {
    try {
        const apiResp = await serverApiRequest("/login", {
            method: "POST",
            body: JSON.stringify(credentials),
            forwardCookies: true,    // Forward Set-Cookie headers from backend
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!apiResp.ok) {
            return { success: false, status: apiResp.status, error: apiResp.message };
        }

        if (!apiResp.data.id) {
            return { success: false, error: "Login failed - no user ID returned" };
        }

        return {
            success: true,
            user_id: apiResp.data.id,
            username: apiResp.data.username,
            avatar_url: apiResp.data.avatar_url || ""
        };
    } catch (error) {
        return { success: false, error: error.message };
    }
}
```

### Step 3: `serverApiRequest()` — The Chokepoint

<!-- src/lib/server-api.js -->
```jsx
"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { propagation, context } from "@opentelemetry/api";
import * as logger from "@/lib/logger.server";

const API_BASE = process.env.GATEWAY;  // e.g., "http://gateway:8080"

export async function serverApiRequest(endpoint, options = {}) {
    const method = options.method || "GET";
    const start = performance.now();

    // Log outgoing request
    logger.info("outgoing request @1 @2", "method", method, "url", endpoint);

    try {
        // 1. Read JWT from cookies
        const cookieStore = await cookies();
        const jwt = cookieStore.get("jwt")?.value;

        // 2. Build headers — inject auth cookie
        const headers = { ...(options.headers || {}) };
        if (jwt) headers["Cookie"] = `jwt=${jwt}`;

        // 3. Inject W3C trace context for distributed tracing
        propagation.inject(context.active(), headers);

        // 4. Make the actual HTTP request to the Go backend
        const res = await fetch(`${API_BASE}${endpoint}`, {
            ...options,
            headers,
            cache: "no-store"   // Never cache API responses
        });

        // 5. Forward cookies back (for login — the backend sets the JWT cookie)
        if (options.forwardCookies) {
            const setCookieHeaders = res.headers.getSetCookie
                ? res.headers.getSetCookie()
                : [];
            // ... parse and set each cookie on the response
        }

        // 6. Handle error responses
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            const duration = Math.round(performance.now() - start);
            const errMsg = err.error || err.message || "Unknown error";

            logger.error("request failed @1 @2 @3 @4 @5",
                "method", method, "url", endpoint, "status", res.status,
                "error", errMsg, "duration_ms", duration);

            // 401 → clear cookie and redirect to login
            if (res.status === 401) {
                cookieStore.delete("jwt");
                redirect("/login");
            }

            // 403, 400 → return error for the UI to handle
            if (res.status === 403) {
                return { ok: false, status: 403, message: errMsg };
            }
            if (res.status === 400) {
                return { ok: false, status: 400, message: errMsg };
            }

            return { ok: false, status: res.status, message: errMsg };
        }

        // 7. Parse successful response
        const duration = Math.round(performance.now() - start);
        const text = await res.text();

        logger.info("request succeeded @1 @2 @3 @4",
            "method", method, "url", endpoint, "status", res.status,
            "duration_ms", duration);

        if (!text || text.trim() === '') {
            return { ok: true, data: null };  // Empty response (e.g., DELETE)
        }
        return { ok: true, data: JSON.parse(text) };

    } catch (e) {
        // Re-throw Next.js redirects (they use a special error mechanism)
        if (e?.digest?.startsWith("NEXT_REDIRECT")) throw e;

        const duration = Math.round(performance.now() - start);
        logger.error("request exception @1 @2 @3 @4",
            "method", method, "url", endpoint,
            "error", e?.message || "unknown", "duration_ms", duration);

        return { ok: false, message: "Network error" };
    }
}
```

---

## Why Server Actions Were Chosen

| Concern | API Routes | Server Actions (this project) |
|---------|-----------|-------------------------------|
| **Calling from client** | `fetch("/api/posts")` — manage URLs yourself | `getPublicPosts()` — import and call |
| **Type safety** | None — URL strings, manual parsing | Function signatures, return types |
| **Auth handling** | Each route reads cookies separately | Centralized in `serverApiRequest()` |
| **Telemetry** | Add tracing to each route | Centralized in `serverApiRequest()` |
| **Error handling** | Each route handles errors | Centralized 401/403/400 handling |
| **Server components** | Need `fetch()` calls | Call directly with `await` |
| **Client components** | Need `fetch()` calls | Import and call like functions |

The key advantage: **one function** (`serverApiRequest`) handles auth, tracing, error handling, logging, and cookie management for the entire app. Every backend call gets these for free.

---

## When You'd Still Want API Routes

Server actions are great for this project, but API routes are better for:

- **Webhooks** — external services need a URL to call (e.g., payment notifications)
- **Public APIs** — if other applications need to call your backend
- **External service callbacks** — OAuth providers redirect to URLs, not function calls
- **File downloads** — streaming responses that need custom headers

This project doesn't have any of these use cases, so server actions are the right choice.

---

## Return Value Convention

Every server action in this project returns an object with a predictable shape:

```jsx
// Success
{ success: true, data: { ... } }

// Failure
{ success: false, error: "Error message" }
{ success: false, status: 403, error: "Forbidden" }
```

The underlying `serverApiRequest()` returns:

```jsx
// Success
{ ok: true, data: { ... } }

// Failure
{ ok: false, status: 400, message: "Bad request" }
```

Server actions translate between these:

```jsx
"use server";
export async function getPublicPosts({ limit, offset }) {
    const result = await serverApiRequest(`/posts?limit=${limit}&offset=${offset}`);
    if (!result.ok) {
        return { success: false, error: result.message };
    }
    return { success: true, data: result.data };
}
```

This consistent pattern means every component can handle responses the same way:

```jsx
const result = await someAction(params);
if (!result.success) {
    setError(result.error);
    return;
}
// Use result.data
```

---

## Next Steps

Continue to [State and Context](./05-state-and-context.md) to learn about Zustand stores, React Context, WebSocket connections, and how data flows through the app in real-time.
