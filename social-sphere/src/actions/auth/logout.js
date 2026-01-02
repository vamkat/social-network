"use server";

import { serverApiRequest } from "@/lib/server-api";
import { redirect } from "next/navigation";

export async function logout() {
    try {
        await serverApiRequest("/logout", {
            method: "POST",
            forwardCookies: true // Forward the cookie clearing headers
        });
    } catch (error) {
        console.error("Logout Action Error:", error);
        return { success: false, error: error.message };
    }

    // Redirect after successful logout
    redirect("/login");
}
