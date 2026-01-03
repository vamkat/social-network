"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function requestJoinGroup({ groupId }) {
    try {
        const response = await serverApiRequest("/group/request-or-cancel", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error requesting to join group:", error);
        return { success: false, error: error.message };
    }
}
