"use client";

import { safeApiCall } from "@/lib/api-wrapper";

export async function getUserData(id) {

    const url = `/api/auth/profile/${id}`;

    // revalidate cache for 2-3 min in api/ endpoint
    const user = await safeApiCall(url, {
        method: "GET",
    })

    if (user.success) {
        return { success: true, userData: user.data };
    }

    return user;
}

export async function toggleFollowUser(id) {
    return { success: true };
}

export async function togglePrivacy(id) {
    return { success: true };
}
