import { getMockPosts } from "@/mock-data/posts";
import PostFeed from "@/components/features/feed/post-feed";

export default function PublicFeedPage() {
    const posts = getMockPosts();

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Public Feed</h1>
                <p className="feed-subtitle">What's happening around the world</p>
            </div>

            <PostFeed posts={posts} pageSize={4} />
        </div >
    );
}
