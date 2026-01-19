"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function sendMsg({interlocutor , msg}) {
    try {
        const url = `/my/chat/${interlocutor}`;
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify({
                message_body: msg
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Create Event Action Error:", error);
        return { success: false, error: error.message };
    }
}