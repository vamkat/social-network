"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function handleFollowRequest({ requesterId, accept }) {
    try {
        const response = await serverApiRequest("/follow/response", {
            method: "POST",
            body: JSON.stringify({
                requester_id: requesterId,
                accept: accept,
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error handling follow request:", error);
        return { success: false, error: error.message };
    }
}
