"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function editComment(commentData) {
    try {
        const apiResp = await serverApiRequest("/comments/edit", {
            method: "POST",
            body: JSON.stringify(commentData),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Edit Comment Action Error:", error);
        return { success: false, error: error.message };
    }
}
