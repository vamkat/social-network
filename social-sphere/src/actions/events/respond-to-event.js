"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function respondToEvent(data) {
    try {
        const apiResp = await serverApiRequest("/events/respond", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Respond to Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
