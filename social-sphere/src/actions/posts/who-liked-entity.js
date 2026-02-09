"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getWhoLikedEntity(entityID) {
    try {
        const url = `/reactions/${entityID}`;
        const response = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        if (!response.ok) {
            return {success: false, status: response.status, error: response.message};
        }

        return { success: true, data: response.data };

    } catch (error) {
        return { success: false, error: error.message };
    }
}
