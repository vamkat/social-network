import { apiRequest } from "@/lib/api";
import { serverApiRequest } from "@/lib/server-api";

export async function getProfileInfo(userId) {
    const isServer = typeof window === 'undefined';
    
    try {
        const apiFn = isServer ? serverApiRequest : apiRequest;
        
        const user = await apiFn(`/profile/${userId}`, {
            method: "GET",
        });

        return user;

    } catch (error) {
        console.error("Error fetching profile:", error);
        return { success: false, error: error };
    }
}