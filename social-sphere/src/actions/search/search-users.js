"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function SearchUsers({ query, limit }) {
    try {
        const response = await serverApiRequest("/users/search", {
            method: "POST",
            body: JSON.stringify({
                query: query,
                limit: limit,
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });
        return response;
    } catch (error) {
        console.error("Error searching users:", error);
        return { success: false, error: error.message };
    }
}
