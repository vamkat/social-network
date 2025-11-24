Authentication state is managed globally via an AuthProvider.

All routes inside app/(main)/ are protected (redirect to /login if not authenticated).

The Navbar is shared across all main pages and now uses real user data from the backend instead of mock data.


File: src/providers/AuthProvider.js: 
    On mount, calls GET /api/v1/auth/me via apiClient to:

        Check if a valid session exists (cookie-based).

        Populate user in context.

    Exposes:

        user – current logged-in user (or null).

        isLoading – true while we are checking /me.

        isAuthenticated – !!user.

        logout() – calls backend logout endpoint and redirects to /login.

File: src/lib/api/client.js:

    Wraps Axios with a base URL (NEXT_PUBLIC_API_URL or http://localhost:8080).

    Enables withCredentials: true so session cookies are sent with each request.

    Has a response interceptor:

        If response status is 401, it redirects to /login (client-side) unless already on /login.

    This gives us a central place to make API calls without repeating config.

File: src/app/(main)/layout.js:
    Uses useAuth() to read isAuthenticated + isLoading.

    While isLoading is true → shows a simple spinner.

    When auth check is done:

        If not authenticated → redirects to /login via useRouter() and renders null.

        If authenticated → renders:

            <Navbar />

            A centered <main> container for the current page.

File: src/components/layout/navbar.js:
    
    Uses useAuth() to:

        Display the current user’s username and avatar (or a default icon).

        Build the link to the profile: /profile/${user.username}.

        Call logout() on sign out.

    Contains:

        Logo linking to /feed/public.

        Navigation items:

            /feed/public (Home)

            /feed/friends (Friends)

            /groups

            /messages

        Active state based on usePathname().

        Notification bell (placeholder, UI only for now).

        User dropdown (desktop) with:

            Profile

            Settings

            Sign Out (calls logout()).

        Mobile menu with the same links + sign out.