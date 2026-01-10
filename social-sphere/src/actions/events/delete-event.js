"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deleteEvent(eventId) {
    try {
        const apiResp = await serverApiRequest("/events/delete", {
            method: "POST",
            body: JSON.stringify({ entity_id: eventId }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
