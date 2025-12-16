"use server";

import { cookies } from "next/headers";

const API_BASE = process.env.SERVER_API_BASE || 'http://localhost:8081'

export async function serverApiRequest(endpoint, options = {}) {
    const cookieStore = await cookies();
    const jwt = cookieStore.get("jwt")?.value;

    const headers = { ...(options.headers || {}) };
    if (jwt) headers["Cookie"] = `jwt=${jwt}`;

    const res = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
        cache: "no-store"
    });

    if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || "API error");
    }

    if (options.forwardCookies) {
        // Handle multiple Set-Cookie headers
        const setCookieHeaders = res.headers.getSetCookie ? res.headers.getSetCookie() : [];

        // Fallback for environments where getSetCookie might not be available
        if (setCookieHeaders.length === 0) {
            const header = res.headers.get('Set-Cookie');
            if (header) setCookieHeaders.push(header);
        }

        if (setCookieHeaders.length > 0) {
            setCookieHeaders.forEach(cookieStr => {
                // Simple parsing to extract name, value and path/httpOnly etc would be complex
                // But cookies().set(name, value, options) requires parsed data.
                // Alternative: Let's try to parse at least the name and value.

                const parts = cookieStr.split(';');
                const [nameValue, ...optionsParts] = parts;
                const [name, ...valueParts] = nameValue.split('=');
                const value = valueParts.join('=');

                if (name && value !== undefined) {
                    const cookieOptions = {
                        secure: true,
                        httpOnly: true,
                        path: '/',
                        sameSite: 'lax', // Default safe value
                    };

                    optionsParts.forEach(part => {
                        const [optKey, optVal] = part.trim().split('=');
                        const keyLower = optKey.toLowerCase();
                        if (keyLower === 'path') cookieOptions.path = optVal;
                        if (keyLower === 'httponly') cookieOptions.httpOnly = true;
                        if (keyLower === 'secure') cookieOptions.secure = true;
                        if (keyLower === 'samesite') cookieOptions.sameSite = optVal.toLowerCase();
                        if (keyLower === 'max-age') cookieOptions.maxAge = parseInt(optVal);
                        if (keyLower === 'expires') cookieOptions.expires = new Date(optVal);
                    });

                    cookieStore.set(name.trim(), value, cookieOptions);
                }
            });
        }
    }

    return res.json();
}
