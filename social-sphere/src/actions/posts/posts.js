"use server";

import { getMockPosts, GetPostsByUserId } from "@/mock-data/posts";

export async function fetchPublicPosts() {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts();
}

export async function fetchFeedPosts() {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts();
}

export async function fetchUserPosts(userID) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return GetPostsByUserId(userID);
}