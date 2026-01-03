"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function updateGroup(data) {
    try {
        const response = await serverApiRequest("/groups/update", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
            forwardCookies: true
        });

        return {
            success: true,
            GroupId: response.GroupId,
            FileId: response.FileId,
            UploadUrl: response.UploadUrl
        };
    } catch (error) {
        console.error("Error updating group:", error);
        return { success: false, error: error.message };
    }
}
