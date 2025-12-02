import { fetchPublicPosts } from "@/actions/posts/posts";
import FeedList from "@/components/feed/feed-list";

export default async function PublicFeedPage() {
    const initialPosts = await fetchPublicPosts(0, 5);

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Public Feed</h1>
                <p className="feed-subtitle">What's happening around the world</p>
            </div>

            <FeedList initialPosts={initialPosts} fetchPosts={fetchPublicPosts} />
        </div >
    );
}
