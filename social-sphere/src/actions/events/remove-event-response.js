"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function removeEventResponse(data) {
    try {
        const apiResp = await serverApiRequest("/events/remove-response", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Remove Event Response Action Error:", error);
        return { success: false, error: error.message };
    }
}
