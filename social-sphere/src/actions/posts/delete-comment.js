"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deleteComment(commentId) {
    try {
        const apiResp = await serverApiRequest("/comments/delete", {
            method: "POST",
            body: JSON.stringify({ entity_id: commentId }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Comment Action Error:", error);
        return { success: false, error: error.message };
    }
}
