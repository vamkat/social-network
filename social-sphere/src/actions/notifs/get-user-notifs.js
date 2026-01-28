"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotifs({ limit = 10, offset = 0 } = {}) {
    try {
        const url = `/notifications?limit=${limit}&offset=${offset}&read_only=false`;
        const notifs = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return notifs;

    } catch (error) {
        console.error("Error fetching notifications:", error);
        return [];
    }
}
