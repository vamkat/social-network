import { LogoutButton } from "@/components/LogoutButton";
import { getPublicPosts } from "@/services/posts/public-posts";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";

export const metadata = {
    title: "Public Feed",
}

export default async function PublicFeedPage() {
    // call backend for public posts
    const limit = 10;
    const offset = 0;
    const posts = await getPublicPosts({ limit, offset });
    console.log("posts", posts);

    return (
        <div>
            <div className="pt-15 flex flex-col px-70">
                <CreatePost />
            </div>
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title">Public Feed</h1>
                <p className="text-center feed-subtitle">What's happening around the world</p>
            </div>
            <div className="section-divider mb-6" />
            <div className="pt-6 flex flex-col px-70">
                {posts.map((post, index) => {
                    if (posts.length === index + 1) {
                        return (
                            <div key={`${post.ID}-${index}`}>
                                <PostCard post={post} />
                            </div>
                        );
                    } else {
                        return <PostCard key={`${post.ID}-${index}`} post={post} />;
                    }
                })}
            </div>
            <LogoutButton />
        </div>
    );
}