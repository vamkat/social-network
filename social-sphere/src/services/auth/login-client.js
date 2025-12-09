"use client";

import { safeApiCall } from "@/lib/api-wrapper";

/**
 * Client-side login function that calls the API route directly.
 * This ensures the backend cookie is properly set in the browser.
 */
export async function loginClient(credentials) {
    const result = await safeApiCall("/api/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(credentials),
    });

    if (result.success) {
        return { success: true, user: result.data };
    }

    return { success: false, error: result.error || "Login failed" };
}