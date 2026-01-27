"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getNotifs({ limit = 10, offset = 0 } = {}) {
    try {
        const url = `/notifications?limit=${limit}&offset=${offset}&read_only=false`;
        const posts = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return posts;

    } catch (error) {
        console.error("Error fetching friends posts:", error);
        return [];
    }
}
