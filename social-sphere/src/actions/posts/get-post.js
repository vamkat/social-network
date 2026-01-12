"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPost(postId) {
    try {
        const post = await serverApiRequest("/post/", {
            method: "POST",
            body: JSON.stringify({
                entity_id: postId
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return {success: true, error: null, post: post};

    } catch (error) {
        console.error("Error fetching post:", error);
        return {success:false, error: error.message, post: null};
    }
}
