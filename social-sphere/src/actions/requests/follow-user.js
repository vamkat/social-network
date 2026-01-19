"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function followUser(userId) {
    try {
        const url = `/users/${userId}/follow`;
        const response = await serverApiRequest(url, {
            method: "POST",
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error following user:", error);
        return { success: false, error: error.message };
    }
}
