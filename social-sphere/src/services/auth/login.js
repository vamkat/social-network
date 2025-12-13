"use client";

import { apiRequest } from "@/lib/api";

export async function login(credentials) {
    try {
        // request
        const apiResp = await apiRequest("/login", {
            method: "POST",
            body: JSON.stringify(credentials),
        });

        // check if user id is provided
        if (!apiResp.UserId) {
            return {
                success: false,
                error: "Login failed - no user ID returned"
            };
        }

        // all good
        return { success: true, user_id: apiResp.UserId };

    } catch (error) {
        console.error("Error: ", error)
        return { success: false }
    }
}