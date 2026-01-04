import { getPost } from "@/actions/posts/get-post";
import SinglePostCard from "@/components/ui/SinglePostCard";
import { notFound } from "next/navigation";

export default async function PostPage({ params }) {
    const { id } = await params;
    const post = await getPost(id);

    console.log("Post by id: ", post);

    if (!post) {
        notFound();
    }

    return (
        <div className="min-h-screen">
            <div className="max-w-2xl mx-auto px-4 py-8">
                <SinglePostCard post={post} />
            </div>
        </div>
    );
}
