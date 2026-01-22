"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotInvited({ groupId, limit = 20, offset = 0 } = {}) {
    try {
        if (!groupId) {
            console.error("Group ID is required to fetch followers");
            return [];
        }
        const url = `/groups/${groupId}/invitable-followers?limit=${limit}&offset=${offset}`;
        const followers = await serverApiRequest(url, {
            method: "GET"
        });

        return followers;

    } catch (error) {
        console.error("Error fetching followers:", error);
        return [];
    }
}