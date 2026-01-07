"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPendingRequests({ groupId, limit = 20, offset = 0 }) {
    try {
        const response = await serverApiRequest("/group/pending", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                limit: limit,
                offset: offset
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return response;
    } catch (error) {
        console.error("Error fetching pending requests:", error);
        return { users: [] };
    }
}
