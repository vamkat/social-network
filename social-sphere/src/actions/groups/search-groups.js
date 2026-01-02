"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function searchGroups({ query, limit = 10, offset = 0 }) {
    try {
        const response = await serverApiRequest("/search/group", {
            method: "POST",
            body: JSON.stringify({
                query: query,
                limit: limit,
                offset: offset,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, groups: response.groups || [] };

    } catch (error) {
        console.error("Error searching groups:", error);
        return { success: false, error: error.message, groups: [] };
    }
}
