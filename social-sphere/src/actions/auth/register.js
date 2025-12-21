"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function register(userData) {
    try {
        // register with a public profile by default
        userData.public = true;

        const apiResp = await serverApiRequest("/register", {
            method: "POST",
            body: JSON.stringify(userData),
            headers: {
                "Content-Type": "application/json"
            },
            forwardCookies: true
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Register Action Error:", error);
        return { success: false, error: error.message };
    }
}
