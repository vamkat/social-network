import { getMockPosts } from "@/mock-data/posts";
import PostFeed from "@/components/features/feed/post-feed";

export default function FriendsFeedPage() {
    const posts = getMockPosts();

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Friends Feed</h1>
                <p className="feed-subtitle">Updates from your friends</p>
            </div>

            <PostFeed posts={posts} pageSize={4} />
        </div>
    );
}
