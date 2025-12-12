import { fetchPublicPosts } from "@/services/posts/posts";
import PublicFeedClient from "./client";

export default async function PublicFeedPage() {
    const initialPosts = await fetchPublicPosts(0, 5);

    return <PublicFeedClient initialPosts={initialPosts} />;
}
