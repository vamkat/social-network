"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function logout() {
    try {
        const apiResp = await serverApiRequest("/logout", {
            method: "POST",
            forwardCookies: true // Forward the cookie clearing headers
        });

        return { success: true };

    } catch (error) {
        console.error("Logout Action Error:", error);
        return { success: false, error: error.message };
    }
}
