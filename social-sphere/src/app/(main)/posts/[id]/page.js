import PostCard from "@/components/ui/post-card";
import { fetchPostById } from "@/actions/posts/posts";
import { notFound } from "next/navigation";

export default async function SinglePostPage({ params }) {
    const { id } = await params;
    const post = await fetchPostById(id);

    if (!post) {
        return notFound();
    }

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Single Post Page</h1>
                <p className="feed-subtitle">View post conversation and activity</p>
            </div>

            <PostCard post={post} />
        </div>
    );
}
