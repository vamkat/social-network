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
                postFormProps={{
                    defaultVisibility: "public",
                    visibilityOptions: [
                        { value: "public", label: "Public", helper: "Shown in both Public and Friends feeds." },
                        { value: "friends", label: "Friends", helper: "Only visible in Friends feed." },
                        { value: "custom", label: "Select members", helper: "Choose specific members who can view this post." },
                    ],
                }}
            />

            <FeedList initialPosts={initialPosts} fetchPosts={fetchPublicPosts} />
        </div >
    );
}
