# State and Context

> **Prerequisites:** [React Fundamentals](./01-react-fundamentals.md) (especially `useState`, `useEffect`, `useRef`, `useCallback`), [Next.js Fundamentals](./02-nextjs-fundamentals.md), [Architecture](./03-architecture.md).
>
> This guide covers all the state management patterns used in SocialSphere.

---

## State Management Overview

The app uses four levels of state management, each for a different scope:

| Level | Tool | Scope | Example |
|-------|------|-------|---------|
| **Local** | `useState` | Single component | Form inputs, loading flags, open/closed toggles |
| **Shared** | React Context | Component subtree | WebSocket connection, toast notifications, conversations |
| **Global** | Zustand | Entire app (persisted) | Current user, unread counts |
| **Real-time** | WebSocket | Server-pushed updates | New messages, notifications |

---

## React Context Fundamentals

React Context lets you share data across a component tree without passing props through every level. The pattern is always:

1. **Create** a context
2. **Provide** a value at the top of the tree
3. **Consume** the value anywhere below with a custom hook

```jsx
// 1. Create
const MyContext = createContext(null);

// 2. Provide
function MyProvider({ children }) {
    const [value, setValue] = useState("hello");
    return (
        <MyContext.Provider value={{ value, setValue }}>
            {children}
        </MyContext.Provider>
    );
}

// 3. Consume (via custom hook for safety)
function useMyContext() {
    const context = useContext(MyContext);
    if (!context) {
        throw new Error("useMyContext must be used within a MyProvider");
    }
    return context;
}

// Usage in any child component:
function SomeComponent() {
    const { value, setValue } = useMyContext();
    return <p>{value}</p>;
}
```

Every context in this project follows this exact pattern. The custom hook ensures you get a helpful error if you forget to wrap your component with the provider.

---

## Zustand: `useStore`

Zustand is a lightweight global state library. Unlike React Context, Zustand stores don't need a provider component — you can use them anywhere.

<!-- src/store/store.js -->
```jsx
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useStore = create(
  persist(
    (set) => ({
      // === User state ===
      user: null,
      setUser: (userData) => set({ user: userData }),
      clearUser: () => set({ user: null }),

      // === Unread messages count ===
      unreadCount: 0,
      setUnreadCount: (count) => set({ unreadCount: count }),
      incrementUnreadCount: () => set((state) => ({
          unreadCount: state.unreadCount + 1
      })),
      decrementUnreadCount: (n = 1) => set((state) => ({
          unreadCount: Math.max(0, state.unreadCount - n)
      })),

      // === Unread notifications count ===
      unreadNotifs: 0,
      setUnreadNotifs: (count) => set({ unreadNotifs: count }),
      incrementNotifs: () => set((state) => ({
          unreadNotifs: state.unreadNotifs + 1
      })),
      decrementNotifs: () => set((state) => ({
          unreadNotifs: Math.max(0, state.unreadNotifs - 1)
      })),

      // === Group message flag ===
      hasMsg: false,
      setHasMsg: (hasMsg) => set({ hasMsg }),
    }),
    {
      name: 'user',  // localStorage key
      partialize: (state) => ({
        user: state.user  // Only persist user — counts are fetched fresh
      }),
    }
  )
);
```

### How to Use Zustand

```jsx
// Read a value (with a selector — only re-renders when this specific value changes)
const user = useStore((state) => state.user);
const unreadCount = useStore((state) => state.unreadCount);

// Call an action
const setUser = useStore((state) => state.setUser);
setUser({ id: "123", username: "john" });

// Or call directly
useStore.getState().incrementUnreadCount();
```

### The `persist` Middleware

The `persist` middleware saves state to `localStorage`. The `partialize` option controls **which fields** are saved:

```jsx
partialize: (state) => ({
    user: state.user  // Only save user data
})
```

This means:
- `user` survives page refreshes (saved to localStorage)
- `unreadCount`, `unreadNotifs`, `hasMsg` reset on refresh (fetched fresh from the server)

---

## Zustand: `useMsgReceiver`

A separate store for tracking the message recipient during navigation:

<!-- src/store/store.js -->
```jsx
export const useMsgReceiver = create(
  persist(
    (set) => ({
      msgReceiver: null,
      setMsgReceiver: (receiverData) => set({ msgReceiver: receiverData }),
      clearMsgReceiver: () => set({ msgReceiver: null }),
    }),
    {
      name: 'msgReceiver',  // Separate localStorage key
      partialize: (state) => ({
        msgReceiver: state.msgReceiver
      }),
    }
  )
);
```

This is persisted so when you click on a user to message them and navigate to `/messages/[id]`, the recipient data survives the page transition.

---

## `LiveSocketContext` — WebSocket Connection

The `LiveSocketContext` manages the WebSocket connection for real-time features. It's the most complex context in the app.

<!-- src/context/LiveSocketContext.js -->

### Connection Lifecycle

