"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getConv({first=false, beforeDate=null, limit}) {
    try {
        let url = null;
        if (first === true) {
            url = `/chat/get-private-conversations?limit=${limit}`;
        } else {
            url = `/chat/get-private-conversations?before-date=${beforeDate}&limit=${limit}`;
        }

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