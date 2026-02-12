# Next.js Fundamentals

> **Prerequisites:** [React Fundamentals](./01-react-fundamentals.md) — components, props, state, effects, JSX.
>
> This guide explains what Next.js adds on top of React, using real code from SocialSphere.

---

## What Next.js Adds

React alone is just a UI library — it doesn't handle routing, server-side rendering, or data fetching patterns. Next.js is a **framework** built on React that provides:

- **File-based routing** — folders and files in `src/app/` become URL routes
- **Server components** — components that run on the server (can access databases, APIs, etc.)
- **Client components** — components that run in the browser (can use state, effects, etc.)
- **Server actions** — server-side functions callable directly from client code
- **Optimized images, fonts, and metadata**

---

## App Router File Conventions

Next.js uses special filenames inside `src/app/` to define pages, layouts, loading states, etc.

| File | Purpose |
|------|---------|
| `page.js` | The UI for a route (required to make a URL accessible) |
| `layout.js` | Shared wrapper that persists across child routes |
| `loading.js` | Loading UI shown while the page is loading |
| `error.js` | Error boundary UI shown when something goes wrong |

### Root Layout

Every Next.js app has a root layout that wraps all pages:

<!-- src/app/layout.js -->
```jsx
import "./globals.css";
import { ThemeProvider } from "next-themes";
import { ThemeToggle } from "@/components/layout/ThemeToggle";

export const metadata = {
  title: "SocialSphere",
  description: "SocialSphere",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          {children}
          <ThemeToggle />
        </ThemeProvider>
      </body>
    </html>
  );
}
```

### Nested Layout

Layouts can be nested. The `(main)` route group has its own layout that adds the Navbar and WebSocket providers:

<!-- src/app/(main)/layout.js -->
```jsx
import Navbar from "@/components/layout/Navbar";
import LiveSocketWrapper from "@/components/providers/LiveSocketWrapper";
import { ToastProvider } from "@/context/ToastContext";

export const dynamic = 'force-dynamic';

export default function MainLayout({ children }) {
    const wsUrl = process.env.LIVE;

    return (
        <LiveSocketWrapper wsUrl={wsUrl}>
            <ToastProvider>
                <div className="min-h-screen flex flex-col bg-(--muted)/6">
                    <Navbar />
                    <main className="flex-1 w-full">
                        {children}
                    </main>
                </div>
            </ToastProvider>
        </LiveSocketWrapper>
    );
}
```

Layout nesting means: Root Layout wraps Main Layout wraps your Page. Each layout renders its `{children}`, which is the next layout or page in the chain.

---

## Route Groups

Folders wrapped in parentheses `()` are **route groups** — they organize code without affecting the URL.

```
src/app/
├── (auth)/              ← Route group (NOT in the URL)
│   ├── login/
│   │   ├── layout.js
│   │   └── page.js      → URL: /login
│   └── register/
│       ├── layout.js
│       └── page.js       → URL: /register
│
├── (main)/              ← Route group (NOT in the URL)
│   ├── layout.js         → Adds Navbar + WebSocket to all main pages
│   ├── feed/
│   │   ├── public/
│   │   │   └── page.js   → URL: /feed/public
│   │   └── friends/
│   │       └── page.js   → URL: /feed/friends
│   ├── groups/
│   │   ├── page.js        → URL: /groups
│   │   └── [id]/
│   │       └── page.js    → URL: /groups/123
│   ├── posts/
│   │   └── [id]/
│   │       └── page.js    → URL: /posts/456
│   ├── profile/
│   │   └── [id]/
│   │       ├── page.js    → URL: /profile/789
│   │       └── settings/
│   │           └── page.js → URL: /profile/789/settings
│   ├── messages/
│   │   ├── layout.js
│   │   ├── page.js        → URL: /messages
│   │   └── [id]/
│   │       └── page.js    → URL: /messages/abc
│   └── notifications/
│       └── page.js        → URL: /notifications
│
├── about/
│   └── page.js            → URL: /about
├── layout.js              → Root layout (wraps everything)
└── page.js                → URL: / (landing page)
```

This separation means:
- `(auth)` pages (login, register) get a minimal layout — no Navbar, no WebSocket
- `(main)` pages get the full app shell — Navbar, WebSocket, Toast notifications

