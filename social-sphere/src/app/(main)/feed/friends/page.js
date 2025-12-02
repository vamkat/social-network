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
                postFormProps={{
                    defaultVisibility: "friends",
                    visibilityOptions: [
                        { value: "friends", label: "Friends", helper: "Only visible in Friends feed." },
                        { value: "public", label: "Public", helper: "Also shows in the Public feed." },
                        { value: "custom", label: "Select members", helper: "Choose specific members who can view this post." },
                    ],
                }}
            />

            <FeedList initialPosts={initialPosts} fetchPosts={fetchFeedPosts} />
        </div>
    );
}
