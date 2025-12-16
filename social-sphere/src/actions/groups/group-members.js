"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getGroupMembers({ group_id, limit, offset }) {
    try {
        const response = await serverApiRequest("/group/members", {
            method: "POST",
            body: JSON.stringify({
                group_id: group_id,
                limit: limit,
                offset: offset
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error fetching group members: ", error);
        return { success: false, error: error.message };
    }
}
