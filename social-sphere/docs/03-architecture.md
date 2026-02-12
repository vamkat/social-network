# Architecture Deep Dive

> **Prerequisites:** [React Fundamentals](./01-react-fundamentals.md) and [Next.js Fundamentals](./02-nextjs-fundamentals.md).
>
> This guide maps out the entire SocialSphere frontend structure.

---

## Project Root Overview

```
social-sphere/
├── docs/                    ← You are here
├── public/                  ← Static assets (images, alert sounds)
├── src/                     ← All source code
│   ├── app/                 ← Routes and pages
│   ├── actions/             ← Server actions (backend API calls)
│   ├── components/          ← Reusable UI components
│   ├── context/             ← React Context providers
│   ├── hooks/               ← Custom React hooks
│   ├── lib/                 ← Shared utilities
│   ├── store/               ← Zustand global state
│   ├── instrumentation.js   ← OTEL SDK bootstrap
│   └── proxy.js             ← Request interceptor (auth redirects)
├── next.config.mjs          ← Next.js configuration
├── jsconfig.json            ← Path alias: @/* → ./src/*
├── postcss.config.mjs       ← Tailwind CSS via PostCSS
├── package.json             ← Dependencies and scripts
└── .env                     ← Environment variables (not committed)
```

### Key Dependencies (`package.json`)

| Package | Purpose |
|---------|---------|
| `next` (^16) | Framework — routing, SSR, server actions |
| `react` / `react-dom` (19.2.1) | UI library |
| `zustand` (^5) | Lightweight global state management |
| `motion` (^12) | Animation library (Motion for React) |
| `lucide-react` | Icon library |
| `next-themes` | Dark/light mode switching |
| `date-fns` | Date formatting utilities |
| `emoji-picker-react` | Emoji picker component |
| `react-datepicker` | Date input component |
| `@opentelemetry/*` | Observability — traces and logs |
| `@grpc/grpc-js` | gRPC transport for OTEL exporters |

---

## `src/app/` — Route Structure

The app directory defines all pages and their layout hierarchy:

```
src/app/
├── layout.js                → Root layout: <html>, ThemeProvider, globals.css
├── page.js                  → Landing page (/)
├── globals.css              → Theme variables, custom utility classes, animations
│
├── (auth)/                  → Auth pages (minimal layout, no Navbar)
│   ├── login/
│   │   ├── layout.js        → Auth page layout (centered card)
│   │   └── page.js          → Login page
│   └── register/
│       ├── layout.js        → Auth page layout
│       └── page.js          → Registration page
│
├── (main)/                  → App pages (Navbar + WebSocket + Toast)
│   ├── layout.js            → MainLayout: LiveSocketWrapper → ToastProvider → Navbar
│   ├── feed/
│   │   ├── public/page.js   → Public feed
│   │   └── friends/page.js  → Friends-only feed
│   ├── groups/
│   │   ├── page.js          → Groups listing
│   │   └── [id]/page.js     → Individual group
│   ├── posts/
│   │   └── [id]/page.js     → Individual post
│   ├── profile/
│   │   └── [id]/
│   │       ├── page.js      → User profile
│   │       ├── loading.js   → Profile loading skeleton
│   │       └── settings/page.js → Profile settings
│   ├── messages/
│   │   ├── layout.js        → Messages layout (with ConversationsProvider)
│   │   ├── page.js          → Conversations list
│   │   └── [id]/page.js     → Chat with specific user
│   └── notifications/
│       └── page.js          → Notifications page
│
└── about/page.js            → About page
```

### Layout Nesting

When you visit `/feed/public`, the layouts nest like this:

```
RootLayout (app/layout.js)
  └─ MainLayout (app/(main)/layout.js)
       └─ PublicFeedPage (app/(main)/feed/public/page.js)
```

The `(main)/layout.js` provider stack is important:

```
LiveSocketWrapper       ← WebSocket connection for real-time features
  └─ ToastProvider      ← Toast notification system (listens to WebSocket)
       └─ Navbar        ← Navigation bar
            └─ {children} ← The actual page content
```

---

## `src/actions/` — Server Actions

All server actions follow the same pattern: `"use server"` → call `serverApiRequest()` → return `{success, data/error}`.

### By Feature

