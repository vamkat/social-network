"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function toggleReaction(postId) {
    try {
        const apiResp = await serverApiRequest("/reactions/", {
            method: "POST",
            body: JSON.stringify({ entity_id: postId }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Toggle Reaction Action Error:", error);
        return { success: false, error: error.message };
    }
}
