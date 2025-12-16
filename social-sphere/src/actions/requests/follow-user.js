"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function followUser(userId) {
    try {
        const response = await serverApiRequest("/user/follow", {
            method: "POST",
            body: JSON.stringify({
                target_user_id: userId,
            }),
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
