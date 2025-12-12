import NextAuth from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"

export const authOptions = {
    secret: process.env.NEXTAUTH_SECRET || "super-secret-secret-key-change-me",
    providers: [
        CredentialsProvider({
            name: "Credentials",
            credentials: {
                userId: { label: "User ID", type: "text" }
            },
            async authorize(credentials, req) {
                try {
                    const apiBase = process.env.API_BASE || "http://localhost:8081";

                    // Both login and registration flows provide userId
                    // At this point, the cookie is ALREADY set in the browser
                    // We just need to fetch the profile
                    
                    if (!credentials?.userId) {
                        console.error("No userId provided to authorize");
                        return null;
                    }

                    const cookieHeader = req.headers.get?.("cookie") || req.headers.cookie;
                    console.log("Fetching profile for userId:", credentials.userId);
                    console.log("With cookies:", cookieHeader);

                    if (!cookieHeader) {
                        console.error("No cookies found in request");
                        return null;
                    }

                    // Fetch user profile using the cookie that's already set
                    const profileRes = await fetch(`${apiBase}/profile/${credentials.userId}`, {
                        method: 'GET',
                        headers: { "Cookie": cookieHeader }
                    });

                    if (profileRes.ok) {
                        const profileData = await profileRes.json();
                        console.log("Profile fetched successfully:", profileData);
                        return { ...profileData, backendCookie: cookieHeader };
                    } else {
                        console.error("Profile fetch failed:", profileRes.status);
                        return null;
                    }

                } catch (e) {
                    console.error("Auth error:", e);
                    return null;
                }
            }
        })
    ],
    callbacks: {
        async jwt({ token, user }) {
            // First login: attach backendCookie and other user details to the token
            if (user) {
                token.backendCookie = user.backendCookie
                token.user = user
            }
            return token
        },
        async session({ session, token }) {
            // Expose user and backendCookie to the session
            session.user = token.user
            session.backendCookie = token.backendCookie
            return session
        },
        async redirect({ url, baseUrl }) {
            console.log("Redirect callback:", { url, baseUrl });
            
            // If url is relative, make it absolute
            if (url.startsWith("/")) {
                return `${baseUrl}${url}`;
            }
            // If url is on same origin, allow it
            else if (new URL(url).origin === baseUrl) {
                return url;
            }
            // Default: go to feed after login
            return `${baseUrl}/feed/public`;
        }
    },
    pages: {
        signIn: '/login',
        error: '/login',
    },
    debug: process.env.NODE_ENV === 'development',
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }