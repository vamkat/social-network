import { getMockPosts } from "@/mock-data/posts";
import PostCard from "@/components/ui/post-card";

export default function PublicFeedPage() {
    const posts = getMockPosts();

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Public Feed</h1>
                <p className="feed-subtitle">What's happening around the world</p>
            </div>

            <div className="flex flex-col">
                {posts.map((post, i) => (
                    <PostCard key={i} post={post} />
                ))}
            </div>
        </div >
    );
}
