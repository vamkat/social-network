"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getFollowers({ userId, limit = 100, offset = 0 } = {}) {
    try {
        if (!userId) {
            console.error("User ID is required to fetch followers");
            return [];
        }

        const followers = await serverApiRequest("/users/followers/paginated", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
                limit: limit,
                offset: offset
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return followers;

    } catch (error) {
        console.error("Error fetching followers:", error);
        return [];
    }
}
