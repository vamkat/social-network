"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function SearchGroups({ query, limit, offset }) {
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
        return response;
    } catch (error) {
        console.error("Error searching groups:", error);
        return { success: false, error: error.message };
    }
}
