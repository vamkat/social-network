import { getPublicPosts } from "@/actions/posts/get-public-posts";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";
import Container from "@/components/layout/Container";

export const metadata = {
    title: "Public Feed",
}

export default async function PublicFeedPage() {
    // call backend for public posts
    const limit = 10;
    const offset = 0;
    const posts = await getPublicPosts({ limit, offset });

    return (
        <div className="w-full">
            {/* Create Post Section */}
            <Container className="pt-6 md:pt-10">
                <CreatePost />
            </Container>

            {/* Feed Header */}
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title px-4">Public Feed</h1>
                <p className="text-center feed-subtitle px-4">What's happening in global sphere?</p>
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
                            Be the first ever to share something on the public sphere!
                        </p>
                    </div>
                )}
            </Container>
        </div>
    );
}