```jsx
export function LiveSocketProvider({ children, wsUrl }) {
    const user = useStore((state) => state.user);
    const wsRef = useRef(null);                    // WebSocket instance
    const reconnectAttemptsRef = useRef(0);        // For exponential backoff
    const subscribedGroupsRef = useRef(new Set()); // Track group subscriptions

    const [connectionState, setConnectionState] = useState(ConnectionState.DISCONNECTED);
```

The connection follows the user's login state:

```jsx
useEffect(() => {
    if (user) {
        connect();       // User logged in → open WebSocket
    } else if (hadUserRef.current) {
        disconnect(true); // User logged out → close WebSocket
    }

    return () => disconnect(false); // Cleanup on unmount
}, [user, connect, disconnect]);
```

### Auto-Reconnect with Exponential Backoff

When the connection drops unexpectedly, it reconnects with increasing delays:

```jsx
ws.onclose = (event) => {
    if (event.code !== 1000 && user) {  // 1000 = clean close
        // Delay: 1s, 2s, 4s, 8s, 16s, max 30s
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
        reconnectAttemptsRef.current++;

        reconnectTimeoutRef.current = setTimeout(() => {
            connect();
        }, delay);
    }
};
```

### Message Routing

Incoming WebSocket messages are routed based on their content:

```jsx
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    const messages = Array.isArray(data) ? data : [data];

    for (const msg of messages) {
        if (msg.group_id || msg.GroupId) {
            // Group message → notify group listeners
            groupMessageListenersRef.current.forEach((listener) => listener(msg));
        } else if (msg.conversation_id || msg.ConversationId) {
            // Private message → notify private listeners
            privateMessageListenersRef.current.forEach((listener) => listener(msg));
        } else if (msg.notification_type || msg.type) {
            // Notification → notify notification listeners
            notificationListenersRef.current.forEach((listener) => listener(msg));
        }
    }
};
```

### Listener Pattern

Any component can subscribe to specific message types:

```jsx
// In the provided context value:
const value = {
    addOnPrivateMessage,     // Register a callback for private messages
    removeOnPrivateMessage,  // Unregister it
    addOnGroupMessage,       // Register for group messages
    removeOnGroupMessage,    // Unregister
    addOnNotification,       // Register for notifications
    removeOnNotification,    // Unregister
    sendPrivateMessage,      // Send a message via WebSocket
    sendGroupMessage,        // Send a group message
    subscribeToGroup,        // Subscribe to a group's messages
    unsubscribeFromGroup,    // Unsubscribe
    // ...
};
```

Components use this pattern to listen for messages:

```jsx
// In Navbar.js — listen for private messages to update badge count
useEffect(() => {
    addOnPrivateMessage(handleNewMessage);
    return () => removeOnPrivateMessage(handleNewMessage);
}, [addOnPrivateMessage, removeOnPrivateMessage, handleNewMessage]);
```

### Group Subscriptions

When you visit a group page, the component subscribes to that group's messages:

```jsx
subscribeToGroup(groupId, isMember);  // Sends "sub:groupId" via WebSocket
unsubscribeFromGroup(groupId);        // Sends "unsub:groupId"
```

Subscriptions are tracked in a `Set` and automatically re-sent on reconnection.

---

## `ToastContext` — Toast Notifications

The toast system provides temporary popup notifications with a queue.

<!-- src/context/ToastContext.js -->

### Key Behavior

- **Max 3 visible toasts** — additional toasts go into a queue
- **Auto-dismiss after 4 seconds** — configurable per-toast
- **Hover pauses the timer** — moving the mouse away resumes countdown
- **Integrated with LiveSocket** — automatically shows toasts for incoming notifications

```jsx
export function ToastProvider({ children }) {
    const [toasts, setToasts] = useState([]);
    const [queue, setQueue] = useState([]);
    const { addOnNotification, removeOnNotification } = useLiveSocket();

    // Listen for notifications from WebSocket
    useEffect(() => {
        const handleNotification = (notification) => {
            showToast(notification);
        };
        addOnNotification(handleNotification);
        return () => removeOnNotification(handleNotification);
    }, [addOnNotification, removeOnNotification, showToast]);
```

### Timer Management

Each toast has its own auto-dismiss timer. Hovering pauses the timer and saves the remaining time:

```jsx
const pauseToast = useCallback((id) => {
    if (timersRef.current.has(id)) {
        clearTimeout(timersRef.current.get(id));
        // Calculate remaining time
        const elapsed = Date.now() - pausedAtRef.current.get(id);
        const remaining = Math.max(originalRemaining - elapsed, 1000);
        remainingTimeRef.current.set(id, remaining);
    }
}, []);

const resumeToast = useCallback((id) => {
    const remaining = remainingTimeRef.current.get(id);
    if (remaining && !timersRef.current.has(id)) {
        const timer = setTimeout(() => dismissToast(id), remaining);
        timersRef.current.set(id, timer);
    }
}, []);
```

