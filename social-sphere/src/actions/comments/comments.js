"use server";

import { getCommentsForPost } from "@/mock-data/comments";

export async function fetchComments(postID, offset = 0, limit = 2) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 300));

    return getCommentsForPost(postID, offset, limit);
}
