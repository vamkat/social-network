"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function markAsRead({convID , lastMsgID}) {
    try {
        const apiResp = await serverApiRequest("/chat/update-last-read-pm", {
            method: "POST",
            body: JSON.stringify({
                conversation_id: convID,
                last_read_message_id: lastMsgID
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Mark as read error: ", error);
        return { success: false, error: error.message };
    }
}