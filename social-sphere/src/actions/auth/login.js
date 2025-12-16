"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function login(credentials) {
    try {
        const apiResp = await serverApiRequest("/login", {
            method: "POST",
            body: JSON.stringify(credentials),
            forwardCookies: true, // IMPORTANT: Forward Set-Cookie from backend to client
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!apiResp.UserId) {
            return {
                success: false,
                error: "Login failed - no user ID returned"
            };
        }

        return { success: true, user_id: apiResp.UserId };

    } catch (error) {
        console.error("Login Action Error:", error);
        return { success: false, error: error.message };
    }
}
