import { getFriendsPosts } from "@/actions/posts/get-friends-posts";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";
import Container from "@/components/layout/Container";

export const metadata = {
    title: "Friends Feed",
}

export default async function FriendsFeedPage() {
    const posts = await getFriendsPosts({ limit: 10, offset: 0 });

    return (
        <div className="w-full">
            {/* Create Post Section */}
            <Container className="pt-6 md:pt-10">
                <CreatePost />
            </Container>

            {/* Feed Header */}
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title px-4">Friends Feed</h1>
                <p className="text-center feed-subtitle px-4">What's happening in your sphere?</p>
            </div>

            <div className="section-divider mb-6" />

            {/* Posts Feed */}
            <Container className="pt-6 pb-12">
                {posts?.length > 0 ? (
                    <div className="flex flex-col">
                        {posts.map((post, index) => (
                            <PostCard key={`${post.post_id}-${index}`} post={post} />
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                        <p className="text-muted text-center max-w-md px-4">
                            Your friends haven't shared anything yet. <br></br> Why not be the first?
                        </p>
                    </div>
                )}
            </Container>
        </div>
    );
}