"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function validateUpload(fileId) {
    try {
        await serverApiRequest("/validate-file-upload", {
            method: "POST",
            body: JSON.stringify({ file_id: fileId }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true };

    } catch (error) {
        console.error("Validate Upload Action Error:", error);
        return { success: false, error: error.message };
    }
}
