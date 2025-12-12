"use client";

import { safeApiCall } from "@/lib/api-wrapper";

/**
 * Client-side registration function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 */
export async function registerClient(formData) {
    const result = await safeApiCall("/api/auth/register", {
        method: "POST",
        body: formData,
    });

    if (result.success) {
        return { success: true, user: result.data };
    }

    return result;
}
