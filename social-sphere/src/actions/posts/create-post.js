"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createPost(postData) {
    try {
        const apiResp = await serverApiRequest("/posts/create", {
            method: "POST",
            body: JSON.stringify(postData),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Create Post Action Error:", error);
        return { success: false, error: error.message };
    }
}
