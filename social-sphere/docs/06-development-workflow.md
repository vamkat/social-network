# Development Workflow

> **Prerequisites:** All previous docs — [React Fundamentals](./01-react-fundamentals.md), [Next.js Fundamentals](./02-nextjs-fundamentals.md), [Architecture](./03-architecture.md), [Actions vs API](./04-actions-vs-api.md), [State and Context](./05-state-and-context.md).
>
> This guide covers how to run, develop, and add features to SocialSphere.

---

## Prerequisites and Setup

### Requirements

- **Node.js** (v18+)
- **npm** (comes with Node.js)
- The Go backend services running (via Docker Compose or locally)

### Installation

```bash
cd social-sphere
npm install
```

### Environment Variables

Create a `.env` file in `social-sphere/`:

```env
# Backend gateway URL (where the Go API is running)
GATEWAY=http://localhost:8080

# WebSocket URL for real-time features
LIVE=ws://localhost:8082

# OpenTelemetry (optional — for observability)
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=social-sphere

# Logging
LOG_LEVEL=INFO
ENABLE_DEBUG_LOGS=false
```

| Variable | Purpose |
|----------|---------|
| `GATEWAY` | URL of the Go API gateway — all `serverApiRequest()` calls go here |
| `LIVE` | WebSocket URL for `LiveSocketContext` — real-time messages and notifications |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Where to send traces/logs (Grafana Alloy in production) |
| `OTEL_SERVICE_NAME` | Service name in telemetry data |
| `LOG_LEVEL` | Minimum log level: `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `ENABLE_DEBUG_LOGS` | Set to `"true"` to enable DEBUG level server logs |

---

## Running the Dev Server

```bash
# Development (hot reload)
npm run dev

# Production build
npm run build

# Start production server
npm run start
```

The dev server runs on `http://localhost:3000` by default.

---

## Project Structure Quick Reference

```
src/
├── app/                        # Pages and routes
│   ├── (auth)/login/page.js    #   Login page
│   ├── (auth)/register/page.js #   Register page
│   ├── (main)/layout.js        #   Main app shell (Navbar, WebSocket, Toast)
│   ├── (main)/feed/            #   Feed pages
│   ├── (main)/groups/          #   Groups pages
│   ├── (main)/profile/         #   Profile pages
│   ├── (main)/messages/        #   Messages pages
│   └── (main)/notifications/   #   Notifications page
│
├── actions/                    # Server actions (one per API call)
│   ├── auth/                   #   login, logout, register
│   ├── posts/                  #   CRUD for posts/comments
│   ├── groups/                 #   CRUD for groups
│   ├── chat/                   #   Conversations and messages
│   ├── events/                 #   Group events
│   ├── notifs/                 #   Notifications
│   ├── profile/                #   Profile management
│   ├── requests/               #   Follow/unfollow
│   ├── search/                 #   User search
│   └── users/                  #   Followers/following lists
│
├── components/                 # Reusable UI
│   ├── ui/                     #   Generic: Modal, Tooltip, PostCard, Toast
│   ├── forms/                  #   LoginForm, RegisterForm, ProfileForm
│   ├── layout/                 #   Navbar, Container, ThemeToggle
│   ├── feed/                   #   PublicFeedContent, FriendsFeedContent
│   ├── groups/                 #   GroupCard, GroupPageContent, EventCard
│   ├── messages/               #   ConversationsContent, MessagesContent
│   ├── notifications/          #   NotificationCard, NotificationsContent
│   ├── profile/                #   ProfileHeader, ProfileContent
│   └── providers/              #   LiveSocketWrapper
│
├── context/                    # React Context providers
│   ├── LiveSocketContext.js    #   WebSocket connection
│   ├── ToastContext.js         #   Toast notifications
│   └── ConversationsContext.js #   Conversation list
│
├── store/store.js              # Zustand global state
├── hooks/useFormValidation.js  # Custom form validation hook
├── lib/                        # Utilities
│   ├── server-api.js           #   Backend request chokepoint
│   ├── validation.js           #   Input validation
│   ├── logger.server.js        #   Server-side logger (OTEL + stdout)
│   ├── logger.client.js        #   Client-side logger (console)
│   ├── logger.js               #   Logger barrel file
│   ├── time.js                 #   Relative time formatting
│   └── notifications.js        #   Notification message construction
│
├── instrumentation.js          # OTEL bootstrap
└── proxy.js                    # Auth redirects
```

