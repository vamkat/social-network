"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getUserGroups({ limit, offset }) {
    try {
        const response = await serverApiRequest("/groups/user/", {
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

        return { success: true, data: response };

    } catch (error) {
        console.error("Error fetching user groups: ", error);
        return { success: false, error: error.message };
    }
}
