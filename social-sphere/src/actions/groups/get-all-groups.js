"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getAllGroups({ limit, offset }) {
    try {
        const response = await serverApiRequest(`/groups/paginated`, {
            method: "POST",
            body: JSON.stringify({
                limit: limit,
                offset: offset,
            }),
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
