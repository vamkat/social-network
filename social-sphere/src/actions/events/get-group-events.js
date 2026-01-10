"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getGroupEvents({ groupId, limit = 10, offset = 0 }) {
    try {
        const apiResp = await serverApiRequest("/events/", {
            method: "POST",
            body: JSON.stringify({
                entity_id: groupId,
                limit: limit,
                offset: offset
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, data: apiResp };

    } catch (error) {
        console.error("Get Group Events Action Error:", error);
        return { success: false, error: error.message };
    }
}
