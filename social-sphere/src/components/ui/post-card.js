"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { fetchComments } from "@/actions/comments/comments";

export default function PostCard({ post }) {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true); // Initial loading for first comment
    const [loadingMore, setLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);

    const handleMouseEnter = async () => {
        if (!loading || post.NumOfComments === 0) return;

        try {
            // Fetch the last comment (offset 0, limit 1)
            const initialComments = await fetchComments(post.ID, 0, 1);
            setComments(initialComments);
            if (initialComments.length === 0) setHasMore(false);
        } catch (error) {
            console.error("Failed to load comments", error);
        } finally {
            setLoading(false);
        }
    };

    const handleLoadMore = async (e) => {
        e.preventDefault(); // Prevent link navigation if inside a link (though button should handle it)
        e.stopPropagation();

        if (loadingMore) return;
        setLoadingMore(true);

        try {
            const currentCount = comments.length;
            const limit = 2;
            // Fetch previous 2 comments
            const newComments = await fetchComments(post.ID, currentCount, limit);

            if (newComments.length < limit) {
                setHasMore(false);
            }

            if (newComments.length > 0) {
                // Prepend new comments (which are older) to the list
                // Note: fetchComments returns [newest-1, newest-2].
                // We want to display them as [newest-2, newest-1, newest].
                // So we reverse the batch before prepending.
                setComments(prev => [...newComments.reverse(), ...prev]);
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Failed to load more comments", error);
        } finally {
            setLoadingMore(false);
        }
    };

    return (
        <div className="post-card group" onMouseEnter={handleMouseEnter}>
            <Link href={`/profile/${post.BasicUserInfo.UserID}`}>
                {/* Avatar Column */}
                <div className="post-avatar-container">
                    <img src={post.BasicUserInfo.Avatar} alt="Post Avatar" className="post-avatar" />
                </div>
            </Link>

            {/* Content Column */}
            <div className="post-content-container">
                {/* Header */}
                <div className="post-header">
                    <Link href={`/profile/${post.BasicUserInfo.UserID}`}>
                        <h3 className="post-username">
                            @{post.BasicUserInfo.Username}
                        </h3>
                    </Link>
                    <span className="post-timestamp">{post.CreatedAt}</span>
                </div>

                {/* Content */}
                <Link href={`/posts/${post.ID}`}>
                    <p className="post-text">
                        {post.Content}
                    </p>
                </Link>

                {/* Post Image - Fixed Height & Cover */}
                {post.PostImage && (
                    <div>
                        <img
                            src={post.PostImage}
                            alt="Post content"
                        />
                    </div>
                )}

                {/* Footer / Actions */}
                <div className="post-actions mt-2">
                    {/* Reaction Button */}
                    <button className="action-btn action-btn-heart group/heart">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-heart">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
                        </svg>
                        <span className="text-sm font-medium">{post.NumOfHearts}</span>
                    </button>

                    {/* Comments */}
                    <button className="action-btn action-btn-comment group/comment">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-comment">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
                        </svg>
                        <span className="text-sm font-medium">{post.NumOfComments}</span>
                    </button>
                </div>

                {/* Hover Comment Preview */}
                {comments.length > 0 && (
                    <div className="hidden group-hover:block mt-4 pt-3 border-t border-(--muted)/15 animate-in fade-in slide-in-from-top-1 duration-200">

                        {hasMore && (
                            <button
                                onClick={handleLoadMore}
                                disabled={loadingMore}
                                className="w-full text-left text-xs text-(--muted) hover:text-(--foreground) mb-3 pl-11 transition-colors"
                            >
                                {loadingMore ? "Loading..." : "View previous comments"}
                            </button>
                        )}

                        <div className="flex flex-col gap-3 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                            {comments.map((comment, index) => (
                                <div key={comment.ID || index} className="flex gap-3">
                                    <div className="w-8 h-8 rounded-full overflow-hidden shrink-0 border border-(--muted)/20">
                                        <img src={comment.BasicUserInfo.Avatar} alt={comment.BasicUserInfo.Username} className="w-full h-full object-cover" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center justify-between">
                                            <span className="text-sm font-semibold">@{comment.BasicUserInfo.Username}</span>
                                            <span className="text-xs text-(--muted)">{comment.CreatedAt}</span>
                                        </div>
                                        <p className="text-sm text-(--foreground)/90 mt-1 leading-relaxed">
                                            {comment.Content}
                                        </p>
                                        <button className="flex items-center gap-1.5 mt-2 text-xs text-(--muted) hover:text-red-500 transition-colors group/comment-heart">
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-3.5 h-3.5 group-hover/comment-heart:fill-current">
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
                                            </svg>
                                            <span className="font-medium">{comment.NumOfHearts}</span>
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}