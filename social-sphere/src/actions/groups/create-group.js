"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createGroup({ group_title, group_description, group_image }) {
    try {
        const response = await serverApiRequest(`/groups/create`, {
            method: "POST",
            body: JSON.stringify({
                group_title,
                group_description,
                ...(group_image && { group_image }),
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error creating group:", error);
        return { success: false, error: error.message };
    }
}
