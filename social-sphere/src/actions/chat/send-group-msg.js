"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function sendGroupMsg({ groupId, msg }) {
    try {
        const url = `/groups/${groupId}/chat`;
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                message_body: msg
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Send Group Message Error:", error);
        return { success: false, error: error.message };
    }
}
