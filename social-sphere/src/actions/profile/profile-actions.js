"use server";

import { getUserByID } from "@/mock-data/users";

export async function fetchUserProfile(id) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getUserByID(id);
}

export async function toggleFollowUser(id) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // In a real app, this would call the backend API
    // For now, we just return success
    return { success: true };
}

export async function togglePrivacy(id) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    // In a real app, this would call the backend API
    return { success: true };
}
