# React Fundamentals

> **Prerequisites:** HTML, CSS, and JavaScript knowledge. No React experience needed.
>
> This guide uses real code from the SocialSphere project to teach React concepts.

---

## What is React?

React is a JavaScript library for building user interfaces. Instead of manually updating the DOM (like `document.getElementById(...).innerHTML = ...`), you describe **what** the UI should look like and React figures out **how** to update it.

Key ideas:

- **Declarative** — You describe the desired UI state, not the steps to get there
- **Component-based** — The UI is split into small, reusable pieces called components
- **Reactive** — When data changes, React automatically re-renders only the parts of the UI that need updating

---

## JSX Syntax

React uses **JSX**, a syntax extension that lets you write HTML-like code inside JavaScript. It looks like HTML but has a few differences:

| HTML | JSX | Why |
|------|-----|-----|
| `class="btn"` | `className="btn"` | `class` is a reserved word in JavaScript |
| `for="email"` | `htmlFor="email"` | `for` is a reserved word in JavaScript |
| `<img>` | `<img />` | All tags must be closed in JSX |
| `<br>` | `<br />` | Self-closing tags need the `/` |
| `style="color: red"` | `style={{ color: "red" }}` | Styles are JavaScript objects |

You can embed any JavaScript expression inside `{}`:

```jsx
const name = "SocialSphere";
const count = 5;

// Inside JSX:
<h1>{name}</h1>              // renders: SocialSphere
<p>{count * 2} posts</p>     // renders: 10 posts
<p>{count > 0 ? "Yes" : "No"}</p>  // renders: Yes
```

---

## Components

A React component is a **function that returns JSX**. Components are always named in **PascalCase** (e.g., `LoginForm`, not `loginForm`).

Here's a simple component from our project:

<!-- src/components/layout/Container.js -->
```jsx
export default function Container({ children, size = "default", className = "" }) {
    const sizeClasses = {
        narrow: "max-w-md",
        default: "max-w-3xl",
        wide: "max-w-7xl",
        full: "max-w-full",
    };

    return (
        <div className={`mx-auto w-full px-4 sm:px-6 lg:px-0 ${sizeClasses[size]} ${className}`}>
            {children}
        </div>
    );
}
```

And an even simpler one:

<!-- src/components/ui/Tooltip.js -->
```jsx
export default function Tooltip({ content, active = true, children }) {
    if (active === false) {
        return (
            <div className="group/tooltip relative inline-flex">
                {children}
            </div>
        );
    }

    return (
        <div className="group/tooltip relative inline-flex">
            {children}
            <div className="absolute top-full left-1/2 -translate-x-1/2 mt-2 px-2.5 py-1 ...">
                {content}
            </div>
        </div>
    );
}
```

Components are used like HTML tags:

```jsx
<Container size="narrow">
    <h1>Hello</h1>
</Container>

<Tooltip content="Click me">
    <button>Hover me</button>
</Tooltip>
```

---

## Props

Props (short for "properties") are the arguments you pass to a component. They work like function parameters.

<!-- src/components/ui/Modal.js -->
```jsx
export default function Modal({
    isOpen,                      // required: boolean
    onClose,                     // required: function
    title,                       // required: string
    description,                 // optional: string
    children,                    // optional: anything nested inside <Modal>...</Modal>
    footer,                      // optional: JSX
    onConfirm,                   // optional: function
    confirmText = "Confirm",     // optional with default value
    cancelText = "Cancel",       // optional with default value
    isLoading = false,           // optional with default value
    showCloseButton = true       // optional with default value
}) {
    // ... component body
}
```

Key patterns:

- **Destructuring** `{ isOpen, onClose, title }` — extracts individual props from the props object
- **Default values** `confirmText = "Confirm"` — used when the prop isn't provided
- **`children`** — a special prop that contains whatever is nested between the opening and closing tags

Using this component:

```jsx
<Modal
    isOpen={showModal}
    onClose={() => setShowModal(false)}
    title="Delete Post?"
    description="This action cannot be undone."
    onConfirm={handleDelete}
    confirmText="Delete"
>
    <p>Are you sure?</p>
</Modal>
```

---

## State (`useState`)

State is data that can **change over time** and trigger a re-render when it does. The `useState` hook gives you a value and a function to update it.

<!-- src/components/forms/LoginForm.js -->
```jsx
import { useState } from "react";

export default function LoginForm() {
    // Declare state variables with their initial values
    const [isLoading, setIsLoading] = useState(false);    // boolean
    const [error, setError] = useState("");                // string
    const [showPassword, setShowPassword] = useState(false); // boolean

    // ... rest of component
}
```

The pattern is always: `const [value, setValue] = useState(initialValue)`

- `value` — the current state value
- `setValue` — function to update it (triggers a re-render)
- `initialValue` — what it starts as