---

## Dynamic Routes

Folders with `[brackets]` create dynamic routes — the bracket name becomes a parameter.

<!-- src/app/(main)/groups/[id]/page.js -->
```jsx
// The URL /groups/abc123 makes params.id = "abc123"
export default async function GroupPage({ params }) {
    const { id } = await params;

    // Use the id to fetch group data from the backend
    const group = await getGroup({ id });

    return <GroupPageContent group={group.data} />;
}
```

Common dynamic routes in this project:
- `/groups/[id]` — a specific group
- `/posts/[id]` — a specific post
- `/profile/[id]` — a user's profile
- `/messages/[id]` — a conversation with a specific user

---

## Server Components vs Client Components

This is the most important concept in Next.js. **By default, all components are server components.**

### Server Components (default)

- Run on the server during the request
- Can directly call server-side code (database queries, API calls with secrets)
- **Cannot** use `useState`, `useEffect`, or browser APIs
- Their code never reaches the browser

### Client Components (`"use client"`)

- Run in the browser
- **Can** use `useState`, `useEffect`, `useRef`, event handlers, browser APIs
- Marked with `"use client"` at the top of the file
- Also render on the server first (for the initial HTML), then "hydrate" in the browser

### The Pattern in This Project

The typical pattern is: a **server component** fetches data, then passes it to a **client component** for interactivity.

<!-- src/app/(main)/feed/public/page.js — SERVER component -->
```jsx
import { getPublicPosts } from "@/actions/posts/get-public-posts";
import PublicFeedContent from "@/components/feed/PublicFeedContent";

export const metadata = {
    title: "Public Feed",
};

// This is a server component — no "use client" at the top
// It can call server actions directly and use async/await
export default async function PublicFeedPage() {
    const limit = 10;
    const offset = 0;
    const result = await getPublicPosts({ limit, offset });

    // Pass the fetched data as a prop to the client component
    return <PublicFeedContent initialPosts={result.success ? result.data : []} />;
}
```

<!-- src/components/feed/PublicFeedContent.js — CLIENT component -->
```jsx
"use client";  // ← This makes it a client component

import { useState, useEffect, useRef, useCallback } from "react";

export default function PublicFeedContent({ initialPosts }) {
    // Can use state, effects, and browser APIs
    const [posts, setPosts] = useState(initialPosts || []);
    const [loading, setLoading] = useState(false);
    // ...
}
```

---

## Server Actions

Server actions are functions marked with `"use server"` that run on the server but can be called from client components as if they were regular functions. Next.js handles the network request transparently.

<!-- src/actions/auth/login.js -->
```jsx
"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function login(credentials) {
    try {
        const apiResp = await serverApiRequest("/login", {
            method: "POST",
            body: JSON.stringify(credentials),
            forwardCookies: true,
            headers: { "Content-Type": "application/json" },
        });

        if (!apiResp.ok) {
            return { success: false, status: apiResp.status, error: apiResp.message };
        }

        return {
            success: true,
            user_id: apiResp.data.id,
            username: apiResp.data.username,
            avatar_url: apiResp.data.avatar_url || "",
        };
    } catch (error) {
        return { success: false, error: error.message };
    }
}
```

The client component imports and calls it like a normal function:

```jsx
import { login } from "@/actions/auth/login";

// Inside a client component:
const resp = await login({ email, password });
```

Behind the scenes, Next.js sends a POST request to the server, runs the function, and returns the result. You never manage URLs or `fetch()` calls. See [Actions vs API Routes](./04-actions-vs-api.md) for a deep dive.

---

## Navigation

Next.js provides several ways to navigate between pages.

### `<Link>` Component

For declarative navigation (like `<a>` tags but with client-side transitions):

```jsx
import Link from "next/link";

<Link href="/feed/public">Public Feed</Link>
<Link href={`/profile/${user.id}`}>My Profile</Link>
```

### `useRouter()` Hook

For programmatic navigation (navigate in response to events):

```jsx
import { useRouter } from "next/navigation";

const router = useRouter();

// Navigate after an action
router.push(`/messages/${userId}`);
```

### `usePathname()` Hook

To get the current URL path (useful for highlighting active nav items):

