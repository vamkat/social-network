"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function inviteToGroup({ groupId, invitedIds }) {
    try {
        const response = await serverApiRequest("/group/invite", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                invited_ids: invitedIds,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error inviting to group:", error);
        return { success: false, error: error.message };
    }
}