When you call the setter, React re-renders the component with the new value:

```jsx
// Toggle password visibility
<button onClick={() => setShowPassword(!showPassword)}>
    {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
</button>

// The input type changes based on state
<input type={showPassword ? "text" : "password"} />
```

---

## Effects (`useEffect`)

`useEffect` lets you run **side effects** — code that interacts with the outside world (DOM manipulation, timers, network requests, etc.). It runs after the component renders.

```jsx
useEffect(() => {
    // This code runs after render
    // ...

    return () => {
        // This cleanup code runs before the next effect or when component unmounts
    };
}, [dependencies]); // Only re-run when these values change
```

**Example 1: Preventing body scroll when a modal is open**

<!-- src/components/ui/Modal.js -->
```jsx
useEffect(() => {
    if (isOpen) {
        document.body.style.overflow = "hidden";  // Prevent scrolling
    } else {
        document.body.style.overflow = "unset";   // Restore scrolling
    }
    return () => {
        document.body.style.overflow = "unset";   // Cleanup on unmount
    };
}, [isOpen]); // Re-run when isOpen changes
```

**Example 2: IntersectionObserver for infinite scroll**

<!-- src/components/feed/PublicFeedContent.js -->
```jsx
useEffect(() => {
    // Create an observer that watches when an element enters the viewport
    const observer = new IntersectionObserver(
        (entries) => {
            if (entries[0].isIntersecting && hasMore && !loading) {
                loadMorePosts();
            }
        },
        { threshold: 0.1 }
    );

    // Start observing our target element
    if (observerTarget.current) {
        observer.observe(observerTarget.current);
    }

    // Cleanup: stop observing when component unmounts or dependencies change
    return () => {
        if (observerTarget.current) {
            observer.unobserve(observerTarget.current);
        }
    };
}, [loadMorePosts, hasMore, loading]);
```

**Dependency array rules:**

| Dependency array | When the effect runs |
|-----------------|---------------------|
| `[a, b]` | When `a` or `b` changes |
| `[]` | Only once, after the first render |
| *(omitted)* | After every render (rarely wanted) |

---

## Event Handling

React handles events with camelCase attributes: `onClick`, `onSubmit`, `onChange`, etc.

<!-- src/components/forms/LoginForm.js -->
```jsx
async function handleSubmit(event) {
    event.preventDefault();   // Stop the browser's default form submission
    setIsLoading(true);
    setError("");

    const formData = new FormData(event.currentTarget);
    const email = formData.get("email");
    const password = formData.get("password");

    try {
        const resp = await login({ email, password });
        if (!resp.success || resp.error) {
            setError(resp.error || "Invalid credentials");
            setIsLoading(false);
            return;
        }
        // Success — store user and redirect
        setUser({ id: resp.user_id, username: resp.username, avatar_url: resp.avatar_url || "" });
        window.location.href = "/feed/public";
    } catch (error) {
        setError("An unexpected error occurred");
        setIsLoading(false);
    }
}

return (
    <form onSubmit={handleSubmit}>
        <input name="email" type="email" required disabled={isLoading} />
        <input name="password" type="password" required disabled={isLoading} />
        <button type="submit" disabled={isLoading}>
            {isLoading ? <LoadingThreeDotsJumping /> : "Sign In"}
        </button>
    </form>
);
```

Common event patterns:

- `onClick={() => setShowPassword(!showPassword)}` — inline handler for simple toggles
- `onSubmit={handleSubmit}` — separate function for complex logic
- `onChange={(e) => setSearchQuery(e.target.value)}` — reading input values
- `event.preventDefault()` — stop browser default behavior

---

## Conditional Rendering

React has several ways to show or hide UI based on conditions.

**Early return — render nothing if a condition isn't met:**

<!-- src/components/ui/Modal.js -->
```jsx
export default function Modal({ isOpen, ... }) {
    // If not open, render absolutely nothing
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 ...">
            {/* Modal content */}
        </div>
    );
}
```

**`&&` operator — show something only when a condition is true:**

<!-- src/components/forms/LoginForm.js -->
```jsx
{/* Only show the error div when there IS an error */}
{error && (
    <div className="form-error animate-fade-in mt-6 text-center pt-5">
        {error}
    </div>
)}
```

**Ternary `? :` — choose between two things:**

```jsx
{/* Show loading dots OR "Sign In" text based on isLoading state */}
<button type="submit" disabled={isLoading}>
    {isLoading ? <LoadingThreeDotsJumping /> : "Sign In"}
</button>
```

---

## Lists and Keys

To render a list of items, use `.map()` to transform an array of data into an array of JSX elements. Each element needs a unique `key` prop so React can efficiently track which items changed.

