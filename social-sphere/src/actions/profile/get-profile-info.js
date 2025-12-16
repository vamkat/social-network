"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getProfileInfo(userId) {
    try {
        const user = await serverApiRequest(`/profile/${userId}`, {
            method: "POST", // API seems to use POST for getting profile info based on previous code
            // No body needed apparently based on previous service
            forwardCookies: true
        });

        return user;

    } catch (error) {
        console.error("Error fetching profile info:", error);
        return { success: false, error: error.message };
    }
}
