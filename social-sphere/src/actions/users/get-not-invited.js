"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotInvited({ groupId, limit = 20, offset = 0 } = {}) {
    try {
        if (!groupId) {
            console.error("Group ID is required to fetch followers");
            return [];
        }

        const followers = await serverApiRequest("/group/notinvited", {
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

        return followers;

    } catch (error) {
        console.error("Error fetching followers:", error);
        return [];
    }
}