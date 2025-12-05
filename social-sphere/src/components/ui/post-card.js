"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { fetchComments } from "@/actions/comments/comments";
import { getUserByID } from "@/mock-data/users";

export default function PostCard({ post }) {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true); // Initial loading for first comment
    const [loadingMore, setLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [showComposer, setShowComposer] = useState(false);
    const [draftComment, setDraftComment] = useState("");
    const [postContent, setPostContent] = useState(post.Content ?? "");
    const [isEditingPost, setIsEditingPost] = useState(false);
    const [postDraft, setPostDraft] = useState(post.Content ?? "");
    const [editingCommentId, setEditingCommentId] = useState(null);
    const [editingText, setEditingText] = useState("");
    const router = useRouter();
    const currentUser = getUserByID("1"); // Mock current user (matches Navbar)
    const isOwnPost = Boolean(
        currentUser &&
        post?.BasicUserInfo?.UserID &&
        String(post.BasicUserInfo.UserID) === String(currentUser.ID)
    );

    useEffect(() => {
        const lastComment = post.LastComment;
        if (lastComment) {
            setComments([lastComment]);
            const total = Number(post.NumOfComments ?? 0);
            setHasMore(total > 1);
        } else {
            setComments([]);
            setHasMore(false);
        }
        setLoading(false);
    }, [post.ID, post.LastComment, post.NumOfComments]);

    const handleToggleComposer = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setShowComposer((v) => !v);
    };

    const handleDraftChange = (e) => {
        setDraftComment(e.target.value);
    };

    const handleCancelComposer = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setDraftComment("");
        setShowComposer(false);
    };

    const handleStartEditPost = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setPostDraft(postContent);
        setIsEditingPost(true);
    };

    const handleCancelEditPost = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setPostDraft(postContent);
        setIsEditingPost(false);
    };

    const handleSaveEditPost = (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (!postDraft.trim()) return;
        setPostContent(postDraft);
        setIsEditingPost(false);
    };

    const handleStartEdit = (comment) => {
        setEditingCommentId(comment.CommentId);
        setEditingText(comment.Body);
    };

    const handleCancelEdit = () => {
        setEditingCommentId(null);
        setEditingText("");
    };

    const handleSaveEdit = (commentId) => {
        if (!editingText.trim()) {
            handleCancelEdit();
            return;
        }
        setComments((prev) =>
            prev.map((c) =>
                c.CommentId === commentId ? { ...c, Body: editingText } : c
            )
        );
        handleCancelEdit();
    };

    const handleLoadMore = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (loadingMore) return;
        setLoadingMore(true);

        try {
            const currentCount = comments.length;
            const limit = 2;
            const newComments = await fetchComments(post.ID, currentCount, limit);
            const existingIds = new Set(comments.map((c) => c.CommentId));
            const uniqueNew = newComments.filter((c) => !existingIds.has(c.CommentId));

            if (uniqueNew.length < limit) {
                setHasMore(false);
            }

            if (uniqueNew.length > 0) {
                setComments((prev) => [...uniqueNew.reverse(), ...prev]);
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Failed to load more comments", error);
        } finally {
            setLoadingMore(false);
        }
    };

    const handleHeartClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        console.log("Heart clicked");
        // TODO: wire up heart reaction
    };

    const handleOpenPost = (e) => {
        const interactive = e.target.closest("button, a, textarea, input, select, option");
        if (interactive) return;
        e.preventDefault();
        router.push(`/posts/${post.ID}`);
    };

    return (
        <div className="post-card group" onClick={handleOpenPost}>
            <Link href={`/profile/${post.BasicUserInfo.UserID}`}>
                <div className="post-avatar-container">
                    <img src={post.BasicUserInfo.Avatar} alt="Post Avatar" className="post-avatar" />
                </div>
            </Link>

            <div className="post-content-container">
                <div className="post-header">
                    <Link href={`/profile/${post.BasicUserInfo.UserID}`}>
                        <h3 className="post-username">
                            @{post.BasicUserInfo.Username}
                        </h3>
                    </Link>
                    <div className="flex items-center justify-between gap-3 text-xs text-(--muted)">
                        <div className="flex items-center gap-2">
                            <span className="post-timestamp">{post.CreatedAt}</span>
                            <span>•</span>
                            <span className="capitalize">{post.Visibility ?? post.visibility ?? "public"}</span>
                        </div>
                        {isOwnPost && (
                            <div className="flex items-center gap-2">
                                <button
                                    className="text-(--muted) hover:text-(--foreground) transition-colors"
                                    onClick={handleStartEditPost}
                                >
                                    Edit
                                </button>
                                <span className="post-delete-btn">Delete</span>
                            </div>
                        )}
                    </div>
                </div>

                {isEditingPost ? (
                    <div className="space-y-2">
                        <textarea
                            className="w-full rounded-md border border-(--muted)/30 px-3 py-2 text-sm bg-transparent focus:outline-none focus:ring-2 focus:ring-(--foreground)/20"
                            rows={4}
                            value={postDraft}
                            onChange={(e) => setPostDraft(e.target.value)}
                        />
                        <div className="flex items-center gap-2 text-xs">
                            <button
                                type="button"
                                className="px-3 py-1 rounded-md bg-(--foreground) text-(--background) disabled:opacity-60"
                                disabled={!postDraft.trim()}
                                onClick={handleSaveEditPost}
                            >
                                Save
                            </button>
                            <button
                                type="button"
                                className="px-3 py-1 rounded-md border border-(--muted)/30 text-(--muted) hover:text-(--foreground)"
                                onClick={handleCancelEditPost}
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                ) : (
                    <Link href={`/posts/${post.ID}`}>
                        <p className="post-text">
                            {postContent}
                        </p>
                    </Link>
                )}

                {post.PostImage && (
                    <div>
                        <img
                            src={post.PostImage}
                            alt="Post content"
                        />
                    </div>
                )}

                <div className="post-footer flex flex-col gap-3">
                    <div className="flex items-center justify-between">
                        <div className="post-actions mt-2">
                            <button className="action-btn action-btn-heart group/heart" onClick={handleHeartClick}>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-heart">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
                                </svg>
                                <span className="text-sm font-medium">{post.NumOfHearts}</span>
                            </button>

                            <button className="action-btn action-btn-comment group/comment">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-comment">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
                                </svg>
                                <span className="text-sm font-medium">{post.NumOfComments}</span>
                            </button>
                        </div>

                        <button
                            onClick={handleToggleComposer}
                            className="action-btn action-btn-comment group/comment"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-4 h-4">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
                            </svg>
                            <span className="text-sm font-small">{showComposer ? "Hide" : "Add Comment"}</span>
                        </button>
                    </div>

                    {showComposer && (
                        <div className="w-full border border-(--muted)/20 rounded-lg p-3 space-y-2">
                            <textarea
                                value={draftComment}
                                onChange={handleDraftChange}
                                rows={3}
                                className="w-full rounded-md border border-(--muted)/30 px-3 py-2 text-sm bg-transparent focus:outline-none focus:ring-2 focus:ring-(--foreground)/20"
                                placeholder="Write a comment..."
                            />
                            <div className="flex items-center justify-end gap-2 text-xs">
                                <button
                                    type="button"
                                    onClick={handleCancelComposer}
                                    className="px-3 py-1 rounded-md border border-(--muted)/30 text-(--muted) hover:text-(--foreground)"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="button"
                                    disabled={!draftComment.trim()}
                                    onClick={(e) => {
                                        e.preventDefault();
                                        e.stopPropagation();
                                        // TODO: wire up createComment(post.ID, draftComment)
                                    }}
                                    className="px-3 py-1 rounded-md bg-(--foreground) text-(--background) text-xs disabled:opacity-60"
                                >
                                    Post
                                </button>
                            </div>
                        </div>
                    )}
                </div>

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
                            {comments.map((comment, index) => {
                                const isOwner = currentUser && String(comment.Creator?.UserID) === String(currentUser.ID);
                                const isEditing = editingCommentId === comment.CommentId;
                                return (
                                    <div key={comment.CommentId || index} className="flex gap-3">
                                        <Link href={`/profile/${comment.Creator.UserID}`} className="w-8 h-8 rounded-full overflow-hidden shrink-0 border border-(--muted)/20">
                                            <img src={comment.Creator.Avatar} alt={comment.Creator.Username} className="w-full h-full object-cover" />
                                        </Link>
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center justify-between">
                                                <Link href={`/profile/${comment.Creator.UserID}`} className="text-sm font-semibold hover:underline">
                                                    @{comment.Creator.Username}
                                                </Link>
                                                <div className="flex items-center gap-2 text-xs text-(--muted)">
                                                    <span>{comment.CreatedAt}</span>
                                                    {isOwner && !isEditing && (
                                                        <button
                                                            className="text-(--muted) hover:text-(--foreground) transition-colors"
                                                            onClick={(e) => {
                                                                e.preventDefault();
                                                                e.stopPropagation();
                                                                handleStartEdit(comment);
                                                            }}
                                                        >
                                                            Edit
                                                        </button>
                                                    )}
                                                </div>
                                            </div>

                                            {isEditing ? (
                                                <div className="mt-2 space-y-2">
                                                    <textarea
                                                        className="w-full rounded-md border border-(--muted)/30 px-3 py-2 text-sm bg-transparent focus:outline-none focus:ring-2 focus:ring-(--foreground)/20"
                                                        rows={3}
                                                        value={editingText}
                                                        onChange={(e) => setEditingText(e.target.value)}
                                                    />
                                                    <div className="flex items-center gap-2 text-xs">
                                                    <button
                                                        type="button"
                                                        className="px-3 py-1 rounded-md bg-(--foreground) text-(--background) disabled:opacity-60"
                                                        disabled={!editingText.trim()}
                                                        onClick={(e) => {
                                                            e.preventDefault();
                                                            e.stopPropagation();
                                                            handleSaveEdit(comment.CommentId);
                                                        }}
                                                    >
                                                        Save
                                                    </button>
                                                        <button
                                                            type="button"
                                                            className="px-3 py-1 rounded-md border border-(--muted)/30 text-(--muted) hover:text-(--foreground)"
                                                            onClick={(e) => {
                                                                e.preventDefault();
                                                                e.stopPropagation();
                                                                handleCancelEdit();
                                                            }}
                                                        >
                                                            Cancel
                                                        </button>
                                                    </div>
                                                </div>
                                            ) : (
                                                <p className="text-sm text-(--foreground)/90 mt-1 leading-relaxed">
                                                    {comment.Body}
                                                </p>
                                            )}

                                            {!isEditing && (
                                                <button className="flex items-center gap-1.5 mt-2 text-xs text-(--muted) hover:text-red-500 transition-colors group/comment-heart">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-3.5 h-3.5 group-hover/comment-heart:fill-current">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
                                                    </svg>
                                                    <span className="font-medium">{comment.ReactionsCount}</span>
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
