"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPendingRequestsCount({ groupId }) {
    try {
        const response = await serverApiRequest("/group/pending-count", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error fetching user groups: ", error);
        return { success: false, error: error.message };
    }
}