---

## How to Add a New Page

Follow the pattern established by the public feed page.

### 1. Create the route directory and `page.js` (server component)

```bash
mkdir -p src/app/(main)/your-feature
```

<!-- src/app/(main)/your-feature/page.js -->
```jsx
import { getYourData } from "@/actions/your-feature/get-data";
import YourFeatureContent from "@/components/your-feature/YourFeatureContent";

export const metadata = {
    title: "Your Feature",
};

export default async function YourFeaturePage() {
    // Fetch initial data on the server
    const result = await getYourData({ limit: 10, offset: 0 });

    // Pass data to client component
    return <YourFeatureContent initialData={result.success ? result.data : []} />;
}
```

### 2. Create the client component for interactivity

<!-- src/components/your-feature/YourFeatureContent.js -->
```jsx
"use client";

import { useState } from "react";
import Container from "@/components/layout/Container";

export default function YourFeatureContent({ initialData }) {
    const [data, setData] = useState(initialData || []);

    return (
        <Container>
            <h1 className="heading-md mt-8">Your Feature</h1>
            {data.map((item) => (
                <div key={item.id}>{item.name}</div>
            ))}
        </Container>
    );
}
```

### 3. Create the server action

<!-- src/actions/your-feature/get-data.js -->
```jsx
"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getYourData({ limit, offset }) {
    const result = await serverApiRequest(`/your-endpoint?limit=${limit}&offset=${offset}`);
    if (!result.ok) {
        return { success: false, error: result.message };
    }
    return { success: true, data: result.data };
}
```

---

## How to Add a Server Action

Every action follows the canonical pattern:

```jsx
"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function yourAction(params) {
    try {
        const result = await serverApiRequest("/your-endpoint", {
            method: "POST",  // or GET, PUT, DELETE
            body: JSON.stringify(params),
            headers: { "Content-Type": "application/json" },
        });

        if (!result.ok) {
            return { success: false, error: result.message };
        }

        return { success: true, data: result.data };
    } catch (error) {
        return { success: false, error: error.message };
    }
}
```

For actions that include file uploads, use `FormData` instead of `JSON.stringify`:

```jsx
export async function createPost(formData) {
    const result = await serverApiRequest("/posts", {
        method: "POST",
        body: formData,  // Don't set Content-Type — fetch sets multipart/form-data automatically
    });
    // ...
}
```

Place the file in `src/actions/<feature>/` matching the feature it belongs to.

---

## How to Add a Component

### Decision: Server or Client?

| Need | Component type | Directive |
|------|---------------|-----------|
| Fetch data, no interactivity | Server component | *(none — default)* |
| `useState`, `useEffect`, onClick, browser APIs | Client component | `"use client"` |
| Both (fetch + interact) | Server page → client component | Split into two files |

### Placement

Place components in `src/components/<feature>/`:

- `ui/` — generic, reusable across features
- `forms/` — form components
- `layout/` — structural components (Navbar, Container)
- `<feature>/` — specific to a feature (groups, messages, etc.)

### Naming

- **PascalCase** file names: `PostCard.js`, `CreateEventModal.js`
- **`export default function ComponentName`** — one component per file
- Match the file name to the component name

---

## How to Add a Zustand Store Field

Edit `src/store/store.js`:

```jsx
export const useStore = create(
  persist(
    (set) => ({
      // ... existing fields ...

      // Add your new field
      yourField: initialValue,
      setYourField: (value) => set({ yourField: value }),
    }),
    {
      name: 'user',
      partialize: (state) => ({
        user: state.user,
        // Add here ONLY if it needs to survive page refreshes:
        // yourField: state.yourField,
      }),
    }
  )
);
```

Then use it in components:

```jsx
const yourField = useStore((state) => state.yourField);
const setYourField = useStore((state) => state.setYourField);
```

**Important:** Only add fields to `partialize` if they truly need to persist across page refreshes. Most data should be fetched fresh from the server.

---

## Common Patterns

### Form Handling

Forms in this project use the native `FormData` API with `onSubmit`:

