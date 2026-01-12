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

        if (err.error.includes("permission denied")) {
            return {ok: true, permission: false};
        } else {
            throw new Error(err.error || "API error");
        }
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

    // Handle empty response bodies (like delete endpoints)
    const text = await res.text();
    if (!text || text.trim() === '') {
        return {};
    }

    try {
        return JSON.parse(text);
    } catch (e) {
        console.error('Failed to parse JSON response:', text);
        throw new Error('Invalid JSON response from server');
    }
}
