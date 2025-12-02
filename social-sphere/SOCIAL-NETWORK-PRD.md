# Social Network - Frontend Product Requirements Document (PRD)

**Project Name:** Social Network Frontend  
**Technology Stack:** Next.js 15, JavaScript, Tailwind CSS  
**Backend:** Golang API (Session-based authentication, WebSocket)  
**Database:** SQLite (Backend responsibility)  
**Version:** 1.0  
**Last Updated:** November 16, 2025

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Core Principles](#core-principles)
3. [Technical Architecture](#technical-architecture)
4. [Feature Specifications](#feature-specifications)
5. [Page Structure](#page-structure)
6. [Component Architecture](#component-architecture)
7. [API Integration](#api-integration)
8. [WebSocket Events](#websocket-events)
9. [Authentication & Authorization](#authentication--authorization)
10. [Data Models](#data-models)
11. [User Flows](#user-flows)
12. [UI/UX Guidelines](#uiux-guidelines)
13. [Acceptance Criteria](#acceptance-criteria)
14. [Development Phases](#development-phases)

---

## Project Overview

### Vision
Build a modern, real-time social network frontend similar to Facebook, featuring followers, profiles, posts, groups, notifications, and chat functionality.

### Goals
- Create an intuitive, responsive social networking experience
- Implement real-time features using WebSocket
- Ensure equal priority for mobile and desktop experiences
- Pass all audit requirements for deployment
- Maintain clean, maintainable code architecture

### Success Metrics
- All audit questions pass ✅
- Real-time message delivery < 500ms
- Page load time < 2 seconds
- Mobile and desktop responsive (100%)
- Zero authentication bugs

---

## Core Principles

### 1. **Mobile & Desktop Equal Priority**
- Design components that work perfectly on both mobile (320px+) and desktop (1920px+)
- Use responsive breakpoints: `sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`, `2xl: 1536px`
- Touch-friendly buttons (min 44x44px on mobile)
- Keyboard navigation support for desktop

### 2. **Real-Time First**
- WebSocket connections for chat and notifications
- Optimistic UI updates for instant feedback
- Fallback to polling if WebSocket disconnects

### 3. **Performance**
- Lazy load images (Google Images style - dynamic sizing)
- Code splitting for heavy components
- TanStack Query for smart caching
- Debounced search inputs

### 4. **Security**
- Session-based authentication (cookies from backend)
- No sensitive data in localStorage
- CSRF protection via backend cookies
- Validate all user inputs

### 5. **User Experience**
- Clear visual feedback for all actions
- Loading states for async operations
- Error messages that help users recover
- Confirmation dialogs for destructive actions

---

## Technical Architecture

### Frontend Stack
```
Next.js 15.5 (App Router)
├── JavaScript (ES2022+)
├── Tailwind CSS (styling)
├── TanStack Query (API state management)
├── Axios (HTTP client)
├── WebSocket API (real-time)
└── React Hook Form (forms)
```

### Folder Structure
```bash
src/
├── app/                        # Next.js App Router
│   ├── (auth)/                # Auth route group
│   │   ├── login/page.js
│   │   └── register/page.js
│   │
│   ├── (main)/                # Main app route group
│   │   ├── layout.js          # Main layout (header, notifications)
│   │   ├── feed/
│   │   │   ├── page.js        # Feed router
│   │   │   ├── public/page.js # Public posts feed
│   │   │   └── friends/page.js # Friends posts feed
│   │   ├── profile/
│   │   │   └── [id]/page.js
│   │   ├── posts/
│   │   │   └── [id]/page.js
│   │   ├── groups/
│   │   │   ├── page.js
│   │   │   └── [id]/page.js
│   │   ├── messages/
│   │   │   ├── page.js
│   │   │   └── [conversationId]/page.js
│   │   └── notifications/page.js
│   │
│   ├── layout.js              # Root layout
│   ├── page.js                # Landing page
│   └── globals.css
│
├── components/
│   ├── ui/                    # Base UI components
│   ├── layout/                # Layout components
│   └── features/              # Feature components
│
├── lib/
│   ├── api/                   # API client
│   ├── auth/                  # Auth utilities
│   ├── websocket/             # WebSocket client
│   └── utils/                 # Utilities
│
├── hooks/                     # Custom hooks
├── providers/                 # React Context providers
├── config/                    # Configuration
└── middleware.js              # Auth middleware
```

---

## Feature Specifications

### Feature 1: Authentication 🔐

#### Registration
**Page:** `/register`

**Required Fields:**
- Email (email validation)
- Password (min 8 chars, hashed by backend)
- First Name (required)
- Last Name (required)
- Date of Birth (date picker, must be 13+ years old)

**Optional Fields:**
- Avatar/Image (JPEG, PNG, GIF - base64 encoded)
- Nickname (alphanumeric, 3-20 chars)
- About Me (textarea, max 500 chars)

**Validation Rules:**
```javascript
{
  email: {
    required: true,
    pattern: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i
  },
  password: {
    required: true,
    minLength: 8
  },
  firstName: {
    required: true,
    minLength: 2,
    maxLength: 50
  },
  lastName: {
    required: true,
    minLength: 2,
    maxLength: 50
  },
  dateOfBirth: {
    required: true,
    validate: (value) => {
      const age = calculateAge(value)
      return age >= 13 || 'Must be at least 13 years old'
    }
  },
  avatar: {
    required: false,
    fileTypes: ['image/jpeg', 'image/png', 'image/gif'],
    maxSize: 5 * 1024 * 1024 // 5MB
  },
  nickname: {
    required: false,
    minLength: 3,
    maxLength: 20,
    pattern: /^[a-zA-Z0-9_]+$/
  },
  aboutMe: {
    required: false,
    maxLength: 500
  }
}
```

**Flow:**
1. User fills form
2. Client-side validation
3. Convert avatar to base64 if provided
4. POST `/api/v1/auth/register`
5. Backend returns session cookie
6. Redirect to `/feed/public`

**Error Handling:**
- Email already exists → "This email is already registered"
- Validation errors → Show inline errors below fields
- Server error → "Registration failed. Please try again."

#### Login
**Page:** `/login`

**Fields:**
- Email/Username
- Password

**Flow:**
1. User enters credentials (identifier and password)
2. POST `/api/v1/auth/login`
3. Backend sets session cookie
4. Redirect to `/feed/public`
5. Middleware protects all `/feed/*`, `/profile/*`, `/groups/*`, `/messages/*` routes

**Error Handling:**
- Invalid credentials → "Invalid email or password"
- Account not found → "No account found with this email"
- Server error → "Login failed. Please try again."

#### Session Management
- Backend manages sessions via cookies
- Frontend reads session cookie to determine auth state
- Middleware redirects unauthenticated users to `/login`
- Logout clears session cookie and redirects to `/login`

**Logout:**
- POST `/api/v1/auth/logout`
- Clear any client-side state
- Redirect to `/login`

---

### Feature 2: Profile 👤

#### Profile Types
1. **Public Profile** - Visible to all users
2. **Private Profile** - Visible only to followers

#### Profile Page
**Route:** `/profile/[username]`

**Sections:**
1. **Profile Header**
   - Avatar (large)
   - Full name
   - Username (@username)
   - Nickname (if set)
   - About Me (if set)
   - Date of Birth
   - Follow/Unfollow button (if not own profile)
   - Privacy toggle (if own profile)
   - Follower count
   - Following count

2. **Profile Stats**
   - Total posts
   - Followers count
   - Following count

3. **User Posts**
   - All posts by this user
   - Reverse chronological order
   - Same post card as feed

4. **Followers/Following Lists**
   - Modal or separate section
   - Clickable to navigate to other profiles

#### Privacy Toggle (Own Profile Only)
**UI:** Toggle switch in profile header

**States:**
- 🌍 Public (Green) - "Anyone can see your profile"
- 🔒 Private (Gray) - "Only followers can see your profile"

**Confirmation:**
```
"Are you sure you want to change your profile to [PUBLIC/PRIVATE]?"
[Cancel] [Confirm]
```

**API:**
```javascript
PATCH /api/v1/users/me/privacy
Body: { isPublic: true/false }
```

#### Access Rules
| Profile Owner | Viewer Relationship | Can View? |
|--------------|---------------------|-----------|
| Self | - | ✅ Always |
| Public | Anyone | ✅ Yes |
| Private | Follower | ✅ Yes |
| Private | Non-follower | ❌ No - Show "This profile is private" |

**Private Profile Blocked View:**
```
┌─────────────────────────────┐
│  🔒                         │
│  This profile is private    │
│  Follow to see their posts  │
│  [Follow Button]            │
└─────────────────────────────┘
```

---

### Feature 3: Followers 👥

#### Follow Mechanism
**Two Types of Users:**

1. **Public User**
   - Click "Follow" → Instantly following (no request)
   - API: `POST /api/v1/users/:id/follow`

2. **Private User**
   - Click "Follow" → Send follow request
   - Button changes to "Pending"
   - Recipient sees notification
   - Recipient accepts/declines
   - API: `POST /api/v1/users/:id/follow` (creates pending request)

#### Follow Button States
```javascript
// Component states
{
  isFollowing: false,        // Not following
  isPending: false,          // Request sent, awaiting approval
  isFollowing: true,         // Currently following
}
```

**Button UI:**
| State | Button Text | Color | Action on Click |
|-------|-------------|-------|-----------------|
| Not following (public) | Follow | Blue | Follow instantly |
| Not following (private) | Follow | Blue | Send request |
| Pending | Pending | Gray | Cancel request |
| Following | Following | Green | Unfollow (with confirmation) |

#### Unfollow Confirmation
```
"Are you sure you want to unfollow @username?"
[Cancel] [Unfollow]
```

#### Follow Requests (Private Users Only)
**Notification:** "[@username] wants to follow you"

**Actions:**
- ✅ Accept → User becomes follower
- ❌ Decline → Request removed

**API:**
```javascript
POST /api/v1/follow-requests/:requestId/accept
POST /api/v1/follow-requests/:requestId/decline
GET  /api/v1/follow-requests  // Get pending requests
```

#### Followers/Following Lists
**Access via Profile:**
- Click "X Followers" or "Y Following"
- Opens modal/sheet with list
- Each item shows:
  - Avatar
  - Name
  - Username
  - Follow button (if applicable)

---

### Feature 4: Posts 📝

#### Feed Pages

**Route: `/feed/public`**
- Shows ALL public posts from ALL users
- Reverse chronological order
- Visible to everyone (even non-logged-in users? - clarify with backend)
- Includes own posts

**Route: `/feed/friends`**
- Shows posts from followed users
- Includes "almost private" and "private" posts you're allowed to see
- Reverse chronological order
- Requires authentication

**Feed Layout:**
```
┌─────────────────────────┐
│  Create Post (Sticky)   │
├─────────────────────────┤
│  Post Card              │
├─────────────────────────┤
│  Post Card              │
├─────────────────────────┤
│  Post Card              │
└─────────────────────────┘
```

#### Create Post
**Location:** Top of feed (sticky component)

**Fields:**
- Text content (required, min 1 char, max 5000 chars)
- Image/GIF upload (optional, JPEG/PNG/GIF, base64 encoded)
- Privacy selector (required)

**Privacy Options (Dropdown Multi-Select):**
1. **🌍 Public** - Everyone can see
2. **👥 Friends** - Only followers can see
3. **🔒 Private** - Choose specific followers

**Privacy Selector UI:**

**Option 1 & 2 (Public/Friends):**
```
[Dropdown: Public ▼]
```

**Option 3 (Private - Choose Followers):**
```
[Dropdown: Private ▼]
[Multi-select dropdown of followers]
✅ John Doe
✅ Jane Smith
☐ Alice Johnson
☐ Bob Williams
```

**API:**
```javascript
POST /api/v1/posts
Body: {
  content: "Post text...",
  image: "data:image/png;base64,...", // Optional
  privacy: "public" | "friends" | "private",
  allowedUsers: ["userId1", "userId2"] // Only if privacy === "private"
}
```

**Create Post Flow:**
1. User types content
2. (Optional) Upload image → Convert to base64
3. Select privacy
4. If "Private", select specific followers from multi-select
5. Click "Post"
6. Show loading state
7. Optimistic update (add post to feed immediately)
8. POST to backend
9. On success: Update with real post data
10. On error: Remove optimistic post, show error

#### Post Card
**Displays:**
- Author avatar
- Author name
- Author username
- Timestamp (relative: "2 hours ago")
- Post content (text)
- Post image (if present - Google Images style sizing)
- Privacy indicator icon (🌍/👥/🔒)
- Like button + count
- Comment button + count
- Share button (future)
- Delete button (if own post)

**Post Card Layout:**
```
┌─────────────────────────────────┐
│ [Avatar] Name @username · 2h    │
│          🔒 Private             │
├─────────────────────────────────┤
│ Post content text here...       │
│                                 │
│ [Image if present]              │
├─────────────────────────────────┤
│ 👍 12  💬 5  📤 Share  🗑️      │
└─────────────────────────────────┘
```

#### Single Post Page
**Route:** `/posts/[id]`

**Shows:**
- Full post (same as post card)
- All comments (reverse chronological)
- Comment input box

#### Comments
**Add Comment:**
- Text content (required)
- Image/GIF (optional, base64)
- API: `POST /api/v1/posts/:id/comments`

**Comment Display:**
- Author avatar (small)
- Author name
- Comment text
- Comment image (if present)
- Timestamp
- Like button (future)

**Comments List:**
```
┌─────────────────────────────┐
│ [Sm Avatar] John Doe        │
│ Comment text here...        │
│ [Image if present]          │
│ 1 hour ago                  │
├─────────────────────────────┤
│ [Sm Avatar] Jane Smith      │
│ Another comment...          │
│ 30 minutes ago              │
└─────────────────────────────┘
```

#### Image Handling (Google Images Style)
**Requirements:**
- Images fit specific size dynamically
- Maintain aspect ratio
- Click to expand (lightbox/modal)
- Lazy loading

**Implementation:**
```javascript
// Dynamic sizing based on container
<div className="relative w-full" style={{ paddingBottom: `${(height/width) * 100}%` }}>
  <img 
    src={imageUrl} 
    alt="Post image"
    className="absolute inset-0 w-full h-full object-cover rounded-lg"
    loading="lazy"
    onClick={openLightbox}
  />
</div>
```

**Lightbox (Click to Expand):**
```
Modal overlay with:
- Full-size image
- Close button (X)
- Click outside to close
- Previous/Next if multiple images (future)
```

---

### Feature 5: Groups 👨‍👩‍👧‍👦

#### Groups Page
**Route:** `/groups`

**Features:**
- Browse all groups (grid/list view)
- Browse MY GROUPS (grid/list view)
- Search bar (search by group title)
- "Create Group" button

**Group Card (Browse View):**
```
┌─────────────────────────┐
│  Group Title            │
│  Description preview... │
│  👥 12 members          │
│  [Join/View] button     │
└─────────────────────────┘
```

**Search Bar:**
```
[🔍 Search groups...                    ]
```

**API:**
```javascript
GET /api/v1/groups                    // All groups
GET /api/v1/groups?search=keyword     // Search groups
```

#### Create Group
**Modal/Page:** Trigger from "Create Group" button

**Fields:**
- Title (required, max 100 chars)
- Description (required, max 1000 chars)

**Flow:**
1. Click "Create Group"
2. Modal opens
3. Fill title & description
4. POST `/api/v1/groups`
5. Redirect to `/groups/[newGroupId]`

#### Group Detail Page
**Route:** `/groups/[id]`

**Sections:**
1. **Group Header**
   - Title
   - Description
   - Member count
   - Your role (Owner/Member/Non-member)
   - Action buttons (based on role)

2. **Group Tabs**
   - 📝 Posts (default)
   - 💬 Chat
   - 👥 Members
   - 📅 Events

3. **Group Posts Feed**
   - Create post (members only)
   - View all group posts
   - Comment on posts
   - Same as regular posts but scoped to group

4. **Group Chat Room**
   - Real-time chat (WebSocket)
   - All members can send/receive
   - Emoji support
   - Message history

5. **Members List**
   - All group members
   - Invite button (members can invite)
   - Accept/Decline requests (owner only)

6. **Events List**
   - All group events
   - Create event button (members only)
   - Vote on events

#### Group Membership States
| State | Available Actions |
|-------|------------------|
| **Non-member** | Request to Join |
| **Request Pending** | Cancel Request |
| **Member** | Leave Group, Create Posts, Invite Others |
| **Owner** | Everything + Accept/Decline Join Requests |

#### Invite to Group
**Flow:**
1. Member clicks "Invite"
2. Modal shows list of followers
3. Select followers to invite (multi-select)
4. Send invitations
5. Invited users get notification

**API:**
```javascript
POST /api/v1/groups/:id/invite
Body: { userIds: ["user1", "user2"] }
```

#### Join Group (Non-member)
**Flow:**
1. Non-member clicks "Request to Join"
2. Request sent to group owner
3. Owner gets notification
4. Owner accepts/declines
5. If accepted → User becomes member

**API:**
```javascript
POST /api/v1/groups/:id/join              // Send request
POST /api/v1/groups/:id/accept/:userId    // Accept (owner only)
POST /api/v1/groups/:id/decline/:userId   // Decline (owner only)
```

#### Group Posts
**Same as regular posts but:**
- Only visible to group members
- Created within group context
- API: `POST /api/v1/groups/:id/posts`

#### Group Chat
**Real-time chat room for all members**

**Features:**
- Real-time messages (WebSocket)
- Message history
- Emoji support
- Typing indicators (future)
- Online members indicator (future)

**WebSocket Event:**
```javascript
{
  type: 'group_message',
  payload: {
    groupId: 'group123',
    messageId: 'msg456',
    senderId: 'user789',
    senderName: 'John Doe',
    content: 'Hello everyone!',
    timestamp: '2025-11-16T10:30:00Z'
  }
}
```

#### Events
**Create Event (Members Only):**

**Fields:**
- Title (required)
- Description (required)
- Date & Time (required)
- Options (minimum 2):
  - Going
  - Not Going
  - (Can add more custom options)

**Event Card:**
```
┌─────────────────────────────┐
│  📅 Event Title              │
│  Description...              │
│  📆 Dec 25, 2025 at 6:00 PM  │
├─────────────────────────────┤
│  ✅ Going (8 people)         │
│  ❌ Not Going (3 people)     │
│  [Your Vote: Going ▼]        │
└─────────────────────────────┘
```

**Voting:**
- Members click to vote
- Can change vote
- See who voted for each option

**API:**
```javascript
POST /api/v1/groups/:id/events
Body: {
  title: "Event title",
  description: "Event description",
  dateTime: "2025-12-25T18:00:00Z",
  options: ["Going", "Not Going"]
}

POST /api/v1/events/:id/vote
Body: { option: "Going" }
```

**Event Notification:**
- All group members get notified when event is created

---

### Feature 6: Chat 💬

#### Messaging Rules
**Can message if:**
- You follow them, OR
- They follow you (at least one-way follow)

**Cannot message if:**
- No follow relationship exists

#### Messages Page
**Route:** `/messages`

**Layout (Split View):**
```
┌─────────────┬──────────────────┐
│             │                  │
│ Convers.    │   Chat Window    │
│ List        │                  │
│             │                  │
│             │                  │
└─────────────┴──────────────────┘
```

**Mobile Layout (Stacked):**
- Conversation list by default
- Click conversation → Full-screen chat
- Back button to return to list

#### Conversation List (Left Sidebar)
**Shows:**
- All conversations (reverse chronological by last message)
- Avatar
- Name
- Last message preview
- Timestamp
- Unread indicator (badge)

**Conversation Item:**
```
┌─────────────────────────────┐
│ [Avatar] John Doe      2h  │
│          Hey, how are you?  │
│          🔵 1              │
└─────────────────────────────┘
```

**Search Conversations:**
```
[🔍 Search messages...         ]
```

#### Start New Conversation
**Button:** "New Message" or "+" icon

**Flow:**
1. Click "New Message"
2. Modal shows list of users you can message (following or followers)
3. Select user
4. Opens chat window
5. OR check if conversation exists → Open existing

**API:**
```javascript
POST /api/v1/conversations
Body: { recipientId: "user123" }
Returns: { conversationId: "conv456" }
```

#### Chat Window
**Route:** `/messages/[conversationId]`

**Shows:**
- Recipient name & avatar (header)
- Message history (scrollable)
- Message input (bottom)

**Message Bubble:**
```
Own messages (right-aligned, blue):
┌─────────────────────────┐
│            Hello! How   │
│            are you?     │
│            10:30 AM     │
└─────────────────────────┘

Other's messages (left-aligned, gray):
┌─────────────────────────┐
│ I'm good, thanks!       │
│ 10:31 AM                │
└─────────────────────────┘
```

#### Send Message
**Input Area:**
```
┌─────────────────────────────┐
│ [😊] Type a message...      │
│                    [📷][➤] │
└─────────────────────────────┘
```

**Features:**
- Text input (auto-resize textarea)
- Emoji picker (click 😊)
- Image upload (click 📷)
- Send button (➤)
- Enter to send, Shift+Enter for new line

**API:**
```javascript
POST /api/v1/conversations/:id/messages
Body: {
  content: "Message text",
  image: "data:image/png;base64,..." // Optional
}
```

**Real-time Delivery:**
1. User types message
2. POST to API
3. WebSocket broadcasts to recipient
4. Recipient receives instantly (if online)
5. Message appears in chat window

**WebSocket Event:**
```javascript
{
  type: 'message',
  payload: {
    conversationId: 'conv123',
    messageId: 'msg456',
    senderId: 'user789',
    senderName: 'John Doe',
    senderAvatar: 'avatar.jpg',
    content: 'Hello!',
    image: null,
    timestamp: '2025-11-16T10:30:00Z'
  }
}
```

#### Emoji Support
**Implementation:**
- Use emoji picker library (e.g., emoji-picker-react)
- Click emoji to insert at cursor position
- Support standard Unicode emojis

#### Typing Indicator (Future - Optional)
```
[Avatar] John is typing...
```

---

### Feature 7: Notifications 🔔

#### Notification Types
1. **Follow Request** (Private profile only)
   - "[@username] wants to follow you"
   - Actions: Accept | Decline

2. **Group Invitation**
   - "[@username] invited you to join [Group Name]"
   - Actions: Accept | Decline

3. **Group Join Request** (Group owner only)
   - "[@username] wants to join [Group Name]"
   - Actions: Accept | Decline

4. **Event Created** (Group members)
   - "New event in [Group Name]: [Event Title]"
   - Action: View Event

5. **(Optional) Post Interaction** (Future)
   - "[@username] liked your post"
   - "[@username] commented on your post"

#### Notification Bell (Header)
**Location:** Top-right of every page

**UI:**
```
[🔔]  or  [🔔] 3  (with unread count)
```

**Click Behavior:**
- Opens dropdown with recent notifications (max 5)
- "View All" link → `/notifications`

**Dropdown:**
```
┌─────────────────────────────────┐
│ Notifications                   │
├─────────────────────────────────┤
│ 🔵 @john wants to follow you    │
│    [Accept] [Decline]       2h  │
├─────────────────────────────────┤
│ 🔵 New event in Group Name      │
│    [View]                  1h   │
├─────────────────────────────────┤
│ View All Notifications          │
└─────────────────────────────────┘
```

**Unread Indicator:**
- Blue dot on unread notifications
- Badge count on bell icon

#### Notifications Page
**Route:** `/notifications`

**Shows:**
- All notifications (paginated)
- Filter tabs: All | Unread | Follow Requests | Groups
- Mark all as read button

**API:**
```javascript
GET   /api/v1/notifications
PATCH /api/v1/notifications/:id/read
PATCH /api/v1/notifications/read-all
```

#### Real-time Notifications
**WebSocket Event:**
```javascript
{
  type: 'notification',
  payload: {
    id: 'notif123',
    type: 'follow_request',
    fromUserId: 'user456',
    fromUserName: 'John Doe',
    fromUserAvatar: 'avatar.jpg',
    message: '@john wants to follow you',
    actionUrl: '/profile/john',
    createdAt: '2025-11-16T10:30:00Z'
  }
}
```

**Flow:**
1. Backend sends notification via WebSocket
2. Frontend receives event
3. Update notification count (badge)
4. Show toast/banner (optional)
5. Add to notifications dropdown
6. If on notifications page → Prepend to list

---

### Feature 8: Global Search 🔍

#### Search Bar Location
**Every page (except auth pages)** has a search bar in the header

**Search Types:**
1. **Users** (primary)
   - Search by name, username, email
   - Results show: Avatar, Name, Username, Follow button

2. **Groups** (on `/groups` page only)
   - Search by group title
   - Results show: Group card

**Search Bar UI:**
```
Header (Desktop):
[Logo] [🔍 Search users...              ] [Feed] [Groups] [Messages] [🔔] [Avatar▼]

Header (Mobile):
[☰] [🔍] [🔔] [Avatar]
```

**Search Results Dropdown:**
```
┌─────────────────────────────┐
│ Users                       │
├─────────────────────────────┤
│ [Avatar] John Doe           │
│          @johndoe  [Follow] │
├─────────────────────────────┤
│ [Avatar] Jane Smith         │
│          @janesmith [Follow]│
├─────────────────────────────┤
│ View all results            │
└─────────────────────────────┘
```

**API:**
```javascript
GET /api/v1/search/users?q=keyword
GET /api/v1/search/groups?q=keyword  // On groups page
```

**Debouncing:**
- Wait 300ms after user stops typing before searching
- Show loading indicator while searching
- Clear results when search is empty

---

## Page Structure

### Authentication Pages

#### `/register`
```javascript
// app/(auth)/register/page.js
export default function RegisterPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white p-8 rounded-lg shadow">
        <h1>Create Account</h1>
        <RegisterForm />
        <p>Already have an account? <Link to="/login">Login</Link></p>
      </div>
    </div>
  )
}
```

**Components:**
- `RegisterForm` (handles all logic)
- Required/Optional field indicators
- Image upload preview
- Password strength indicator
- Form validation errors

#### `/login`
```javascript
// app/(auth)/login/page.js
export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white p-8 rounded-lg shadow">
        <h1>Welcome Back</h1>
        <LoginForm />
        <p>Don't have an account? <Link to="/register">Sign up</Link></p>
      </div>
    </div>
  )
}
```

---

### Main App Pages

#### `/` (Landing Page)
**For non-authenticated users:**
- Hero section
- Features showcase
- CTA: "Get Started" → `/register`

**For authenticated users:**
- Redirect to `/feed/public`

#### `/feed/public`
**Public Feed - All public posts**

```javascript
// app/(main)/feed/public/page.js
export default function PublicFeedPage() {
  return (
    <div className="max-w-2xl mx-auto py-4">
      <CreatePost />
      <div className="mt-4">
        <PostFeed feedType="public" />
      </div>
    </div>
  )
}
```

**Components:**
- `CreatePost` (sticky at top)
- `PostFeed` (infinite scroll)
- `PostCard` (for each post)

#### `/feed/friends`
**Friends Feed - Posts from followed users**

```javascript
// app/(main)/feed/friends/page.js
export default function FriendsFeedPage() {
  return (
    <div className="max-w-2xl mx-auto py-4">
      <CreatePost />
      <div className="mt-4">
        <PostFeed feedType="friends" />
      </div>
    </div>
  )
}
```

**Shows:**
- Almost private posts (followers only)
- Private posts (if you're in allowed list)

#### `/profile/[username]`
**User Profile Page**

```javascript
// app/(main)/profile/[username]/page.js
export default async function ProfilePage({ params }) {
  const { username } = await params
  
  // Check if can view (public or following)
  // If private and not following → Show blocked view
  
  return (
    <div className="max-w-4xl mx-auto py-4">
      <ProfileHeader user={user} />
      <ProfileTabs>
        <Tab name="Posts">
          <ProfilePosts userId={user.id} />
        </Tab>
        <Tab name="Followers">
          <FollowersList userId={user.id} />
        </Tab>
        <Tab name="Following">
          <FollowingList userId={user.id} />
        </Tab>
      </ProfileTabs>
    </div>
  )
}
```

**Components:**
- `ProfileHeader` (avatar, info, follow button, privacy toggle)
- `ProfileTabs` (Posts, Followers, Following)
- `PrivacyToggle` (own profile only)

#### `/posts/[id]`
**Single Post View**

```javascript
// app/(main)/posts/[id]/page.js
export default async function PostPage({ params }) {
  const { id } = await params
  
  return (
    <div className="max-w-2xl mx-auto py-4">
      <BackButton />
      <PostCard postId={id} expanded={true} />
      <CommentsList postId={id} />
      <CommentInput postId={id} />
    </div>
  )
}
```

#### `/groups`
**Browse Groups**

```javascript
// app/(main)/groups/page.js
export default function GroupsPage() {
  return (
    <div className="max-w-6xl mx-auto py-4">
      <div className="flex justify-between items-center mb-6">
        <h1>Groups</h1>
        <button>Create Group</button>
      </div>
      
      <SearchBar placeholder="Search groups..." />
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
        <GroupsList />
      </div>
    </div>
  )
}
```

#### `/groups/[id]`
**Group Detail Page**

```javascript
// app/(main)/groups/[id]/page.js
export default async function GroupPage({ params }) {
  const { id } = await params
  
  return (
    <div className="max-w-6xl mx-auto py-4">
      <GroupHeader groupId={id} />
      
      <Tabs>
        <Tab name="Posts">
          <GroupPosts groupId={id} />
        </Tab>
        <Tab name="Chat">
          <GroupChat groupId={id} />
        </Tab>
        <Tab name="Members">
          <GroupMembers groupId={id} />
        </Tab>
        <Tab name="Events">
          <GroupEvents groupId={id} />
        </Tab>
      </Tabs>
    </div>
  )
}
```

#### `/messages`
**Messages List + Chat**

```javascript
// app/(main)/messages/page.js
export default function MessagesPage() {
  return (
    <div className="h-[calc(100vh-64px)] flex">
      {/* Left sidebar - conversations */}
      <div className="w-80 border-r">
        <ConversationList />
      </div>
      
      {/* Right - chat window */}
      <div className="flex-1">
        <EmptyState text="Select a conversation" />
      </div>
    </div>
  )
}
```

#### `/messages/[conversationId]`
**Active Conversation**

```javascript
// app/(main)/messages/[conversationId]/page.js
export default async function ConversationPage({ params }) {
  const { conversationId } = await params
  
  return (
    <div className="h-[calc(100vh-64px)] flex">
      <div className="w-80 border-r">
        <ConversationList activeId={conversationId} />
      </div>
      
      <div className="flex-1">
        <ChatWindow conversationId={conversationId} />
      </div>
    </div>
  )
}
```

#### `/notifications`
**All Notifications**

```javascript
// app/(main)/notifications/page.js
export default function NotificationsPage() {
  return (
    <div className="max-w-2xl mx-auto py-4">
      <div className="flex justify-between items-center mb-6">
        <h1>Notifications</h1>
        <button>Mark all as read</button>
      </div>
      
      <NotificationsList />
    </div>
  )
}
```

---

## Component Architecture

### UI Components (`components/ui/`)

#### Button
```javascript
// components/ui/Button/Button.js
export function Button({ 
  children, 
  variant = 'primary', 
  size = 'md',
  onClick,
  disabled,
  ...props 
}) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={cn(
        'btn-base',
        variant === 'primary' && 'bg-blue-500 text-white hover:bg-blue-600',
        variant === 'secondary' && 'bg-gray-200 text-gray-800 hover:bg-gray-300',
        size === 'sm' && 'px-3 py-1.5 text-sm',
        size === 'md' && 'px-4 py-2',
        disabled && 'opacity-50 cursor-not-allowed'
      )}
      {...props}
    >
      {children}
    </button>
  )
}
```

**Variants:**
- `primary` - Blue (main actions)
- `secondary` - Gray (secondary actions)
- `danger` - Red (destructive actions)
- `ghost` - Transparent (subtle actions)

#### Input
```javascript
// components/ui/Input/Input.js
export function Input({ 
  label, 
  error, 
  type = 'text',
  ...props 
}) {
  return (
    <div className="mb-4">
      {label && (
        <label className="block text-sm font-medium mb-1">
          {label}
        </label>
      )}
      <input
        type={type}
        className={cn(
          'w-full px-3 py-2 border rounded-lg',
          error ? 'border-red-500' : 'border-gray-300',
          'focus:outline-none focus:ring-2 focus:ring-blue-500'
        )}
        {...props}
      />
      {error && (
        <p className="text-red-500 text-sm mt-1">{error}</p>
      )}
    </div>
  )
}
```

#### Avatar
```javascript
// components/ui/Avatar/Avatar.js
export function Avatar({ 
  src, 
  alt, 
  size = 'md',
  online = false 
}) {
  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-10 h-10',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24'
  }
  
  return (
    <div className="relative inline-block">
      <img
        src={src || '/default-avatar.png'}
        alt={alt}
        className={cn(
          'rounded-full object-cover',
          sizeClasses[size]
        )}
      />
      {online && (
        <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white rounded-full" />
      )}
    </div>
  )
}
```

#### Modal
```javascript
// components/ui/Modal/Modal.js
export function Modal({ isOpen, onClose, title, children }) {
  if (!isOpen) return null
  
  return (
    <div 
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      onClick={onClose}
    >
      <div 
        className="bg-white rounded-lg p-6 max-w-md w-full mx-4"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">{title}</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            ✕
          </button>
        </div>
        {children}
      </div>
    </div>
  )
}
```

#### Dropdown
```javascript
// components/ui/Dropdown/Dropdown.js
export function Dropdown({ 
  trigger, 
  items,
  align = 'left' 
}) {
  const [isOpen, setIsOpen] = useState(false)
  
  return (
    <div className="relative">
      <div onClick={() => setIsOpen(!isOpen)}>
        {trigger}
      </div>
      
      {isOpen && (
        <div className={cn(
          'absolute top-full mt-2 bg-white rounded-lg shadow-lg border py-2 z-50',
          align === 'left' ? 'left-0' : 'right-0'
        )}>
          {items.map((item, index) => (
            <button
              key={index}
              onClick={() => {
                item.onClick()
                setIsOpen(false)
              }}
              className="w-full text-left px-4 py-2 hover:bg-gray-100"
            >
              {item.label}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
```

---

### Layout Components (`components/layout/`)

#### Header
```javascript
// components/layout/Header/Header.js
export function Header() {
  const { user } = useAuth()
  
  return (
    <header className="sticky top-0 bg-white border-b z-40">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        {/* Logo */}
        <Link href="/feed/public" className="text-2xl font-bold">
          SocialNet
        </Link>
        
        {/* Search */}
        <GlobalSearch />
        
        {/* Navigation */}
        <nav className="flex items-center gap-4">
          <Link href="/feed/public">Feed</Link>
          <Link href="/groups">Groups</Link>
          <Link href="/messages">Messages</Link>
          
          <NotificationBell />
          
          <UserMenu user={user} />
        </nav>
      </div>
    </header>
  )
}
```

#### NotificationBell
```javascript
// components/layout/NotificationBell/NotificationBell.js
export function NotificationBell() {
  const [isOpen, setIsOpen] = useState(false)
  const { data: notifications } = useNotifications()
  
  const unreadCount = notifications?.filter(n => !n.isRead).length || 0
  
  return (
    <div className="relative">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 hover:bg-gray-100 rounded-lg"
      >
        <Bell size={20} />
        {unreadCount > 0 && (
          <span className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>
      
      {isOpen && (
        <NotificationDropdown 
          notifications={notifications} 
          onClose={() => setIsOpen(false)} 
        />
      )}
    </div>
  )
}
```

---

### Feature Components (`components/features/`)

#### CreatePost
```javascript
// components/features/posts/CreatePost/CreatePost.js
export function CreatePost() {
  const [content, setContent] = useState('')
  const [image, setImage] = useState(null)
  const [privacy, setPrivacy] = useState('public')
  const [allowedUsers, setAllowedUsers] = useState([])
  
  const createPost = useCreatePost()
  
  const handleImageUpload = (e) => {
    const file = e.target.files[0]
    if (!file) return
    
    // Convert to base64
    const reader = new FileReader()
    reader.onloadend = () => {
      setImage(reader.result)
    }
    reader.readAsDataURL(file)
  }
  
  const handleSubmit = async () => {
    await createPost.mutateAsync({
      content,
      image,
      privacy,
      allowedUsers: privacy === 'private' ? allowedUsers : []
    })
    
    // Reset form
    setContent('')
    setImage(null)
    setPrivacy('public')
    setAllowedUsers([])
  }
  
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="What's on your mind?"
        className="w-full p-3 border rounded-lg resize-none"
        rows={3}
      />
      
      {image && (
        <div className="mt-2 relative">
          <img src={image} alt="Preview" className="max-h-64 rounded-lg" />
          <button 
            onClick={() => setImage(null)}
            className="absolute top-2 right-2 bg-black/50 text-white p-1 rounded-full"
          >
            ✕
          </button>
        </div>
      )}
      
      <div className="mt-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <input
            type="file"
            id="image-upload"
            accept="image/jpeg,image/png,image/gif"
            onChange={handleImageUpload}
            className="hidden"
          />
          <label 
            htmlFor="image-upload"
            className="cursor-pointer p-2 hover:bg-gray-100 rounded-lg"
          >
            📷 Image
          </label>
          
          <PrivacySelector
            value={privacy}
            onChange={setPrivacy}
            allowedUsers={allowedUsers}
            onAllowedUsersChange={setAllowedUsers}
          />
        </div>
        
        <Button
          onClick={handleSubmit}
          disabled={!content.trim() || createPost.isPending}
        >
          {createPost.isPending ? 'Posting...' : 'Post'}
        </Button>
      </div>
    </div>
  )
}
```

#### PostCard
```javascript
// components/features/posts/PostCard/PostCard.js
export function PostCard({ post }) {
  const { user } = useAuth()
  const likePost = useLikePost()
  const deletePost = useDeletePost()
  
  const isOwn = post.authorId === user?.id
  
  return (
    <article className="bg-white p-4 rounded-lg shadow mb-4">
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <Avatar src={post.author.avatar} alt={post.author.name} />
          <div>
            <Link href={`/profile/${post.author.username}`} className="font-semibold hover:underline">
              {post.author.name}
            </Link>
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <span>@{post.author.username}</span>
              <span>·</span>
              <time>{formatRelativeTime(post.createdAt)}</time>
              <PrivacyIcon privacy={post.privacy} />
            </div>
          </div>
        </div>
        
        {isOwn && (
          <Dropdown
            trigger={<button>⋯</button>}
            items={[
              { label: 'Delete', onClick: () => deletePost.mutate(post.id) }
            ]}
          />
        )}
      </div>
      
      {/* Content */}
      <div className="mb-3">
        <p className="whitespace-pre-wrap">{post.content}</p>
      </div>
      
      {/* Image */}
      {post.image && (
        <DynamicImage src={post.image} alt="Post image" />
      )}
      
      {/* Actions */}
      <div className="flex items-center gap-6 pt-3 border-t">
        <button
          onClick={() => likePost.mutate(post.id)}
          className="flex items-center gap-2 text-gray-600 hover:text-blue-500"
        >
          <Heart className={post.isLiked ? 'fill-red-500 text-red-500' : ''} size={18} />
          <span>{post.likesCount}</span>
        </button>
        
        <Link 
          href={`/posts/${post.id}`}
          className="flex items-center gap-2 text-gray-600 hover:text-blue-500"
        >
          <MessageCircle size={18} />
          <span>{post.commentsCount}</span>
        </Link>
        
        <button className="flex items-center gap-2 text-gray-600 hover:text-blue-500">
          <Share2 size={18} />
        </button>
      </div>
    </article>
  )
}
```

---

## API Integration

### API Client Setup
```javascript
// lib/api/client.js
import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Important for cookies!
  timeout: 10000,
})

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Unauthorized - redirect to login
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
```

### API Endpoints Reference

**Note:** Backend team is still working on API. These are proposed endpoints. Update when backend provides actual API documentation.

```javascript
// lib/api/endpoints.js

// Auth
export const AUTH_ENDPOINTS = {
  register: '/api/v1/auth/register',
  login: '/api/v1/auth/login',
  logout: '/api/v1/auth/logout',
  me: '/api/v1/auth/me',
}

// Users
export const USER_ENDPOINTS = {
  list: '/api/v1/users',
  profile: (username) => `/api/v1/users/${username}`,
  updateProfile: '/api/v1/users/me',
  togglePrivacy: '/api/v1/users/me/privacy',
  followers: (userId) => `/api/v1/users/${userId}/followers`,
  following: (userId) => `/api/v1/users/${userId}/following`,
}

// Followers
export const FOLLOW_ENDPOINTS = {
  follow: (userId) => `/api/v1/users/${userId}/follow`,
  unfollow: (userId) => `/api/v1/users/${userId}/follow`,
  requests: '/api/v1/follow-requests',
  acceptRequest: (requestId) => `/api/v1/follow-requests/${requestId}/accept`,
  declineRequest: (requestId) => `/api/v1/follow-requests/${requestId}/decline`,
}

// Posts
export const POST_ENDPOINTS = {
  feed: (type) => `/api/v1/posts/${type}`, // /posts/public or /posts/friends
  create: '/api/v1/posts',
  single: (id) => `/api/v1/posts/${id}`,
  delete: (id) => `/api/v1/posts/${id}`,
  like: (id) => `/api/v1/posts/${id}/like`,
  comments: (id) => `/api/v1/posts/${id}/comments`,
  createComment: (id) => `/api/v1/posts/${id}/comments`,
}

// Groups
export const GROUP_ENDPOINTS = {
  list: '/api/v1/groups',
  search: (query) => `/api/v1/groups?search=${query}`,
  create: '/api/v1/groups',
  single: (id) => `/api/v1/groups/${id}`,
  invite: (id) => `/api/v1/groups/${id}/invite`,
  join: (id) => `/api/v1/groups/${id}/join`,
  accept: (groupId, userId) => `/api/v1/groups/${groupId}/accept/${userId}`,
  decline: (groupId, userId) => `/api/v1/groups/${groupId}/decline/${userId}`,
  members: (id) => `/api/v1/groups/${id}/members`,
  posts: (id) => `/api/v1/groups/${id}/posts`,
  createPost: (id) => `/api/v1/groups/${id}/posts`,
  events: (id) => `/api/v1/groups/${id}/events`,
  createEvent: (id) => `/api/v1/groups/${id}/events`,
  voteEvent: (eventId) => `/api/v1/events/${eventId}/vote`,
  messages: (id) => `/api/v1/groups/${id}/messages`,
  sendMessage: (id) => `/api/v1/groups/${id}/messages`,
}

// Messages
export const MESSAGE_ENDPOINTS = {
  conversations: '/api/v1/conversations',
  createConversation: '/api/v1/conversations',
  messages: (convId) => `/api/v1/conversations/${convId}/messages`,
  sendMessage: (convId) => `/api/v1/conversations/${convId}/messages`,
}

// Notifications
export const NOTIFICATION_ENDPOINTS = {
  list: '/api/v1/notifications',
  markRead: (id) => `/api/v1/notifications/${id}/read`,
  markAllRead: '/api/v1/notifications/read-all',
}

// Search
export const SEARCH_ENDPOINTS = {
  users: (query) => `/api/v1/search/users?q=${query}`,
  groups: (query) => `/api/v1/search/groups?q=${query}`,
}
```

### API Hooks (TanStack Query)

```javascript
// hooks/api/usePosts.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api/client'
import { POST_ENDPOINTS } from '@/lib/api/endpoints'

export function useFeed(type = 'public') {
  return useQuery({
    queryKey: ['posts', type],
    queryFn: async () => {
      const response = await apiClient.get(POST_ENDPOINTS.feed(type))
      return response.data.posts
    },
  })
}

export function useCreatePost() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data) => {
      const response = await apiClient.post(POST_ENDPOINTS.create, data)
      return response.data.post
    },
    onSuccess: () => {
      // Invalidate both feeds
      queryClient.invalidateQueries({ queryKey: ['posts', 'public'] })
      queryClient.invalidateQueries({ queryKey: ['posts', 'friends'] })
    },
  })
}

export function useLikePost() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (postId) => {
      await apiClient.post(POST_ENDPOINTS.like(postId))
    },
    onMutate: async (postId) => {
      // Optimistic update
      // ... (same pattern as before)
    },
  })
}
```

---

## WebSocket Events

### WebSocket Client
```javascript
// lib/websocket/client.js
export class WebSocketClient {
  constructor(url) {
    this.url = url
    this.ws = null
    this.eventHandlers = new Map()
  }
  
  connect() {
    this.ws = new WebSocket(this.url)
    
    this.ws.onopen = () => {
      console.log('WebSocket connected')
    }
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.handleEvent(data)
    }
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
    
    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      // Auto-reconnect after 3 seconds
      setTimeout(() => this.connect(), 3000)
    }
  }
  
  on(eventType, handler) {
    if (!this.eventHandlers.has(eventType)) {
      this.eventHandlers.set(eventType, new Set())
    }
    this.eventHandlers.get(eventType).add(handler)
  }
  
  off(eventType, handler) {
    this.eventHandlers.get(eventType)?.delete(handler)
  }
  
  handleEvent(event) {
    const handlers = this.eventHandlers.get(event.type)
    if (handlers) {
      handlers.forEach(handler => handler(event.payload))
    }
  }
  
  send(data) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }
}

// Singleton
let wsClient = null

export function getWebSocketClient() {
  if (!wsClient) {
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'
    wsClient = new WebSocketClient(wsUrl)
  }
  return wsClient
}
```

### Event Types

#### Private Message Event
```javascript
{
  type: 'message',
  payload: {
    conversationId: 'conv123',
    messageId: 'msg456',
    senderId: 'user789',
    senderName: 'John Doe',
    senderAvatar: '/avatars/john.jpg',
    content: 'Hello!',
    image: null,
    timestamp: '2025-11-16T10:30:00Z'
  }
}
```

#### Group Message Event
```javascript
{
  type: 'group_message',
  payload: {
    groupId: 'group123',
    messageId: 'msg789',
    senderId: 'user456',
    senderName: 'Jane Smith',
    content: 'Hi everyone!',
    timestamp: '2025-11-16T10:35:00Z'
  }
}
```

#### Notification Event
```javascript
{
  type: 'notification',
  payload: {
    id: 'notif123',
    type: 'follow_request',
    fromUserId: 'user456',
    fromUserName: 'John Doe',
    fromUserAvatar: '/avatars/john.jpg',
    message: '@john wants to follow you',
    actionUrl: '/profile/john',
    createdAt: '2025-11-16T10:30:00Z'
  }
}
```

#### Typing Indicator (Future)
```javascript
{
  type: 'typing',
  payload: {
    conversationId: 'conv123',
    userId: 'user456',
    isTyping: true
  }
}
```

---

## Data Models

### User
```javascript
{
  id: 'user123',
  email: 'john@example.com',
  firstName: 'John',
  lastName: 'Doe',
  username: 'johndoe',
  nickname: 'JD',
  dateOfBirth: '1995-05-15',
  avatar: 'data:image/jpeg;base64,...',
  aboutMe: 'Software developer...',
  isPublic: true,
  followersCount: 150,
  followingCount: 200,
  createdAt: '2025-01-01T00:00:00Z'
}
```

### Post
```javascript
{
  id: 'post123',
  authorId: 'user123',
  author: {
    id: 'user123',
    name: 'John Doe',
    username: 'johndoe',
    avatar: 'data:image/jpeg;base64,...'
  },
  content: 'This is a post...',
  image: 'data:image/png;base64,...', // Optional
  privacy: 'public' | 'friends' | 'private',
  allowedUsers: ['user456', 'user789'], // Only if privacy === 'private'
  likesCount: 25,
  commentsCount: 5,
  isLiked: false,
  createdAt: '2025-11-16T10:00:00Z'
}
```

### Comment
```javascript
{
  id: 'comment123',
  postId: 'post123',
  authorId: 'user456',
  author: {
    id: 'user456',
    name: 'Jane Smith',
    username: 'janesmith',
    avatar: 'data:image/jpeg;base64,...'
  },
  content: 'Great post!',
  image: null, // Optional
  createdAt: '2025-11-16T10:15:00Z'
}
```

### Group
```javascript
{
  id: 'group123',
  title: 'Photography Lovers',
  description: 'A group for photography enthusiasts...',
  ownerId: 'user123',
  membersCount: 50,
  isMember: true,
  role: 'owner' | 'member' | null,
  createdAt: '2025-11-01T00:00:00Z'
}
```

### Event
```javascript
{
  id: 'event123',
  groupId: 'group123',
  title: 'Photography Meetup',
  description: 'Join us for a photo walk...',
  dateTime: '2025-12-25T14:00:00Z',
  options: ['Going', 'Not Going', 'Maybe'],
  votes: {
    'Going': ['user123', 'user456'],
    'Not Going': ['user789'],
    'Maybe': []
  },
  userVote: 'Going',
  createdAt: '2025-11-16T10:00:00Z'
}
```

### Conversation
```javascript
{
  id: 'conv123',
  participantIds: ['user123', 'user456'],
  participants: [
    {
      id: 'user456',
      name: 'Jane Smith',
      username: 'janesmith',
      avatar: 'data:image/jpeg;base64,...',
      isOnline: true
    }
  ],
  lastMessage: {
    id: 'msg789',
    content: 'See you tomorrow!',
    senderId: 'user456',
    timestamp: '2025-11-16T12:00:00Z'
  },
  unreadCount: 2,
  updatedAt: '2025-11-16T12:00:00Z'
}
```

### Message
```javascript
{
  id: 'msg123',
  conversationId: 'conv123',
  senderId: 'user123',
  content: 'Hello!',
  image: null, // Optional
  isRead: true,
  createdAt: '2025-11-16T10:00:00Z'
}
```

### Notification
```javascript
{
  id: 'notif123',
  userId: 'user123', // Recipient
  type: 'follow_request' | 'group_invite' | 'group_join_request' | 'event_created',
  fromUserId: 'user456',
  fromUserName: 'Jane Smith',
  fromUserAvatar: 'data:image/jpeg;base64,...',
  message: '@janesmith wants to follow you',
  actionUrl: '/profile/janesmith',
  isRead: false,
  createdAt: '2025-11-16T10:00:00Z'
}
```

---

## Authentication & Authorization

### Session-Based Auth (Cookies)

**Backend Responsibility:**
- Create session on login
- Set HTTP-only cookie
- Validate session on each request
- Destroy session on logout

**Frontend Responsibility:**
- Send credentials to `/api/v1/auth/login`
- Store no auth data (backend handles via cookies)
- Check auth status via `/api/v1/auth/me`
- Clear any client state on logout

### Auth Provider
```javascript
// providers/AuthProvider.js
'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { apiClient } from '@/lib/api/client'
import { AUTH_ENDPOINTS } from '@/lib/api/endpoints'

const AuthContext = createContext(undefined)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()
  
  // Check if user is authenticated on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await apiClient.get(AUTH_ENDPOINTS.me)
        setUser(response.data.user)
      } catch (error) {
        setUser(null)
      } finally {
        setIsLoading(false)
      }
    }
    
    checkAuth()
  }, [])
  
  const login = async (email, password) => {
    const response = await apiClient.post(AUTH_ENDPOINTS.login, {
      email,
      password,
    })
    
    setUser(response.data.user)
    router.push('/feed/public')
  }
  
  const logout = async () => {
    await apiClient.post(AUTH_ENDPOINTS.logout)
    setUser(null)
    router.push('/login')
  }
  
  const register = async (data) => {
    const response = await apiClient.post(AUTH_ENDPOINTS.register, data)
    setUser(response.data.user)
    router.push('/feed/public')
  }
  
  return (
    <AuthContext.Provider value={{ user, isLoading, isAuthenticated: !!user, login, logout, register }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
```

### Middleware (Route Protection)
```javascript
// middleware.js
import { NextResponse } from 'next/server'

export function middleware(request) {
  // Check for session cookie (set by backend)
  const hasSession = request.cookies.has('session') // Adjust cookie name
  
  const { pathname } = request.nextUrl
  
  // Public routes (no auth needed)
  const publicRoutes = ['/', '/login', '/register']
  const isPublicRoute = publicRoutes.includes(pathname)
  
  // Protected routes (auth required)
  const isProtectedRoute = pathname.startsWith('/feed') ||
                          pathname.startsWith('/profile') ||
                          pathname.startsWith('/groups') ||
                          pathname.startsWith('/messages') ||
                          pathname.startsWith('/notifications')
  
  // Redirect to login if accessing protected route without session
  if (isProtectedRoute && !hasSession) {
    const loginUrl = new URL('/login', request.url)
    loginUrl.searchParams.set('redirect', pathname)
    return NextResponse.redirect(loginUrl)
  }
  
  // Redirect to feed if accessing auth pages with session
  if ((pathname === '/login' || pathname === '/register') && hasSession) {
    return NextResponse.redirect(new URL('/feed/public', request.url))
  }
  
  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
}
```

---

## UI/UX Guidelines

### Responsive Design

**Breakpoints:**
```javascript
{
  sm: '640px',   // Mobile landscape
  md: '768px',   // Tablet
  lg: '1024px',  // Desktop
  xl: '1280px',  // Large desktop
  '2xl': '1536px' // Extra large
}
```

**Mobile-First Approach:**
```javascript
// Base styles for mobile
className="text-sm p-2"

// Add tablet styles
className="text-sm md:text-base p-2 md:p-4"

// Add desktop styles
className="text-sm md:text-base lg:text-lg p-2 md:p-4 lg:p-6"
```

**Touch Targets:**
- Minimum 44x44px for clickable elements on mobile
- Spacing between interactive elements (min 8px)

### Loading States

**Button Loading:**
```javascript
<Button disabled={isLoading}>
  {isLoading ? 'Loading...' : 'Submit'}
</Button>
```

**Skeleton Loaders:**
```javascript
// Post skeleton
<div className="animate-pulse">
  <div className="flex items-center gap-3 mb-3">
    <div className="w-10 h-10 bg-gray-200 rounded-full" />
    <div className="flex-1">
      <div className="h-4 bg-gray-200 rounded w-1/3 mb-2" />
      <div className="h-3 bg-gray-200 rounded w-1/4" />
    </div>
  </div>
  <div className="h-20 bg-gray-200 rounded" />
</div>
```

**Infinite Scroll Loader:**
```javascript
<div className="flex justify-center py-4">
  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500" />
</div>
```

### Error States

**Inline Field Errors:**
```javascript
<div>
  <input className={error ? 'border-red-500' : 'border-gray-300'} />
  {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
</div>
```

**Toast Notifications:**
```javascript
// Success
<Toast type="success" message="Post created successfully!" />

// Error
<Toast type="error" message="Failed to create post. Please try again." />

// Info
<Toast type="info" message="New messages available" />
```

**Empty States:**
```javascript
<div className="text-center py-12">
  <div className="text-6xl mb-4">📭</div>
  <h3 className="text-xl font-semibold mb-2">No messages yet</h3>
  <p className="text-gray-500 mb-4">Start a conversation with your followers</p>
  <Button>New Message</Button>
</div>
```

### Confirmation Dialogs

**Destructive Actions:**
```javascript
// Unfollow confirmation
<Modal isOpen={showConfirm} onClose={() => setShowConfirm(false)}>
  <h2>Unfollow @{username}?</h2>
  <p>You will stop seeing their posts in your feed.</p>
  <div className="flex gap-3 mt-4">
    <Button variant="secondary" onClick={() => setShowConfirm(false)}>
      Cancel
    </Button>
    <Button variant="danger" onClick={handleUnfollow}>
      Unfollow
    </Button>
  </div>
</Modal>
```

### Accessibility

**Semantic HTML:**
```javascript
<nav role="navigation">
  <ul role="list">
    <li role="listitem">
      <a href="/feed">Feed</a>
    </li>
  </ul>
</nav>
```

**ARIA Labels:**
```javascript
<button aria-label="Close modal">✕</button>
<input aria-describedby="email-error" />
<div id="email-error" role="alert">Invalid email</div>
```

**Keyboard Navigation:**
- All interactive elements accessible via Tab
- Enter/Space to activate buttons
- Escape to close modals
- Arrow keys for dropdown navigation

---

## Acceptance Criteria

### Authentication ✅

- [ ] User can register with all required fields
- [ ] Optional fields (avatar, nickname, about me) can be skipped
- [ ] Email validation works correctly
- [ ] Password meets minimum requirements
- [ ] Age validation (13+ years old)
- [ ] User can login with correct credentials
- [ ] Login fails with incorrect credentials
- [ ] User stays logged in until logout
- [ ] Session persists across browser refresh
- [ ] Logout clears session and redirects to login
- [ ] Cannot register with existing email
- [ ] Protected routes redirect to login when not authenticated
- [ ] Auth pages redirect to feed when already authenticated

### Profile ✅

- [ ] Profile displays all user information (except password)
- [ ] Public profiles visible to everyone
- [ ] Private profiles only visible to followers
- [ ] Non-followers see "This profile is private" for private profiles
- [ ] Own profile always visible
- [ ] Can toggle between public and private
- [ ] Confirmation dialog shown when changing privacy
- [ ] Profile shows all posts by user
- [ ] Profile shows followers count
- [ ] Profile shows following count
- [ ] Can view followers list
- [ ] Can view following list

### Followers ✅

- [ ] Can follow public users instantly
- [ ] Following private user sends request (shows "Pending")
- [ ] Private user receives follow request notification
- [ ] Can accept follow request
- [ ] Can decline follow request
- [ ] Accepted request makes user a follower
- [ ] Declined request removes the request
- [ ] Can unfollow users
- [ ] Unfollow shows confirmation dialog
- [ ] Follow button states correct (Follow, Pending, Following)

### Posts ✅

- [ ] Can create post with text
- [ ] Can add image/GIF to post (JPEG, PNG, GIF)
- [ ] Can select privacy (Public, Friends, Private)
- [ ] Private posts allow selecting specific followers
- [ ] Posts appear in correct feed (public or friends)
- [ ] Public feed shows all public posts
- [ ] Friends feed shows posts from followed users
- [ ] Post displays author info, content, image, privacy icon
- [ ] Can like posts
- [ ] Like count updates immediately (optimistic)
- [ ] Can comment on posts
- [ ] Can add image to comments
- [ ] Comments display correctly
- [ ] Single post page shows post + all comments
- [ ] Can delete own posts
- [ ] Delete shows confirmation dialog
- [ ] Images use Google Images style sizing
- [ ] Click image to expand (lightbox)
- [ ] Images lazy load

### Groups ✅

- [ ] Can browse all groups
- [ ] Can search groups by title
- [ ] Can create group with title and description
- [ ] Can invite followers to group
- [ ] Invited users receive notification
- [ ] Can accept/decline group invitation
- [ ] Non-members can request to join
- [ ] Group owner receives join request notification
- [ ] Group owner can accept/decline join requests
- [ ] Members can create posts in group
- [ ] Group posts only visible to members
- [ ] Members can comment on group posts
- [ ] Members can invite other users
- [ ] Can create events in group
- [ ] Event has title, description, date/time, options
- [ ] Members receive event creation notification
- [ ] Can vote on event options
- [ ] Can change vote
- [ ] Event shows vote counts
- [ ] Group has real-time chat
- [ ] All members receive group messages
- [ ] Can send emojis in group chat

### Chat ✅

- [ ] Can only message users with follow relationship
- [ ] Cannot message users with no follow connection
- [ ] Conversation list shows all conversations
- [ ] Shows last message preview
- [ ] Shows unread count
- [ ] Messages delivered in real-time via WebSocket
- [ ] Can send text messages
- [ ] Can send emojis
- [ ] Can attach images to messages
- [ ] Chat window shows message history
- [ ] Own messages appear on right (blue)
- [ ] Other's messages appear on left (gray)
- [ ] Messages show timestamp
- [ ] Can create new conversation
- [ ] Search conversations works
- [ ] Typing in message input doesn't lag

### Notifications ✅

- [ ] Notification bell visible on every page
- [ ] Bell shows unread count
- [ ] Click bell opens dropdown with recent notifications
- [ ] Real-time notifications via WebSocket
- [ ] Follow request notification (private profile)
- [ ] Group invitation notification
- [ ] Group join request notification (owner only)
- [ ] Event created notification (group members)
- [ ] Can accept/decline from notification
- [ ] Notifications page shows all notifications
- [ ] Can mark notification as read
- [ ] Can mark all as read
- [ ] Unread notifications have blue dot

### Global Search ✅

- [ ] Search bar visible on all main pages
- [ ] Can search users by name, username
- [ ] Search results show in dropdown
- [ ] Debounced search (300ms delay)
- [ ] Click result navigates to profile
- [ ] Can search groups on /groups page
- [ ] Group search shows group cards

### Docker ✅

- [ ] Frontend runs in Docker container
- [ ] Can access app via browser after docker run
- [ ] Container is not empty (non-zero size)
- [ ] `docker ps -a` shows frontend container
- [ ] Application works correctly in container

### General ✅

- [ ] All pages responsive (mobile + desktop)
- [ ] No console errors
- [ ] Loading states for all async operations
- [ ] Error messages are helpful
- [ ] Confirmation dialogs for destructive actions
- [ ] Images load correctly
- [ ] WebSocket connects and reconnects
- [ ] Session persists across page refresh
- [ ] All forms validate inputs
- [ ] Optimistic UI updates work correctly

---

## Development Phases

### Phase 1: Foundation (Week 1)
**Goal:** Setup + Authentication

- [ ] Project setup (Next.js, Tailwind, dependencies)
- [ ] Folder structure
- [ ] API client setup
- [ ] Auth provider
- [ ] Middleware
- [ ] Login page + form
- [ ] Register page + form
- [ ] Session management
- [ ] Base UI components (Button, Input, Avatar)
- [ ] Layout components (Header, Footer)

### Phase 2: Core Features (Week 2-3)
**Goal:** Posts, Profile, Followers

- [ ] Feed pages (public + friends)
- [ ] Create post component
- [ ] Post card component
- [ ] Image upload + base64 encoding
- [ ] Privacy selector
- [ ] Post feed with infinite scroll
- [ ] Single post page
- [ ] Comments
- [ ] Profile page
- [ ] Profile privacy toggle
- [ ] Follow/unfollow functionality
- [ ] Follow requests
- [ ] Followers/Following lists

### Phase 3: Groups (Week 4)
**Goal:** Group functionality

- [ ] Groups browse page
- [ ] Group search
- [ ] Create group
- [ ] Group detail page
- [ ] Group invitations
- [ ] Join requests
- [ ] Group posts
- [ ] Group members
- [ ] Events creation
- [ ] Event voting
- [ ] Group chat (real-time)

### Phase 4: Real-Time (Week 5)
**Goal:** WebSocket + Messaging

- [ ] WebSocket client
- [ ] WebSocket provider
- [ ] Messages page
- [ ] Conversation list
- [ ] Chat window
- [ ] Send messages
- [ ] Real-time message delivery
- [ ] Emoji support
- [ ] Image messages
- [ ] Notifications system
- [ ] Notification bell
- [ ] Real-time notifications

### Phase 5: Polish (Week 6)
**Goal:** UX improvements + Testing

- [ ] Global search
- [ ] Loading states everywhere
- [ ] Error handling
- [ ] Confirmation dialogs
- [ ] Empty states
- [ ] Responsive design review
- [ ] Performance optimization
- [ ] Image lazy loading
- [ ] Code splitting
- [ ] Accessibility audit
- [ ] Testing (manual + automated)

### Phase 6: Docker + Deployment (Week 7)
**Goal:** Containerization + Deploy

- [ ] Dockerfile for frontend
- [ ] Docker compose (if needed)
- [ ] Environment variables
- [ ] Build optimization
- [ ] Deployment setup
- [ ] Final testing
- [ ] Audit checklist verification

---

## Environment Variables

```bash
# .env.local (Development)
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# .env.production (Production)
NEXT_PUBLIC_API_URL=https://api.socialnetwork.com
NEXT_PUBLIC_WS_URL=wss://api.socialnetwork.com/ws
```

---

## Docker Configuration

### Dockerfile
```dockerfile
# frontend/Dockerfile
FROM node:20-alpine AS base

# Dependencies
FROM base AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

# Builder
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Runner
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]
```

### Build Script
```bash
#!/bin/bash
# build.sh

echo "Building frontend Docker image..."
docker build -t social-network-frontend .

echo "Build complete!"
echo "To run: docker run -p 3000:3000 social-network-frontend"
```

---

## Notes for Backend Team

**API Requirements:**

1. **Session cookies must be:**
   - HTTP-only
   - Secure (in production)
   - SameSite=Lax or Strict
   - Include credentials in CORS

2. **CORS Configuration:**
   ```
   Access-Control-Allow-Origin: http://localhost:3000 (dev) / https://app.socialnetwork.com (prod)
   Access-Control-Allow-Credentials: true
   Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE
   Access-Control-Allow-Headers: Content-Type
   ```

3. **Image Handling:**
   - Frontend sends images as base64 strings
   - Backend stores them locally
   - Backend returns image URLs or base64 in responses

4. **WebSocket:**
   - Authenticate connection via session cookie or token in URL
   - Broadcast messages to specific users/groups
   - Handle reconnections gracefully

5. **API Response Format:**
   ```json
   {
     "success": true,
     "data": { ... },
     "message": "Operation successful"
   }
   ```

   Error format:
   ```json
   {
     "success": false,
     "error": "Error message",
     "code": "ERROR_CODE"
   }
   ```

---

## Success Criteria

This project is considered **successful** when:

1. ✅ All audit questions pass
2. ✅ Real-time features work reliably
3. ✅ No authentication bugs
4. ✅ Responsive on mobile and desktop
5. ✅ All user flows complete without errors
6. ✅ Docker container runs correctly
7. ✅ Code is clean and maintainable
8. ✅ Team can easily understand and extend

---

**End of PRD**

*This document is the single source of truth for the Social Network frontend. All features, components, and requirements are defined here. Update this document when requirements change.*