```jsx
async function handleSubmit(event) {
    event.preventDefault();
    setIsLoading(true);
    setError("");

    const formData = new FormData(event.currentTarget);
    const value = formData.get("fieldName");

    const result = await someAction({ value });
    if (!result.success) {
        setError(result.error);
        setIsLoading(false);
        return;
    }

    // Success handling
    setIsLoading(false);
}
```

### Infinite Scroll (IntersectionObserver)

Used in feeds to load more content when the user scrolls near the bottom:

```jsx
const [items, setItems] = useState(initialItems);
const [offset, setOffset] = useState(10);
const [hasMore, setHasMore] = useState(initialItems.length >= 10);
const [loading, setLoading] = useState(false);
const observerTarget = useRef(null);

const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;
    setLoading(true);

    const result = await getItems({ limit: 5, offset });
    const newItems = result.success ? result.data : [];

    if (newItems.length > 0) {
        setItems((prev) => [...prev, ...newItems]);
        setOffset((prev) => prev + 5);
        if (newItems.length < 5) setHasMore(false);
    } else {
        setHasMore(false);
    }
    setLoading(false);
}, [offset, loading, hasMore]);

useEffect(() => {
    const observer = new IntersectionObserver(
        (entries) => {
            if (entries[0].isIntersecting && hasMore && !loading) {
                loadMore();
            }
        },
        { threshold: 0.1 }
    );
    if (observerTarget.current) observer.observe(observerTarget.current);
    return () => {
        if (observerTarget.current) observer.unobserve(observerTarget.current);
    };
}, [loadMore, hasMore, loading]);

// In JSX:
{hasMore && <div ref={observerTarget} />}
```

### Error Handling in Components

```jsx
const [error, setError] = useState("");

// After a server action call:
const result = await someAction(params);
if (!result.success) {
    setError(result.error || "Something went wrong");
    return;
}

// In JSX:
{error && <div className="form-error">{error}</div>}
```

### File Upload

```jsx
const [file, setFile] = useState(null);

// Validate before uploading
const validation = isValidImage(file);
if (!validation.valid) {
    setError(validation.error);
    return;
}

// Build FormData
const formData = new FormData();
formData.append("content", postText);
if (file) formData.append("image", file);

// Send via server action
const result = await createPost(formData);
```

---

## Debugging Tips

### Server-Side Logs

Server-side logs appear in the terminal where you ran `npm run dev`. Look for the `[SOC]` prefix:

```
14:32:01.123 [SOC]: INFO outgoing request method=POST url=/login
14:32:01.456 [SOC]: INFO request succeeded method=POST url=/login status=200 duration_ms=333
```

Set `LOG_LEVEL=DEBUG` and `ENABLE_DEBUG_LOGS=true` in `.env` for verbose output.

### Client-Side Logs

Client-side logs appear in the browser's developer console (F12 → Console). Look for `[SOC]:`:

```
[SOC]: INFO connected
[SOC]: INFO private message received conversation_id=abc123
```

### Network Tab

Server actions appear in the browser's Network tab as POST requests to the current page URL (not `/api/...`). Look for requests with:
- **Type:** `fetch`
- **Payload:** contains the function name and arguments
- **Response:** the server action's return value

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| "Hydration mismatch" error | Server and client rendered different HTML | Ensure initial state matches between server and client |
| `useContext` returns `null` | Component outside its Provider | Check that the provider wraps the component in the layout chain |
| WebSocket won't connect | Wrong `LIVE` URL or backend not running | Check `.env` and backend logs |
| 401 redirect loop | Expired JWT or missing cookie | Clear cookies and log in again |
| "use server" function not working | Missing `"use server"` directive at top of file | Add `"use server"` as the first line |
| State lost on navigation | State not persisted in Zustand `partialize` | Add the field to `partialize` or use Context at the right layout level |

---

## Summary

To add a new feature end-to-end:

1. **Server action** in `src/actions/<feature>/` — wraps `serverApiRequest()`
2. **Server page** in `src/app/(main)/<feature>/page.js` — fetches initial data
3. **Client component** in `src/components/<feature>/` — handles interactivity
4. **Zustand field** (if needed) in `src/store/store.js` — for global state
5. **Context** (if needed) in `src/context/` — for shared subtree state

Every piece follows the patterns documented in this series. When in doubt, find an existing feature that does something similar and follow its structure.
