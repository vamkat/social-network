"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import { Heart, MessageCircle, Pencil, Trash2, MoreHorizontal, Share2, Globe, Lock, Users, User, ChevronDown } from "lucide-react";
import PostImage from "./PostImage";
import DeleteConfirmModal from "./DeleteConfirmModal";
import { useStore } from "@/store/store";
import { editPost } from "@/actions/posts/edit-post";
import { deletePost } from "@/actions/posts/delete-post";
import { validateUpload } from "@/actions/auth/validate-upload";
import { getFollowers } from "@/actions/users/get-followers";

export default function PostCard({ post }) {
    const user = useStore((state) => state.user);
    const [image, setImage] = useState(post.image_url);
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [isExpanded, setIsExpanded] = useState(false);
    const [draftComment, setDraftComment] = useState("");
    const [postContent, setPostContent] = useState(post.post_body ?? "");
    const [isEditingPost, setIsEditingPost] = useState(false);
    const [postDraft, setPostDraft] = useState(post.post_body ?? "");
    const [editingCommentId, setEditingCommentId] = useState(null);
    const [editingText, setEditingText] = useState("");
    const [error, setError] = useState("");
    const [imageFile, setImageFile] = useState(null);
    const [imagePreview, setImagePreview] = useState(null);
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [removeExistingImage, setRemoveExistingImage] = useState(false);
    const [privacy, setPrivacy] = useState(post.audience || "everyone");
    const [isPrivacyOpen, setIsPrivacyOpen] = useState(false);
    const [selectedFollowers, setSelectedFollowers] = useState([]);
    const [followers, setFollowers] = useState([]);
    const [isLoadingFollowers, setIsLoadingFollowers] = useState(false);
    const composerRef = useRef(null);
    const cardRef = useRef(null);
    const fileInputRef = useRef(null);
    const dropdownRef = useRef(null);
    const router = useRouter();

    const isOwnPost = Boolean(
        user &&
        post?.post_user?.id &&
        String(post.post_user.id) === String(user.id)
    );

    // Close dropdown when clicking outside
    useEffect(() => {
        function handleClickOutside(event) {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsPrivacyOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    // useEffect(() => {
    //     const lastComment = post.LastComment;
    //     if (lastComment) {
    //         setComments([lastComment]);
    //         const total = Number(post.NumOfComments ?? 0);
    //         setHasMore(total > 1);
    //     } else {
    //         setComments([]);
    //         setHasMore(false);
    //     }
    //     setLoading(false);
    // }, [post.ID, post.LastComment, post.NumOfComments]);

    // useEffect(() => {
    //     if (isExpanded && composerRef.current) {
    //         composerRef.current.focus();
    //     }
    // }, [isExpanded]);

    // useEffect(() => {
    //     const handleClickOutside = (event) => {
    //         if (cardRef.current && !cardRef.current.contains(event.target)) {
    //             setIsExpanded(false);
    //         }
    //     };

    //     if (isExpanded) {
    //         document.addEventListener("mousedown", handleClickOutside);
    //     }

    //     return () => {
    //         document.removeEventListener("mousedown", handleClickOutside);
    //     };
    // }, [isExpanded]);

    const fetchFollowers = async () => {
        if (!user?.id || isLoadingFollowers) return;

        setIsLoadingFollowers(true);
        const followersData = await getFollowers({
            userId: user.id,
            limit: 100,
            offset: 0
        });
        setFollowers(followersData || []);
        setIsLoadingFollowers(false);
    };

    const handlePrivacySelect = (newPrivacy) => {
        setPrivacy(newPrivacy);
        setIsPrivacyOpen(false);
        if (newPrivacy !== "selected") {
            setSelectedFollowers([]);
        } else {
            fetchFollowers();
        }
    };

    const toggleFollower = (followerId) => {
        setSelectedFollowers((prev) =>
            prev.includes(followerId)
                ? prev.filter((id) => id !== followerId)
                : [...prev, followerId]
        );
    };

    const handleStartEditPost = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setPostDraft(postContent);
        setPrivacy(post.audience || "everyone");
        setIsEditingPost(true);
        setError("");
        setRemoveExistingImage(false);
    };

    const handleCancelEditPost = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setPostDraft(postContent);
        setPrivacy(post.audience || "everyone");
        setSelectedFollowers([]);
        setIsEditingPost(false);
        setImageFile(null);
        setImagePreview(null);
        setRemoveExistingImage(false);
        setError("");
        if (fileInputRef.current) {
            fileInputRef.current.value = "";
        }
    };

    const handleImageSelect = (e) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setImageFile(file);
        setError("");

        const reader = new FileReader();
        reader.onloadend = () => {
            setImagePreview(reader.result);
        };
        reader.readAsDataURL(file);
    };

    const handleRemoveImage = () => {
        setImageFile(null);
        setImagePreview(null);
        if (fileInputRef.current) {
            fileInputRef.current.value = "";
        }
    };

    const handleRemoveExistingImage = () => {
        setRemoveExistingImage(true);
    };

    const handleSaveEditPost = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (!postDraft.trim()) return;

        // Validate selected privacy
        if (privacy === "selected" && selectedFollowers.length === 0) {
            setError("Please select at least one follower for selected posts");
            return;
        }

        try {
            setError("");

            const editData = {
                post_id: post.post_id,
                post_body: postDraft.trim(),
                audience: privacy,
                audience_ids: privacy === "selected" ? selectedFollowers.map(id => parseInt(id)) : []
            };

            // Handle new image upload
            if (imageFile) {
                editData.image_name = imageFile.name;
                editData.image_size = imageFile.size;
                editData.image_type = imageFile.type;
            }
            // Handle explicit image removal
            else if (removeExistingImage) {
                editData.delete_image = true;
            }

            console.log("edit data to backend motherfuckers: ", editData);

            const resp = await editPost(editData);

            console.log("WHAT I GOT BACK", resp)

            if (!resp.success) {
                setError(resp.error || "Failed to edit post");
                return;
            }

            if (imageFile && resp.FileId && resp.UploadUrl) {
                console.log("Sending to url with id provided.... ")
                const uploadRes = await fetch(resp.UploadUrl, {
                    method: "PUT",
                    body: imageFile,
                });

                if (!uploadRes.ok) {
                    setError("Failed to upload image");
                    return;
                }

                const validateResp = await validateUpload(resp.FileId);
                if (!validateResp.success) {
                    setError("Failed to validate image upload");
                    return;
                }

                console.log("ID" , resp.FileId);
                console.log("New URL", validateResp.download_url)
                setImage(validateResp.download_url)
            }

            setPostContent(postDraft);
            setIsEditingPost(false);
            //window.location.reload();

        } catch (err) {
            console.error("Failed to edit post:", err);
            setError("Failed to edit post. Please try again.");
        }
    };

    const handleDeleteClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setShowDeleteModal(true);
    };

    const handleDeleteConfirm = async () => {
        setIsDeleting(true);
        setError("");

        try {
            const resp = await deletePost(post.post_id);

            if (!resp.success) {
                setError(resp.error || "Failed to delete post");
                setIsDeleting(false);
                setShowDeleteModal(false);
                return;
            }

            // Successfully deleted - force reload
            window.location.reload();
        

        } catch (err) {
            console.error("Failed to delete post:", err);
            setError("Failed to delete post. Please try again.");
            setIsDeleting(false);
            setShowDeleteModal(false);
        }
    };

    // const handleStartEdit = (comment) => {
    //     setEditingCommentId(comment.CommentId);
    //     setEditingText(comment.Body);
    // };

    // const handleCancelEdit = () => {
    //     setEditingCommentId(null);
    //     setEditingText("");
    // };

    // const handleSaveEdit = (commentId) => {
    //     if (!editingText.trim()) {
    //         handleCancelEdit();
    //         return;
    //     }
    //     setComments((prev) =>
    //         prev.map((c) =>
    //             c.CommentId === commentId ? { ...c, Body: editingText } : c
    //         )
    //     );
    //     handleCancelEdit();
    // };

    // const handleLoadMore = async (e) => {
    //     e.preventDefault();
    //     e.stopPropagation();
    //     if (loadingMore) return;
    //     setLoadingMore(true);

    //     try {
    //         const currentCount = comments.length;
    //         const limit = 2;
    //         const newComments = await fetchComments(post.ID, currentCount, limit);
    //         const existingIds = new Set(comments.map((c) => c.CommentId));
    //         const uniqueNew = newComments.filter((c) => !existingIds.has(c.CommentId));

    //         if (uniqueNew.length < limit) {
    //             setHasMore(false);
    //         }

    //         if (uniqueNew.length > 0) {
    //             setComments((prev) => [...uniqueNew.reverse(), ...prev]);
    //         } else {
    //             setHasMore(false);
    //         }
    //     } catch (error) {
    //         console.error("Failed to load more comments", error);
    //     } finally {
    //         setLoadingMore(false);
    //     }
    // };

    // const handleHeartClick = (e) => {
    //     e.preventDefault();
    //     e.stopPropagation();
    //     console.log("Heart clicked");
    //     // TODO: wire up heart reaction
    // };

    // const handleOpenPost = (e) => {
    //     const interactive = e.target.closest("button, a, textarea, input, select, option");
    //     if (interactive) return;
    //     e.preventDefault();
    //     router.push(`/posts/${post.ID}`);
    // };

    return (
        <div
            ref={cardRef}
            className="group bg-background border border-(--border) rounded-2xl overflow-hidden transition-all hover:border-(--muted)/40 hover:shadow-sm mb-6"
        >
            {/* Header */}
            <div className="p-5 flex items-start justify-between">
                <div className="flex items-center gap-3">
                    <Link href={`/profile/${post.post_user.id}`} className="w-10 h-10 rounded-full bg-(--muted)/10 flex items-center justify-center overflow-hidden shrink-0">
                        { post.post_user.avatar_url ? (<div className="w-10 h-10 rounded-full overflow-hidden border border-(--border)">
                            <img
                                src={post.post_user.avatar_url}
                                alt="Avatar"
                                className="w-full h-full object-cover"
                            />
                        </div>) : ( 
                            <User className="w-5 h-5 text-(--muted)" />)}
                    </Link>
                    <div>
                        <Link href={`/profile/${post.post_user.id}`}>
                            <h3 className="font-semibold text-foreground hover:underline decoration-2 underline-offset-2">
                                @{post.post_user.username}
                            </h3>
                        </Link>
                        <div className="flex items-center gap-2 text-xs text-(--muted) mt-0.5">
                            <span>{post.created_at}</span>
                        </div>
                    </div>
                </div>

                {isOwnPost && !isEditingPost && (
                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                            onClick={handleStartEditPost}
                            className="p-2 text-(--muted) hover:text-(--accent) hover:bg-(--accent)/5 rounded-full transition-colors"
                            title="Edit Post"
                        >
                            <Pencil className="w-4 h-4" />
                        </button>
                        <button
                            onClick={handleDeleteClick}
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

                        {/* Privacy Selector */}
                        <div className="relative" ref={dropdownRef}>
                            <button
                                type="button"
                                onClick={() => setIsPrivacyOpen(!isPrivacyOpen)}
                                className="flex items-center gap-1.5 bg-(--muted)/5 border border-(--border) rounded-full px-3 py-1.5 text-sm text-foreground hover:border-foreground focus:border-(--accent) transition-colors cursor-pointer"
                            >
                                <span className="capitalize">{privacy}</span>
                                <ChevronDown size={14} className={`transition-transform duration-200 ${isPrivacyOpen ? "rotate-180" : ""}`} />
                            </button>
                            {isPrivacyOpen && (
                                <div className="absolute top-full left-0 mt-1 w-32 bg-background border border-(--border) rounded-xl z-50 shadow-lg">
                                    <div className="flex flex-col p-1">
                                        <button
                                            type="button"
                                            onClick={() => handlePrivacySelect("everyone")}
                                            className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "everyone" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"}`}
                                        >
                                            Everyone
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => handlePrivacySelect("followers")}
                                            className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "followers" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"}`}
                                        >
                                            Followers
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => handlePrivacySelect("selected")}
                                            className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "selected" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"}`}
                                        >
                                            Selected
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Follower Selection for "Selected" Privacy */}
                        {privacy === "selected" && (
                            <div className="border border-(--border) rounded-xl p-4 space-y-2 bg-(--muted)/5">
                                <p className="text-xs font-medium text-(--muted)">
                                    Select followers who can see this post:
                                </p>
                                <div className="space-y-1.5 max-h-32 overflow-y-auto">
                                    {followers.length > 0 ? (
                                        followers.map((follower) => (
                                            <label
                                                key={follower.UserId}
                                                className="flex items-center gap-2 cursor-pointer hover:bg-(--muted)/10 rounded-lg px-2 py-1.5 transition-colors"
                                            >
                                                <input
                                                    type="checkbox"
                                                    checked={selectedFollowers.includes(String(follower.UserId))}
                                                    onChange={() => toggleFollower(String(follower.UserId))}
                                                    className="rounded border-gray-300"
                                                />
                                                <span className="text-sm">
                                                    @{follower.Username}
                                                </span>
                                            </label>
                                        ))
                                    ) : (
                                        <p className="text-xs text-(--muted) text-center py-2">
                                            {isLoadingFollowers ? "Loading followers..." : "No followers to select"}
                                        </p>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Image Preview for Edit - New Image */}
                        {imagePreview && (
                            <div className="relative inline-block">
                                <img
                                    src={imagePreview}
                                    alt="Upload preview"
                                    className="max-w-full max-h-64 rounded-xl border border-(--border)"
                                />
                                <button
                                    type="button"
                                    onClick={handleRemoveImage}
                                    className="absolute -top-2 -right-2 bg-background text-(--muted) hover:text-red-500 rounded-full p-1.5 border border-(--border) shadow-sm transition-colors"
                                >
                                    <Trash2 className="w-4 h-4" />
                                </button>
                            </div>
                        )}

                        {/* Existing Image in Edit Mode */}
                        {!imagePreview && post?.image_url && !removeExistingImage && (
                            <div className="relative inline-block">
                                <img
                                    src={image}
                                    alt="Post image"
                                    className="max-w-full max-h-64 rounded-xl border border-(--border)"
                                />
                                <button
                                    type="button"
                                    onClick={handleRemoveExistingImage}
                                    className="absolute -top-2 -right-2 bg-background text-(--muted) hover:text-red-500 rounded-full p-1.5 border border-(--border) shadow-sm transition-colors"
                                >
                                    <Trash2 className="w-4 h-4" />
                                </button>
                            </div>
                        )}

                        {/* Error Message */}
                        {error && (
                            <div className="text-red-500 text-sm bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-4 py-2.5">
                                {error}
                            </div>
                        )}

                        <div className="flex items-center justify-between gap-2">
                            <input
                                ref={fileInputRef}
                                type="file"
                                accept="image/jpeg,image/png,image/gif"
                                onChange={handleImageSelect}
                                className="hidden"
                            />
                            <button
                                type="button"
                                onClick={() => fileInputRef.current?.click()}
                                className="px-3 py-1.5 text-xs font-medium text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors"
                            >
                                Change Image
                            </button>

                            <div className="flex items-center gap-2">
                                <button
                                    type="button"
                                    className="px-3 py-1.5 text-xs font-medium text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors"
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
                    </div>
                ) : (
                    <Link href={`/posts/${post.post_id}`}>
                        <p className="text-[15px] leading-relaxed text-(--foreground)/90 whitespace-pre-wrap">
                            {postContent}
                        </p>
                    </Link>
                )}

                {/* Error Message (for non-edit operations like delete) */}
                {error && !isEditingPost && (
                    <div className="px-5 pb-3">
                        <div className="text-red-500 text-sm bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-4 py-2.5">
                            {error}
                        </div>
                    </div>
                )}
            </div>

            
            {post?.image_url && (
                <PostImage src={image} alt="He" />
            )}

            {/* Actions Footer */}
            <div className="px-5 py-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-6">
                        <button
                            className="flex items-center gap-2 text-(--muted) hover:text-red-500 transition-colors group/heart"

                        >
                            <Heart className={`w-5 h-5 transition-transform group-hover/heart:scale-110 ${post.liked_by_user ? "fill-red-500 text-red-500" : ""}`} />
                            <span className="text-sm font-medium">{post.reactions_count}</span>
                        </button>

                        <button
                            className={`flex items-center gap-2 transition-colors group/comment ${isExpanded ? "text-(--accent)" : "text-(--muted) hover:text-(--accent)"}`}

                        >
                            <MessageCircle className={`w-5 h-5 transition-transform group-hover/comment:scale-110 ${isExpanded ? "fill-(--accent)/10" : ""}`} />
                            <span className="text-sm font-medium">{post.comments_count}</span>
                        </button>
                    </div>
                </div>
            </div>

            {/* Expanded Section: Comments + Composer */}
            {/* {isExpanded && (
                <div className="animate-in fade-in slide-in-from-top-2 duration-200"> */}
            {/* Comments List
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
                                                                    className="text-xs text-(--muted) hover:text-foreground"
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
            {/* <div className="border-t border-(--border) p-4 bg-(--muted)/5">
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
                                        className="px-3 py-1.5 text-xs font-medium text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors"
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
            )} */}

            {/* Delete Confirmation Modal */}
            <DeleteConfirmModal
                isOpen={showDeleteModal}
                onClose={() => setShowDeleteModal(false)}
                onConfirm={handleDeleteConfirm}
                isDeleting={isDeleting}
            />
        </div>
    );
}
