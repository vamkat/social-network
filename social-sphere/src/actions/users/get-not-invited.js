"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotInvited({ groupId, limit = 20, offset = 0 } = {}) {
    try {
        if (!groupId) {
            return { success: false, error: "Group ID is required" };
        }
        const url = `/groups/${groupId}/invitable-followers?limit=${limit}&offset=${offset}`;
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