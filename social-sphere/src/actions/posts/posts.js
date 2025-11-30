"use server";

import { getMockPosts, GetPostsByUserId } from "@/mock-data/posts";

export async function fetchPublicPosts(offset = 0, limit = 5) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts(offset, limit);
}

export async function fetchFeedPosts(offset = 0, limit = 5) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts(offset, limit);
}

import { unstable_cache } from "next/cache";

export async function fetchUserPosts(userID, offset = 0, limit = 5) {
    const getCachedPosts = unstable_cache(
        async () => {
            console.log(`[CACHE MISS] Fetching posts for user ${userID}, offset ${offset}, limit ${limit}`);
            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 100));
            return GetPostsByUserId(userID, offset, limit);
        },
        [`user-posts-${userID}-${offset}-${limit}`],
        { revalidate: 60 }
    );

    return getCachedPosts();
}