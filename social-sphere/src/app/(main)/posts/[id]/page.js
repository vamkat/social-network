import PostCard from "@/components/ui/post-card";
import { fetchPostById } from "@/services/posts/posts";
import { notFound } from "next/navigation";

export default async function SinglePostPage({ params }) {
    const { id } = await params;
    const post = await fetchPostById(id);

    if (!post) {
        return notFound();
    }

    return (
        <div className="w-full py-8">
            <div className="max-w-7xl mx-auto px-6">
                <div className="flex gap-6">
                    <aside className="hidden xl:block w-48 shrink-0" />

                    <main className="flex-1 max-w-2xl mx-auto min-w-0">
                        {/* <div className="feed-header text-center">
                            <h1 className="feed-title">Post</h1>
                            <p className="feed-subtitle">Full conversation and activity</p>
                        </div> */}

                        <PostCard post={post} />
                    </main>

                    <aside className="hidden lg:block w-80 shrink-0" />
                </div>
            </div>
        </div>
    );
}
