"use server";

import { getMockUser } from "@/data/mock-users";

export async function fetchUserProfile(username) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockUser(username);
}

export async function toggleFollowUser(username) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // In a real app, this would call the backend API
    // For now, we just return success
    return { success: true };
}

export async function togglePrivacy(username) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // In a real app, this would call the backend API
    return { success: true };
}
