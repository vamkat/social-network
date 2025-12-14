import { apiRequest } from "@/lib/api";

export async function follow(userId) {
    try {
        await apiRequest('/users/follow', {
            method: "POST",
            body: JSON.stringify({
                target_user_id: userId,
            }),
        });

        return { success: true };

    } catch (error) {
        console.error("Error sending follow request:", error);
        return { success: false, error: error };
    }
}

export async function unfollow(userId) {
    try {
        await apiRequest('/users/unfollow', {
            method: "POST",
            body: JSON.stringify({
                target_user_id: userId
            }),
        });

        return { success: true };
    } catch (error) {
        console.error("Error unfollow request: ", error);
        return { success: false, error: error };
    }
}