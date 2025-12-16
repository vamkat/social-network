"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPublicPosts({ limit = 10, offset = 0 } = {}) {
    try {
        const posts = await serverApiRequest("/public-feed", {
            method: "POST",
            body: JSON.stringify({
                limit: limit,
                offset: offset
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return posts;

    } catch (error) {
        console.error("Error fetching public posts:", error);
        // Returning empty array or null might be better than an object with error property depending on how component handles it
        // based on page.js, it expects an array: posts.map...
        // so let's throw or return empty array if it fails, or maybe the existing component expects a specific error structure?
        // Looking at page.js: const posts = await getPublicPosts(...) -> posts.map
        // So it expects an array.

        // If serverApiRequest throws, we catch it here.
        // If we return an object like { success: false }, map will fail.
        // Let's return an empty array for now to prevent crash, or rethrow if we want the error boundary to catch it.
        // However, the original service returned { success: false, error } on error, but page.js doesn't seem to handle that check before mapping?
        // Let's re-read page.js quickly to be sure.
        // page.js: const posts = await getPublicPosts(...) -> posts.map
        // It blindly maps. So if original service returned {success:false}, page would crash.
        // The original service caught error and returned { success: false, error: error }.
        // So the current page.js is likely crashing on error anyway or I missed something.
        // Let's assume for now we want to return [] on error to be safe, or just let error bubble up to Next.js error boundary.

        return [];
    }
}