<!-- src/components/feed/PublicFeedContent.js -->
```jsx
{posts.map((post, index) => (
    <motion.div
        key={post.post_id + index}   // Unique key for each item
        initial={{ opacity: 0, scale: 0.8 }}
        animate={{ opacity: 1, scale: 1 }}
    >
        <PostCard
            post={post}
            onDelete={(postId) => setPosts(prev => prev.filter(p => p.post_id !== postId))}
        />
    </motion.div>
))}
```

Why `key` matters:
- Without keys, React can't tell which item changed, so it re-renders everything
- With keys, React only updates the items that actually changed
- Use a unique identifier from your data (like `post_id`), not the array index if possible

---

## Refs (`useRef`)

`useRef` gives you a way to hold a value that **persists across renders** but **doesn't cause re-renders** when it changes. It's commonly used for two things:

1. **Accessing DOM elements directly**
2. **Storing mutable values that shouldn't trigger re-renders**

<!-- src/components/feed/PublicFeedContent.js -->
```jsx
const observerTarget = useRef(null);  // Will hold a reference to a DOM element

// Later in JSX, attach the ref to an element
<div ref={observerTarget} className="flex justify-center py-8">
    {loading && <div>Loading more posts...</div>}
</div>

// Now observerTarget.current is the actual DOM element
// Used by IntersectionObserver to detect when it's visible
```

The key difference from state:
- `useState` — changing it triggers a re-render
- `useRef` — changing `.current` does NOT trigger a re-render

---

## `useCallback`

`useCallback` memoizes a function so it keeps the same reference between renders. This is important when you pass functions to `useEffect` dependency arrays or to child components.

<!-- src/components/feed/PublicFeedContent.js -->
```jsx
const loadMorePosts = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    try {
        const result = await getPublicPosts({ limit: 5, offset });
        const newPosts = result.success ? result.data : [];

        if (newPosts && newPosts.length > 0) {
            setPosts((prevPosts) => [...prevPosts, ...newPosts]);
            setOffset((prevOffset) => prevOffset + 5);
            if (newPosts.length < 5) {
                setHasMore(false);
            }
        } else {
            setHasMore(false);
        }
    } catch (error) {
        return;
    } finally {
        setLoading(false);
    }
}, [offset, loading, hasMore]);
```

Without `useCallback`, a new function would be created every render. This would cause the `useEffect` that depends on `loadMorePosts` to re-run unnecessarily, potentially creating infinite loops.

---

## Tailwind CSS

This project uses **Tailwind CSS** — a utility-first CSS framework. Instead of writing CSS classes in separate files, you apply small utility classes directly in JSX.

```jsx
// Traditional CSS approach:
// .card { padding: 1rem; margin-top: 0.5rem; border-radius: 0.5rem; }
// <div className="card">

// Tailwind approach — classes describe the styles directly:
<div className="p-4 mt-2 rounded-lg">
```

Common Tailwind patterns in this project:

| Class | Meaning |
|-------|---------|
| `p-4` | padding: 1rem |
| `mt-2` | margin-top: 0.5rem |
| `flex items-center gap-2` | flexbox row, vertically centered, gap between items |
| `text-sm font-medium` | small text, medium weight |
| `rounded-full` | fully rounded (pill shape) |
| `w-full` | width: 100% |
| `hidden md:flex` | hidden on mobile, flex on medium+ screens |

### Custom Classes

The project defines reusable custom classes in `src/app/globals.css` using Tailwind's `@apply` directive. These bundle common patterns:

```css
/* Form inputs — consistent rounded pill style */
.form-input {
  @apply text-sm px-4 rounded-3xl border border-(--border) hover:border-foreground
         focus:border-(--accent) py-3 w-full bg-(--muted)/7 text-foreground
         placeholder:text-(--muted)/50 focus:outline-none transition-colors;
}

/* Primary button — purple, rounded */
.btn-primary {
  @apply bg-(--accent) text-white hover:bg-(--accent-hover) rounded-full
         focus:outline-none focus:bg-(--accent-hover) text-sm;
}

/* Medium heading */
.heading-md {
  @apply text-4xl font-bold tracking-tight text-foreground;
}
```

### Theme Variables

The project uses CSS variables for theming. Light and dark modes are controlled via the `.dark` class on `<html>`:

```css
:root {
  --accent: #a855f7;      /* Purple accent color */
  --foreground: #0a0a0a;  /* Text color */
  --muted: #6b7280;       /* Secondary text */
  --background: #ffffff;  /* Background color */
  --border: #e5e5e5;      /* Border color */
}

.dark {
  --foreground: #f5f5f5;
  --background: #0a0a0a;
  --border: #262626;
  /* ... */
}
```

In Tailwind classes, these variables are referenced as `text-(--accent)`, `bg-(--background)`, `border-(--border)`, etc.

---

## Next Steps

Now that you understand React fundamentals, continue to [Next.js Fundamentals](./02-nextjs-fundamentals.md) to learn how the framework builds on React with routing, server components, and server actions.
