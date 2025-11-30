import { fetchFeedPosts } from "@/actions/posts/posts";
import FeedActions from "@/components/ui/feed-actions";
import FeedList from "@/components/feed/feed-list";

export default async function FriendsFeedPage() {
    const initialPosts = await fetchFeedPosts(0, 5);

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Friends Feed</h1>
                <p className="feed-subtitle">Updates from your friends</p>
            </div>

            <FeedActions
                ctaProps={{
                    title: "Create new friends post",
                    subtitle: "Share what is on your mind with your friends.",
                    actionLabel: "+ Post",
                }}
            />

            <FeedList initialPosts={initialPosts} fetchPosts={fetchFeedPosts} />
        </div>
    );
}
