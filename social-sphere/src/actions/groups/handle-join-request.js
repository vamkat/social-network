"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function handleJoinRequest({ groupId, requesterId, accepted }) {
    try {
        const response = await serverApiRequest("/group/handle-request", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                requester_id: requesterId,
                accepted: accepted
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error handling join request:", error);
        return { success: false, error: error.message };
    }
}
