"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deletePost(postId) {
    try {
        const apiResp = await serverApiRequest("/posts/delete/", {
            method: "POST",
            body: JSON.stringify({ entity_id: postId }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Post Action Error:", error);
        return { success: false, error: error.message };
    }
}
