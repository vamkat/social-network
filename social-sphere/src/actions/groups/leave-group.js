"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function leaveGroup({ userId, groupId }) {
    try {
        const response = await serverApiRequest("/group/leave", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
                group_id: groupId,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error leaving group:", error);
        return { success: false, error: error.message };
    }
}
