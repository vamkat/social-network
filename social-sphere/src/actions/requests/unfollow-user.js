"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function unfollowUser(userId) {
    try {
        const response = await serverApiRequest("/user/unfollow", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error unfollowing user:", error);
        return { success: false, error: error.message };
    }
}
