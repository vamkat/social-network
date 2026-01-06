"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function respondToGroupInvite({ groupId, accept }) {
    try {
        const response = await serverApiRequest("/group/invite/response", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                accept: accept,
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error responding to group invite:", error);
        return { success: false, error: error.message };
    }
}
