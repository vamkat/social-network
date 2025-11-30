import { fetchPublicPosts } from "@/actions/posts/posts";
import FeedActions from "@/components/ui/feed-actions";
import FeedList from "@/components/feed/feed-list";

export default async function PublicFeedPage() {
    const initialPosts = await fetchPublicPosts(0, 5);

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Public Feed</h1>
                <p className="feed-subtitle">What's happening around the world</p>
            </div>

            <FeedActions
                ctaProps={{
                    title: "New public post",
                    subtitle: "Share what is happening around you with the world.",
                    actionLabel: "+ Post",
                }}
            />

            <FeedList initialPosts={initialPosts} fetchPosts={fetchPublicPosts} />
        </div >
    );
}