<!-- src/components/layout/Navbar.js -->
```jsx
import { usePathname } from "next/navigation";

const pathname = usePathname();

const isActive = (path) => pathname === path;

// Use it to style the active nav item differently
<a
    href={item.href}
    className={isActive(item.href)
        ? "bg-(--accent)/10 text-(--accent)"    // Active style
        : "text-(--muted) hover:text-foreground"  // Inactive style
    }
>
```

---

## Image Component

Next.js provides an `<Image>` component with automatic optimization. The project configures allowed remote image sources in `next.config.mjs`:

<!-- next.config.mjs -->
```js
images: {
    remotePatterns: [
        {
            protocol: 'http',
            hostname: 'localhost',
            port: '9000',
            pathname: '/uploads/**',
        },
        // ... more patterns for Docker/Minio
    ],
},
```

This allows `next/image` to optimize images from these external sources. For avatar images, the project often uses standard `<img>` tags with error handling for dynamic URLs.

---

## Metadata

Each page can export a `metadata` object to set the page title and description:

```jsx
export const metadata = {
    title: "Public Feed",
};
```

The root layout sets the default metadata, and pages can override it:

```jsx
// src/app/layout.js
export const metadata = {
  title: "SocialSphere",
  description: "SocialSphere",
};
```

---

## Proxy (Request Interceptor)

> **Note:** Next.js 16 renamed `middleware.js` to `proxy.js` and the export from `middleware` to `proxy`. The `config.matcher` syntax is unchanged.

The proxy runs on every matching request **before** the page loads. It's used for authentication checks and redirects.

<!-- src/proxy.js -->
```jsx
import { NextResponse } from "next/server";

const protectedRoutes = ["/feed", "/groups", "/posts", "/profile"];
const authRoutes = ["/login", "/register", "/"];

function isTokenExpired(token) {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.exp && payload.exp < Math.floor(Date.now() / 1000);
  } catch {
    return true;
  }
}

export function proxy(request) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get("jwt")?.value;

  // Expired token → delete it and redirect to login
  if (token && isTokenExpired(token)) {
    const response = NextResponse.redirect(new URL("/login", request.url));
    response.cookies.delete("jwt");
    return response;
  }

  const isAuthenticated = !!token;
  const isProtectedRoute = protectedRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`)
  );
  const isAuthRoute = authRoutes.some(
    (route) => pathname === route || (route !== "/" && pathname.startsWith(`${route}/`))
  );

  // Not logged in + protected route → redirect to /login
  if (!isAuthenticated && isProtectedRoute) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("callbackUrl", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Logged in + auth route (login/register) → redirect to /feed/public
  if (isAuthenticated && isAuthRoute) {
    return NextResponse.redirect(new URL("/feed/public", request.url));
  }

  return NextResponse.next(); // Continue normally
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
```

Key points:
- Runs on the server for every request matching the `config.matcher` pattern
- Checks the JWT cookie for authentication
- Redirects unauthenticated users away from protected routes
- Redirects authenticated users away from login/register pages
- Runs in the Node.js runtime only (edge runtime is no longer supported in Next.js 16)

---

## Config Files

| File | Purpose |
|------|---------|
| `next.config.mjs` | Next.js configuration — output mode, image domains, external packages |
| `jsconfig.json` | Path aliases — `@/*` maps to `./src/*` so you can write `import X from "@/lib/..."` |
| `postcss.config.mjs` | PostCSS plugins — configures Tailwind CSS |

<!-- jsconfig.json -->
```json
{
  "compilerOptions": {
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

This lets you write `import { login } from "@/actions/auth/login"` instead of `import { login } from "../../../actions/auth/login"`.

<!-- next.config.mjs -->
```js
const nextConfig = {
  reactStrictMode: false,
  reactCompiler: true,         // React Compiler enabled
  output: 'standalone',        // Standalone output for Docker
  serverExternalPackages: ["@grpc/grpc-js"], // Don't bundle gRPC (used by OTEL)
  images: {
    remotePatterns: [ /* ... */ ],
  },
};
```

---

## Next Steps

Now that you understand Next.js concepts, continue to [Architecture](./03-architecture.md) for a deep dive into how the SocialSphere project is structured.
