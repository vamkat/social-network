"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getConvByID({interlocutorId , convId}) {
    try {
        const url = `/my/chat/${convId}/preview?interlocutor_id=${interlocutorId}`;

        const response = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        // Return success wrapper
        return { success: true, data: response };
    } catch (error) {
        console.error("Error fetching groups:", error);
        return { success: false, error: error.message };
    }
}