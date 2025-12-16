"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function register(formData) {
    try {
        // register with a public profile by default
        formData.append('public', 'true');

        // We need to act as a proxy for FormData.
        // serverApiRequest expects a body. fetch can handle FormData.

        const apiResp = await serverApiRequest("/register", {
            method: "POST",
            body: formData,
            forwardCookies: true
        });

        // The specific error handling logic from service
        if (!apiResp.UserId) {
            return {
                success: false,
                error: "Registration failed - no user ID returned"
            };
        }

        return { success: true, user_id: apiResp.UserId };

    } catch (error) {
        console.error("Register Action Error:", error);
        return { success: false, error: error.message };
    }
}
