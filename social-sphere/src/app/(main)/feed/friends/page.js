import { fetchFeedPosts } from "@/services/posts/posts";
import FriendsFeedClient from "./client";

export default async function FriendsFeedPage() {
    const initialPosts = await fetchFeedPosts(0, 5);

    return <FriendsFeedClient initialPosts={initialPosts} />;
}