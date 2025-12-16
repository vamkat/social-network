"use server";

import { serverApiRequest } from "@/lib/server-api";

/**
 * Updates the user's profile information.
 * Expects { username, first_name, last_name, date_of_birth, avatar_id, about }
 * Or specific fields for other updates if we merge them, but let's stick to what updateProfileInfo did.
 * Note: The previous file had multiple exports: updateProfilePrivacy, updateProfileEmail, updateProfilePassword, updateProfileInfo.
 * I should probably implement all of them or create separate files?
 * The plan said "updateProfile server action". I'll put them all here as exports or separate files.
 * Putting them in one file `update-profile.js` seems cleaner for now as they are related.
 */

export async function updateProfileInfo(data) {
    try {
        const response = await serverApiRequest("/profile/update", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile info:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfilePrivacy({ bool }) {
    try {
        const response = await serverApiRequest("/account/update/public", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                public: bool,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile privacy:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfileEmail({ email }) {
    try {
        const response = await serverApiRequest("/account/update/email", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile email:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfilePassword({ oldPassword, newPassword }) {
    try {
        const response = await serverApiRequest("/account/update/password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile password:", error);
        return { success: false, error: error.message };
    }
}
