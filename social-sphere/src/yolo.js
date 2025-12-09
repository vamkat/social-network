// import { NextResponse } from 'next/server';

// // Define route categories
// const AUTH_ROUTES = ['/login', '/registration'];
// const PUBLIC_ROUTES = ['/', '/about', '/terms', '/privacy'];
// const PROTECTED_ROUTES = ['/feed', '/profile', '/groups', '/messages', '/notifications', '/settings'];

// export function proxy(request) {
//     const { pathname } = request.nextUrl;

//     // Get JWT cookie
//     const token = request.cookies.get('jwt')?.value;
//     const isAuthenticated = !!token;

//     // Check if current path matches any route pattern
//     const isAuthRoute = AUTH_ROUTES.some(route => pathname.startsWith(route));
//     const isPublicRoute = PUBLIC_ROUTES.some(route => pathname === route);
//     const isProtectedRoute = PROTECTED_ROUTES.some(route => pathname.startsWith(route));

//     // Redirect authenticated users away from auth pages
//     if (isAuthenticated && isAuthRoute) {
//         return NextResponse.redirect(new URL('/feed/public', request.url));
//     }

//     // Redirect unauthenticated users to login for protected routes
//     if (!isAuthenticated && isProtectedRoute) {
//         const loginUrl = new URL('/login', request.url);
//         // Store the original URL to redirect back after login
//         loginUrl.searchParams.set('redirect', pathname);
//         return NextResponse.redirect(loginUrl);
//     }

//     // Allow the request to continue
//     return NextResponse.next();
// }

// // Configure which routes the middleware should run on
// export const config = {
//     matcher: [
//         /*
//          * Match all request paths except:
//          * - api routes (already handled by API route handlers)
//          * - _next/static (static files)
//          * - _next/image (image optimization files)
//          * - favicon.ico (favicon file)
//          * - public files (public folder)
//          */
//         '/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
//     ],
// };
