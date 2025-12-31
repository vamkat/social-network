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

        return post;

    } catch (error) {
        console.error("Error fetching post:", error);
        return null;
    }
}
