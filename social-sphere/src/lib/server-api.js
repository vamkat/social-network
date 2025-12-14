"use server";

import { cookies } from "next/headers";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://api-gateway:8081'

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

    return res.json();
}
