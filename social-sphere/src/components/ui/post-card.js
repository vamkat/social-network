"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import { fetchComments } from "@/services/comments/comments";
import { getUserByID } from "@/mock-data/users";
import { Heart, MessageCircle, Pencil, Trash2, MoreHorizontal, Share2, Globe, Lock, Users } from "lucide-react";
import PostImage from "./post-image";

export default function PostCard({ post }) {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [isExpanded, setIsExpanded] = useState(false);
    const [draftComment, setDraftComment] = useState("");
    const [postContent, setPostContent] = useState(post.Content ?? "");
    const [isEditingPost, setIsEditingPost] = useState(false);
    const [postDraft, setPostDraft] = useState(post.Content ?? "");
    const [editingCommentId, setEditingCommentId] = useState(null);
    const [editingText, setEditingText] = useState("");
    const composerRef = useRef(null);
    const cardRef = useRef(null);
    const router = useRouter();
    const currentUser = getUserByID("1"); // Mock current user

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

    useEffect(() => {
        if (isExpanded && composerRef.current) {
            composerRef.current.focus();
        }
    }, [isExpanded]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (cardRef.current && !cardRef.current.contains(event.target)) {
                setIsExpanded(false);
            }
        };

        if (isExpanded) {
            document.addEventListener("mousedown", handleClickOutside);
        }

        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [isExpanded]);

    const handleToggleExpanded = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsExpanded((v) => !v);
    };

    const handleDraftChange = (e) => {
        setDraftComment(e.target.value);
    };

    const handleCancelComposer = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setDraftComment("");
        setIsExpanded(false);
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
        <div
            ref={cardRef}
            className="group bg-(--background) border border-(--border) rounded-2xl overflow-hidden transition-all hover:border-(--muted)/40 hover:shadow-sm mb-6"
            onClick={handleOpenPost}
        >
            {/* Header */}
            <div className="p-5 flex items-start justify-between">
                <div className="flex items-center gap-3">
                    <Link href={`/profile/${post.BasicUserInfo.UserID}`} className="shrink-0">
                        <div className="w-10 h-10 rounded-full overflow-hidden border border-(--border)">
                            <img
                                src={post.BasicUserInfo.Avatar}
                                alt="Avatar"
                                className="w-full h-full object-cover"
                            />
                        </div>
                    </Link>
                    <div>
                        <Link href={`/profile/${post.BasicUserInfo.UserID}`}>
                            <h3 className="font-semibold text-(--foreground) hover:underline decoration-2 underline-offset-2">
                                @{post.BasicUserInfo.Username}
                            </h3>
                        </Link>
                        <div className="flex items-center gap-2 text-xs text-(--muted) mt-0.5">
                            <span>{post.CreatedAt}</span>
                        </div>
                    </div>
                </div>

                {isOwnPost && (
                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                            onClick={handleStartEditPost}
                            className="p-2 text-(--muted) hover:text-(--accent) hover:bg-(--accent)/5 rounded-full transition-colors"
                            title="Edit Post"
                        >
                            <Pencil className="w-4 h-4" />
                        </button>
                        <button
                            className="p-2 text-(--muted) hover:text-red-500 hover:bg-red-500/5 rounded-full transition-colors"
                            title="Delete Post"
                        >
                            <Trash2 className="w-4 h-4" />
                        </button>
                    </div>
                )}
            </div>

            {/* Content */}
            <div className="px-5 pb-3">
                {isEditingPost ? (
                    <div className="space-y-3 mb-4">
                        <textarea
                            className="w-full rounded-xl border border-(--muted)/30 px-4 py-3 text-sm bg-(--muted)/5 focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all resize-none"
                            rows={4}
                            value={postDraft}
                            onChange={(e) => setPostDraft(e.target.value)}
                        />
                        <div className="flex items-center justify-end gap-2">
                            <button
                                type="button"
                                className="px-3 py-1.5 text-xs font-medium text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10 rounded-full transition-colors"
                                onClick={handleCancelEditPost}
                            >
                                Cancel
                            </button>
                            <button
                                type="button"
                                className="px-4 py-1.5 text-xs font-medium bg-(--accent) text-white hover:bg-(--accent-hover) rounded-full transition-colors disabled:opacity-50"
                                disabled={!postDraft.trim()}
                                onClick={handleSaveEditPost}
                            >
                                Save Changes
                            </button>
                        </div>
                    </div>
                ) : (
                    <Link href={`/posts/${post.ID}`}>
                        <p className="text-[15px] leading-relaxed text-(--foreground)/90 whitespace-pre-wrap">
                            {postContent}
                        </p>
                    </Link>
                )}
            </div>

            {/* Image */}
            {post.PostImage && (
                <PostImage src={post.PostImage} />
            )}

            {/* Actions Footer */}
            <div className="px-5 py-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-6">
                        <button
                            className="flex items-center gap-2 text-(--muted) hover:text-red-500 transition-colors group/heart"
                            onClick={handleHeartClick}
                        >
                            <Heart className={`w-5 h-5 transition-transform group-hover/heart:scale-110 ${post.IsHearted ? "fill-red-500 text-red-500" : ""}`} />
                            <span className="text-sm font-medium">{post.NumOfHearts}</span>
                        </button>

                        <button
                            className={`flex items-center gap-2 transition-colors group/comment ${isExpanded ? "text-(--accent)" : "text-(--muted) hover:text-(--accent)"}`}
                            onClick={handleToggleExpanded}
                        >
                            <MessageCircle className={`w-5 h-5 transition-transform group-hover/comment:scale-110 ${isExpanded ? "fill-(--accent)/10" : ""}`} />
                            <span className="text-sm font-medium">{post.NumOfComments}</span>
                        </button>

                        <button className="text-(--muted) hover:text-(--foreground) transition-colors">
                            <Share2 className="w-5 h-5" />
                        </button>
                    </div>
                </div>
            </div>

            {/* Expanded Section: Comments + Composer */}
            {isExpanded && (
                <div className="animate-in fade-in slide-in-from-top-2 duration-200">
                    {/* Comments List */}
                    {comments.length > 0 && (
                        <div className="bg-(--muted)/5 border-t border-(--border) px-5 py-4">
                            {hasMore && (
                                <button
                                    onClick={handleLoadMore}
                                    disabled={loadingMore}
                                    className="w-full text-left text-xs font-medium text-(--accent) hover:underline mb-4 pl-11 transition-colors"
                                >
                                    {loadingMore ? "Loading..." : "View previous comments"}
                                </button>
                            )}

                            <div className="flex flex-col gap-4 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                                {comments.map((comment, index) => {
                                    const isOwner = currentUser && String(comment.Creator?.UserID) === String(currentUser.ID);
                                    const isEditing = editingCommentId === comment.CommentId;
                                    return (
                                        <div key={comment.CommentId || index} className="flex gap-3 group/comment-item">
                                            <Link href={`/profile/${comment.Creator.UserID}`} className="shrink-0">
                                                <div className="w-8 h-8 rounded-full overflow-hidden border border-(--border)">
                                                    <img src={comment.Creator.Avatar} alt={comment.Creator.Username} className="w-full h-full object-cover" />
                                                </div>
                                            </Link>
                                            <div className="flex-1 min-w-0">
                                                <div className="bg-white dark:bg-black/20 rounded-2xl rounded-tl-none px-4 py-2 border border-(--border)">
                                                    <div className="flex items-center justify-between mb-1">
                                                        <Link href={`/profile/${comment.Creator.UserID}`}>
                                                            <span className="text-xs font-bold hover:underline">@{comment.Creator.Username}</span>
                                                        </Link>
                                                        <span className="text-[10px] text-(--muted)">{comment.CreatedAt}</span>
                                                    </div>

                                                    {isEditing ? (
                                                        <div className="space-y-2 mt-2">
                                                            <textarea
                                                                className="w-full rounded-lg border border-(--muted)/30 px-3 py-2 text-sm bg-transparent focus:outline-none focus:border-(--accent)"
                                                                rows={2}
                                                                value={editingText}
                                                                onChange={(e) => setEditingText(e.target.value)}
                                                            />
                                                            <div className="flex justify-end gap-2">
                                                                <button
                                                                    type="button"
                                                                    className="text-xs text-(--muted) hover:text-(--foreground)"
                                                                    onClick={handleCancelEdit}
                                                                >
                                                                    Cancel
                                                                </button>
                                                                <button
                                                                    type="button"
                                                                    className="text-xs font-medium text-(--accent)"
                                                                    onClick={() => handleSaveEdit(comment.CommentId)}
                                                                >
                                                                    Save
                                                                </button>
                                                            </div>
                                                        </div>
                                                    ) : (
                                                        <p className="text-sm text-(--foreground)/90 leading-relaxed">
                                                            {comment.Body}
                                                        </p>
                                                    )}
                                                </div>

                                                <div className="flex items-center gap-4 mt-1 pl-2">
                                                    <button className="text-[10px] font-medium text-(--muted) hover:text-red-500 transition-colors flex items-center gap-1">
                                                        <Heart className="w-3 h-3" />
                                                        {comment.ReactionsCount > 0 && comment.ReactionsCount}
                                                    </button>
                                                    {isOwner && !isEditing && (
                                                        <button
                                                            className="text-[10px] font-medium text-(--muted) hover:text-(--accent) transition-colors opacity-0 group-hover/comment-item:opacity-100"
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
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    )}

                    {/* Comment Composer */}
                    <div className="border-t border-(--border) p-4 bg-(--muted)/5">
                        <div className="flex gap-3">
                            <div className="w-8 h-8 rounded-full overflow-hidden bg-(--muted)/10 shrink-0">
                                {currentUser?.Avatar && (
                                    <img src={currentUser.Avatar} alt="My Avatar" className="w-full h-full object-cover" />
                                )}
                            </div>
                            <div className="flex-1 space-y-2">
                                <textarea
                                    ref={composerRef}
                                    value={draftComment}
                                    onChange={handleDraftChange}
                                    rows={1}
                                    className="w-full rounded-2xl border border-(--muted)/30 px-4 py-2.5 text-sm bg-transparent focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all resize-none min-h-[42px]"
                                    placeholder="Write a comment..."
                                />
                                <div className="flex justify-end gap-2">
                                    <button
                                        type="button"
                                        onClick={handleCancelComposer}
                                        className="px-3 py-1.5 text-xs font-medium text-(--muted) hover:text-(--foreground) hover:bg-(--muted)/10 rounded-full transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="button"
                                        disabled={!draftComment.trim()}
                                        onClick={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            // TODO: wire up createComment
                                        }}
                                        className="px-4 py-1.5 text-xs font-medium bg-(--accent) text-white hover:bg-(--accent-hover) rounded-full transition-colors disabled:opacity-50"
                                    >
                                        Reply
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
