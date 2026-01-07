"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function removeFromGroup({ groupId, memberId }) {
    try {
        const response = await serverApiRequest("/group/remove", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                member_id: memberId,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error removing member from group:", error);
        return { success: false, error: error.message };
    }
}
