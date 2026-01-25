"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getMostPopular(groupId) {
    try {
        const url = `/groups/${groupId}/popular-post`
        const response = await serverApiRequest(url, {
            method: "GET"
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error getting group:", error);
        return { success: false, error: error.message };
    }
}
