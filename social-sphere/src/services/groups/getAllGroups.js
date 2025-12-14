import { serverApiRequest } from "@/lib/server-api";

export async function getAllGroups(limit, offset) {
    try {
        const groups = await serverApiRequest(`/groups/paginated`, {
            method: "POST",
            body: JSON.stringify({
                limit: limit,
                offset: offset,
            }),
        });

        return groups;

    } catch (error) {
        console.error("Error fetching groups:", error);
        return { success: false, error: error };
    }
}