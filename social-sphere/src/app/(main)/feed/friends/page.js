import { getMockPosts } from '@/mock-data/posts';
import PostCard from "@/components/ui/post-card";

export default function FriendsFeedPage() {
    const posts = getMockPosts();

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Friends Feed</h1>
                <p className="feed-subtitle">Updates from your friends</p>
            </div>

            <div className="flex flex-col">
                {posts.map((post, i) => (
                    <PostCard key={i} post={post} />
                ))}
            </div>
        </div>
    );
}