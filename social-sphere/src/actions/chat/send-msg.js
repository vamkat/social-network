"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function sendMsg({interlocutor , msg}) {
    try {
        const apiResp = await serverApiRequest("/chat/create-pm", {
            method: "POST",
            body: JSON.stringify({
                interlocutor_id: interlocutor,
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