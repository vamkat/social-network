import { NextResponse } from 'next/server';

export async function POST(request) {
    try {
        const cookieHeader = request.headers.get('cookie');

        const apiBase = process.env.API_BASE || "http://localhost:8081";
        const logoutEndpoint = process.env.LOGOUT || "/logout";

        const headers = {};
        if (cookieHeader) {
            headers['Cookie'] = cookieHeader;
        }

        const backendResponse = await fetch(`${apiBase}${logoutEndpoint}`, {
            method: "POST",
            headers: headers,
        });

        const responseData = await backendResponse.json().catch(() => null);
        const setCookieHeader = backendResponse.headers.get('set-cookie');

        const response = NextResponse.json(
            responseData || { error: "Logout failed" },
            { status: backendResponse.status }
        );

        if (setCookieHeader) {
            response.headers.set('Set-Cookie', setCookieHeader);
        }

        return response;
    } catch (error) {
        console.error("Logout API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}