**auth/**
| File | Purpose |
|------|---------|
| `login.js` | Authenticate user, forward JWT cookie |
| `logout.js` | Clear JWT cookie, redirect to login |
| `register.js` | Create new account with avatar upload |
| `validate-upload.js` | Server-side image validation |
| `get-image-url.js` | Fetch fresh image URL when cached one expires |

**posts/**
| File | Purpose |
|------|---------|
| `get-public-posts.js` | Paginated public feed |
| `get-friends-posts.js` | Paginated friends-only feed |
| `get-user-posts.js` | Posts by a specific user |
| `get-post.js` | Single post by ID |
| `create-post.js` | Create new post (with optional image) |
| `edit-post.js` | Edit existing post |
| `delete-post.js` | Delete a post |
| `get-comments.js` | Comments on a post |
| `create-comment.js` | Add comment to a post |
| `edit-comment.js` | Edit a comment |
| `delete-comment.js` | Delete a comment |
| `toggle-reaction.js` | Like/unlike a post |
| `who-liked-entity.js` | Get users who liked a post/comment |

**groups/**
| File | Purpose |
|------|---------|
| `get-all-groups.js` | All groups listing |
| `get-user-groups.js` | Groups the current user belongs to |
| `get-most-popular.js` | Most popular groups |
| `get-group.js` | Single group details |
| `get-group-posts.js` | Posts within a group |
| `create-group.js` | Create new group |
| `update-group.js` | Update group details |
| `group-members.js` | List group members |
| `invite-to-group.js` | Invite user to group |
| `respond-to-invite.js` | Accept/reject group invitation |
| `request-join-group.js` | Request to join a group |
| `handle-join-request.js` | Accept/reject join request (for group owners) |
| `cancel-join-request.js` | Cancel pending join request |
| `leave-group.js` | Leave a group |
| `remove-from-group.js` | Remove member from group |
| `get-pending-requests.js` | Pending join requests for a group |
| `get-pening-count.js` | Count of pending requests |
| `search-groups.js` | Search groups by name |

**chat/**
| File | Purpose |
|------|---------|
| `get-conv.js` | List conversations (paginated) |
| `get-conv-by-id.js` | Single conversation details |
| `get-messages.js` | Messages in a private conversation |
| `get-group-messages.js` | Messages in a group chat |
| `send-msg.js` | Send private message (HTTP fallback) |
| `send-group-msg.js` | Send group message (HTTP fallback) |
| `get-unread-count.js` | Total unread message count |
| `mark-read.js` | Mark conversation as read |

**events/**
| File | Purpose |
|------|---------|
| `get-group-events.js` | Events for a group |
| `create-event.js` | Create new group event |
| `edit-event.js` | Edit an event |
| `delete-event.js` | Delete an event |
| `respond-to-event.js` | RSVP to an event |
| `remove-event-response.js` | Remove RSVP |

**notifs/**
| File | Purpose |
|------|---------|
| `get-user-notifs.js` | Paginated notifications |
| `get-notif-count.js` | Unread notification count |
| `mark-as-read.js` | Mark single notification as read |
| `mark-all-as-read.js` | Mark all notifications as read |
| `delete-notification.js` | Delete a notification |

**Other**
| File | Purpose |
|------|---------|
| `profile/get-profile-info.js` | User profile data |
| `profile/update-profile.js` | Update profile fields |
| `profile/settings.js` | Privacy/account settings |
| `requests/follow-user.js` | Send follow request |
| `requests/unfollow-user.js` | Unfollow a user |
| `requests/handle-request.js` | Accept/reject follow request |
| `search/search-users.js` | Search users by name/username |
| `users/get-followers.js` | List a user's followers |
| `users/get-following.js` | List who a user follows |
| `users/get-not-invited.js` | Users not yet invited to a group |

---

## `src/components/` — UI Components

### `layout/`
| Component | Purpose |
|-----------|---------|
| `Navbar.js` | Top navigation bar — search, nav links, messages dropdown, notifications, user menu |
| `Container.js` | Responsive width wrapper with size variants (narrow/default/wide/full) |
| `ThemeToggle.js` | Dark/light mode toggle button |

### `ui/`
| Component | Purpose |
|-----------|---------|
| `Modal.js` | Universal modal dialog with header, body, footer, confirm/cancel |
| `Tooltip.js` | Hover tooltip with arrow |
| `PostCard.js` | Post display card — author, content, image, reactions, comments |
| `SinglePostCard.js` | Full-page single post view |
| `PostImage.js` | Post image with loading and error states |
| `CreatePost.js` | Post creation form with image upload |
| `Toast.js` | Individual toast notification |
| `ToastContainer.js` | Toast stack manager |
| `LoadingDots.js` | Animated loading indicator |

### `forms/`
| Component | Purpose |
|-----------|---------|
| `LoginForm.js` | Login form with email/password |
| `RegisterForm.js` | Registration form with avatar upload |
| `ProfileForm.js` | Profile editing form |
| `PrivacyForm.js` | Privacy settings form |
| `SecurityForm.js` | Security/password settings form |

### `feed/`
| Component | Purpose |
|-----------|---------|
| `PublicFeedContent.js` | Public feed with infinite scroll |
| `FriendsFeedContent.js` | Friends feed with infinite scroll |

### `groups/`
| Component | Purpose |
|-----------|---------|
| `GroupsContent.js` | Groups listing page content |
| `GroupsPagination.js` | Paginated groups navigation |
| `GroupCard.js` | Group preview card |
| `GroupPageContent.js` | Individual group page content |
| `GroupHeader.js` | Group page header with actions |
| `CreateGroup.js` | Group creation form |
| `UpdateGroupModal.js` | Group editing modal |
| `GroupPostCard.js` | Post card for group context |
| `CreatePostGroup.js` | Post creation for groups |
| `EventCard.js` | Group event display card |
| `CreateEventModal.js` | Event creation modal |
| `EditEventModal.js` | Event editing modal |

### `messages/`
| Component | Purpose |
|-----------|---------|
| `ConversationsContent.js` | Conversations list with real-time updates |
| `MessagesContent.js` | Chat view with message history and sending |

### `notifications/`
| Component | Purpose |
|-----------|---------|
| `NotificationCard.js` | Individual notification with actions |
| `NotificationsContent.js` | Notifications page content |

### `profile/`
| Component | Purpose |
|-----------|---------|
| `ProfileHeader.js` | Profile page header with avatar, stats, follow button |
| `ProfileContent.js` | Profile page content (posts, followers, following) |
| `ProfileStats.js` | Followers/following count display |

### `providers/`
| Component | Purpose |
|-----------|---------|
| `LiveSocketWrapper.js` | Client component wrapper for `LiveSocketProvider` |

---

## `src/lib/` — Shared Utilities

### `server-api.js` — The Single Chokepoint

Every backend request flows through `serverApiRequest()`. This is the most important utility file.

```
Component → Server Action → serverApiRequest() → Go Backend
```

What it does:
1. Reads the JWT cookie and injects it into the request
2. Injects W3C trace context headers for distributed tracing
3. Makes the `fetch()` call to the Go backend
4. Forwards `Set-Cookie` headers back (for login)
5. Handles error responses:
   - **401 Unauthorized** → deletes JWT cookie, redirects to `/login`
   - **403 Forbidden** → returns `{ok: false, message: "Forbidden"}`
   - **400 Bad Request** → returns `{ok: false, message: "..."}`
6. Logs every request (outgoing, succeeded, failed, exception) with timing
7. Returns `{ok: true, data: ...}` or `{ok: false, message: "..."}` — never throws

See [Actions vs API Routes](./04-actions-vs-api.md) for a full annotated walkthrough.

### `validation.js` — Input Validation

Shared validators used in both forms and server actions:

- `EMAIL_PATTERN` — email regex
- `STRONG_PASSWORD_PATTERN` — at least 1 lowercase, 1 uppercase, 1 number, 1 symbol
- `USERNAME_PATTERN` — letters, numbers, dots, underscores, dashes
- `MAX_FILE_SIZE` — 5MB limit for images
- `ALLOWED_FILE_TYPES` — JPEG, PNG, GIF, WebP
- `validateRegistrationForm()` — full registration validation
- `validateProfileForm()` — profile update validation
- `validateLoginForm()` — login validation
- `validatePostContent()` — post content length check
- `validateImage()` — file type, size, and dimension check

### `logger.server.js` — Server Logger

Dual-output logger: OTEL LogRecord + stdout. Uses the `@N` template syntax that matches the Go backend:

```js
logger.info("outgoing request @1 @2", "method", "POST", "url", "/login");
// Output: 14:32:01.123 [SOC]: INFO outgoing request method=POST url=/login
```

- `@1` is replaced with the first key-value pair (`method=POST`)
- `@2` is replaced with the second pair (`url=/login`)
- Also emits an OTEL log record for centralized log collection

### `logger.client.js` — Client Logger

Console-only logger with the same API and template syntax. Used in browser-side components.

### `logger.js` — Barrel File

Re-exports the client logger. Server code imports `@/lib/logger.server` directly.

### `time.js` — Relative Time

`getRelativeTime(timestamp)` — converts ISO timestamps to "2 mins ago", "1 hour ago", etc.

### `notifications.js` — Notification Formatting

`constructNotif(notif)` — transforms raw notification objects into display-friendly format with message text, user links, and action callbacks.

---

## `src/context/` — React Context Providers

| Context | File | Purpose |
|---------|------|---------|
| `LiveSocketContext` | `LiveSocketContext.js` | WebSocket connection, message routing, group subscriptions |
| `ToastContext` | `ToastContext.js` | Toast notification queue (max 3, auto-dismiss 4s, hover pause) |
| `ConversationsContext` | `ConversationsContext.js` | Conversation list state, real-time updates, pagination |

See [State and Context](./05-state-and-context.md) for detailed walkthroughs.

---

## `src/store/` — Zustand Global State

`store.js` exports two stores:

- **`useStore`** — main app state: `user`, `unreadCount`, `unreadNotifs`, `hasMsg`. Only `user` is persisted to `localStorage`.
- **`useMsgReceiver`** — message recipient state: `msgReceiver`. Persisted to `localStorage` so it survives page navigation.

---

## `src/hooks/` — Custom Hooks

- **`useFormValidation.js`** — field-level validation with error tracking. Provides `errors` state and `validateField(name, value, validator)`.

---

## Special Files

### `instrumentation.js` + `instrumentation-node.js`

Next.js automatically calls `register()` from `instrumentation.js` on server startup. It conditionally imports the Node.js OTEL SDK setup:

```js
export async function register() {
    if (process.env.NEXT_RUNTIME === "nodejs") {
        await import("./instrumentation-node.js");
    }
}
```

`instrumentation-node.js` initializes OpenTelemetry with gRPC exporters for traces and logs, sending telemetry data to the Grafana Alloy collector.

### `proxy.js`

The request interceptor that runs before every page load. Handles authentication redirects — see [Next.js Fundamentals](./02-nextjs-fundamentals.md#proxy-request-interceptor) for the full code.

---

## `globals.css` — Theme and Custom Classes

Defines the project's visual system:

- **CSS variables** — `--accent`, `--foreground`, `--muted`, `--background`, `--border` for light/dark themes
- **Form elements** — `.form-input`, `.form-label`, `.form-error`, `.form-toggle-btn`
- **Buttons** — `.btn`, `.btn-primary`, `.btn-secondary`
- **Typography** — `.heading-xl`, `.heading-lg`, `.heading-md`, `.heading-sm`
- **Layout** — `.section-divider`, `.page-container`
- **Post styles** — `.post-card`, `.post-avatar`, `.post-text`, `.post-actions`
- **Animations** — `.animate-fade-in`, `.animate-pulse-glow`
- **Custom scrollbars** — thin, themed scrollbars for Webkit and Firefox

---

## Data Flow Diagram

```
┌─────────┐     ┌──────────────────┐     ┌───────────────┐     ┌──────────────────┐     ┌─────────────┐
│ Browser  │────→│ Server Component │────→│ Server Action │────→│ serverApiRequest │────→│ Go Backend  │
│          │     │ (page.js)        │     │ (actions/*.js)│     │ (server-api.js)  │     │ (gateway)   │
│          │     │                  │     │               │     │                  │     │             │
│          │←────│  Returns JSX     │←────│ {success,data}│←────│ {ok, data}       │←────│ JSON resp   │
└─────────┘     └──────────────────┘     └───────────────┘     └──────────────────┘     └─────────────┘
                                                                       │
                                                                       ├─ Injects JWT cookie
                                                                       ├─ Injects W3C traceparent
                                                                       ├─ Handles 401 → redirect
                                                                       └─ Logs all requests
```

For real-time features, the data flow is different:

```
┌─────────┐     ┌─────────────────────┐     ┌──────────────┐
│ Browser  │←──→│ LiveSocketContext    │←──→│ Go LiveService│
│          │ WS │ (WebSocket client)   │ WS │ (via gateway) │
│          │     │                     │     │              │
│ Toast    │←────│ Routes messages to: │     │              │
│ Navbar   │←────│ - private listeners │     │              │
│ Messages │←────│ - group listeners   │     │              │
│          │←────│ - notif listeners   │     │              │
└─────────┘     └─────────────────────┘     └──────────────┘
```

---

## Next Steps

Continue to [Actions vs API Routes](./04-actions-vs-api.md) to understand why this project uses server actions instead of traditional API routes, with a full login flow walkthrough.
