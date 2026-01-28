"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotifCount() {
    try {
        const url = `/notifications-count`;
        const count = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });
        console.log("Count", count)
        return count;

    } catch (error) {
        console.error("Error fetching notif count:", error);
        return [];
    }
}
