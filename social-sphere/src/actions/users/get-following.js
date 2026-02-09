"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getFollowers({ userId, limit = 100, offset = 0 } = {}) {
    try {
        if (!userId) {
            return { success: false, error: "User ID is required" };
        }
        const url = `/users/${userId}/following?limit=${limit}&offset=${offset}`;
        const response = await serverApiRequest(url, {
            method: "GET"
        });

        if (!response.ok) {
            return {success: false, status: response.status, error: response.message};
        }

        return { success: true, data: response.data };

    } catch (error) {
        return { success: false, error: error.message };
    }
}
