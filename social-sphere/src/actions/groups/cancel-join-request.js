"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function cancelJoinRequest({ groupId }) {
    try {
        const response = await serverApiRequest("/group/cancel-request", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error canceling join request:", error);
        return { success: false, error: error.message };
    }
}
