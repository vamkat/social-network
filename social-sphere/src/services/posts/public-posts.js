import { serverApiRequest } from "@/lib/server-api";
/*
* @ params -> limit, offset
*/
export async function getPublicPosts(params) {
    try {
        const posts = await serverApiRequest("/public-feed", {
            method: "GET",
            //body: JSON.stringify(params)
        })

        console.log("posts got:", posts);

        return posts;

    } catch (error) {
        console.error("Error fetching public posts: ", error);
        return { success:false, error: error };
    }
}