### Queue System

When 3 toasts are already visible, new ones are queued. When a toast is dismissed, the next one from the queue is promoted:

```jsx
const dismissToast = useCallback((id) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));

    // Promote next queued toast
    setQueue((prevQueue) => {
        if (prevQueue.length === 0) return prevQueue;
        const [next, ...rest] = prevQueue;
        setToasts((prevToasts) => {
            if (prevToasts.length < MAX_VISIBLE_TOASTS) {
                startTimer(next.id);
                return [...prevToasts, next];
            }
            return prevToasts;
        });
        return rest;
    });
}, [startTimer]);
```

---

## `ConversationsContext` — Conversation List

Manages the list of conversations on the messages page with real-time updates.

<!-- src/context/ConversationsContext.js -->

```jsx
export function ConversationsProvider({ initialConversations = [], children }) {
    const [conversations, setConversations] = useState(initialConversations);
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(() => initialConversations.length >= 15);

    const { addOnPrivateMessage, removeOnPrivateMessage } = useLiveSocket();
```

### Provided Operations

| Function | Purpose |
|----------|---------|
| `addConversation(conv)` | Add a new conversation to the list |
| `updateConversation(id, updates)` | Update an existing conversation (new message, etc.) |
| `markAsRead(id)` | Set unread count to 0 |
| `loadMore()` | Load next page of conversations |

### Real-time Updates

When a private message arrives via WebSocket, the conversation list updates automatically:

```jsx
const handlePrivateMessage = useCallback(async (msg) => {
    const senderId = msg.sender?.id;
    if (senderId === user?.id) return; // Skip own messages

    const existingConv = conversationsRef.current.find(
        (conv) => conv.Interlocutor?.id === senderId
    );

    if (existingConv) {
        // Update existing conversation — new last message, bump unread count, re-sort
        setConversations((prev) =>
            prev.map((conv) => {
                if (conv.Interlocutor?.id === senderId) {
                    return {
                        ...conv,
                        LastMessage: { id: msg.id, message_text: msg.message_text, sender: msg.sender },
                        UpdatedAt: msg.created_at,
                        UnreadCount: (conv.UnreadCount || 0) + 1,
                    };
                }
                return conv;
            }).sort((a, b) => new Date(b.UpdatedAt) - new Date(a.UpdatedAt))
        );
        incrementUnreadCount();
    } else {
        // New conversation — fetch from server and add to list
        const result = await getConvByID({ interlocutorId: senderId, convId: msg.conversation_id });
        if (result.success && result.data) {
            setConversations((prev) => [{ ...result.data, UnreadCount: 1 }, ...prev]);
            incrementUnreadCount();
        }
    }
}, [user?.id, incrementUnreadCount]);
```

---

## Data Flow Diagrams

### WebSocket Message Flow

```
Go Backend (LiveService)
    │
    ▼ WebSocket
LiveSocketContext
    │
    ├──→ Navbar
    │     ├─ Updates unread message badge
    │     ├─ Updates conversation list in dropdown
    │     └─ Plays alert sound
    │
    ├──→ ToastContext
    │     └─ Shows notification popup
    │
    ├──→ ConversationsContext
    │     └─ Updates conversation list (new message, re-sort)
    │
    └──→ MessagesContent (if open)
          └─ Appends new message to chat view
```

### Login → State Restoration Flow

```
1. User submits LoginForm
    │
    ▼
2. login() server action → Go Backend → JWT cookie set
    │
    ▼
3. setUser({ id, username, avatar_url })
    │
    ├──→ Zustand store updates
    └──→ persist middleware saves to localStorage (key: "user")
    │
    ▼
4. window.location.href = "/feed/public" (full page load)
    │
    ▼
5. Page loads → Zustand reads from localStorage
    │
    ├──→ user is restored → Navbar renders with username
    └──→ LiveSocketContext sees user → opens WebSocket
          │
          └──→ Real-time features active
```

### Component → Backend Data Flow

```
1. PublicFeedContent calls loadMorePosts()
    │
    ▼
2. getPublicPosts({ limit: 5, offset: 10 })  ← server action
    │
    ▼
3. serverApiRequest("/posts?limit=5&offset=10")
    │
    ├─ Reads JWT from cookies
    ├─ Injects traceparent header
    └─ fetch("http://gateway:8080/posts?limit=5&offset=10")
    │
    ▼
4. Go Backend returns JSON
    │
    ▼
5. serverApiRequest returns { ok: true, data: [...] }
    │
    ▼
6. Server action returns { success: true, data: [...] }
    │
    ▼
7. setPosts((prev) => [...prev, ...newPosts])  ← React state update
    │
    ▼
8. Component re-renders with new posts
```

---

## Next Steps

Continue to [Development Workflow](./06-development-workflow.md) to learn how to run the project, add new features, and debug issues.
