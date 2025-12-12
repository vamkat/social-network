"use client";

import { signOut } from "next-auth/react";
import { safeApiCall } from "@/lib/api-wrapper";

/**
 * Client-side logout function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 */
export async function logoutClient() {
    try {
        // Call API route directly from client
        const response = await safeApiCall("/api/auth/logout", {
            method: "POST",
        });

        if (!response.success || response.error) {
            console.error("Logout failed");
            return { success: false };
        }

        // Clear client-side sessionLL
        signOut({ callbackUrl: "/" });

        return { success: true };
    } catch (error) {
        console.error("Logout error:", error);
        return { success: false };
    }
}
