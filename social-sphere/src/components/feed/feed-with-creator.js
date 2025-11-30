"use client";

import FeedList from "@/components/feed/feed-list";
import FeedPostCTA from "../ui/feed-post-creator";

export default function FeedWithCreator({
    initialPosts = [],
    fetchPosts,
    ctaProps = {},
    containerClassName = "",
}) {
    return (
        <div className={`space-y-4 ${containerClassName}`}>
            <FeedPostCTA {...ctaProps} />
            <FeedList initialPosts={initialPosts} fetchPosts={fetchPosts} />
        </div>
    );
}
