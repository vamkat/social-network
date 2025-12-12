import { NextResponse } from 'next/server';

export async function POST(request) {
    try {
        const body = await request.json();
        
        const apiBase = process.env.API_BASE || "http://api-gateway:8081";
        const loginEndpoint = process.env.LOGIN || "/login";

        // Call Golang backend
        const backendResponse = await fetch(`${apiBase}${loginEndpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body),
        });

        const responseData = await backendResponse.json().catch(() => null);
        const setCookieHeader = backendResponse.headers.get('set-cookie');

        console.log("Login response data:", responseData);
        console.log("Set-Cookie header:", setCookieHeader);

        // Create response
        const response = NextResponse.json(
            responseData || { error: "Login failed" },
            { status: backendResponse.status }
        );

        // Forward the backend cookie to the browser (same as registration)
        if (setCookieHeader) {
            const cookieParts = setCookieHeader.split(';').map(part => part.trim());
            const [nameValue] = cookieParts;
            const [name, value] = nameValue.split('=');

            const attributes = {};
            cookieParts.slice(1).forEach(part => {
                const [key, val] = part.split('=');
                attributes[key.toLowerCase()] = val || true;
            });

            const cookieOptions = {
                path: attributes.path || '/',
                httpOnly: attributes.httponly === true,
                secure: attributes.secure === true,
                domain: 'localhost',
            };

            if (attributes.expires) {
                cookieOptions.expires = new Date(attributes.expires);
            }

            response.cookies.set(name, value, cookieOptions);
            console.log("Cookie set:", name);
        }

        return response;
    } catch (error) {
        console.error("Login API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}