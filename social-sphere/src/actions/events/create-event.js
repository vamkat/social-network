"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createEvent(eventData) {
    try {
        const apiResp = await serverApiRequest("/events/create", {
            method: "POST",
            body: JSON.stringify(eventData),
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
