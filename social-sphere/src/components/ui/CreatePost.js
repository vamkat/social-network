"use client";

import { useState, useRef, useEffect } from "react";
import { X, Image as ImageIcon, ChevronDown } from "lucide-react";
import Tooltip from "@/components/ui/Tooltip";
// import { isValidImage } from "@/lib/validation";

export default function CreatePost() {
    const [content, setContent] = useState("");
    const [privacy, setPrivacy] = useState("public");
    const [isPrivacyOpen, setIsPrivacyOpen] = useState(false);
    const [selectedFollowers, setSelectedFollowers] = useState([]);
    const [imageFile, setImageFile] = useState(null);
    const [imagePreview, setImagePreview] = useState(null);
    const [error, setError] = useState("");
    const fileInputRef = useRef(null);
    const dropdownRef = useRef(null);

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

    // const currentUser =  // Mock current user
    const MAX_CHARS = 5000;
    const MIN_CHARS = 1;

    // Mock followers list (will be replaced with real data)
    const mockFollowers = [
        { id: "2", name: "Georgia Toaka", username: "gtoaka" },
        { id: "4", name: "Watermelon Musk", username: "watermelon_musk" },
        { id: "5", name: "Trumpet Trump", username: "trumpet" },
    ];

    const handleImageSelect = (e) => {
        const file = e.target.files?.[0];
        if (!file) return;

        // Validate image file
        const validation = isValidImage(file);
        if (!validation.valid) {
            setError(validation.error);
            return;
        }

        setImageFile(file);
        setError("");

        // Create preview
        const reader = new FileReader();
        reader.onloadend = () => {
            setImagePreview(reader.result);
        };
        reader.readAsDataURL(file);
    };

    // const handleRemoveImage = () => {
    //     setImageFile(null);
    //     setImagePreview(null);
    //     if (fileInputRef.current) {
    //         fileInputRef.current.value = "";
    //     }
    // };

    // const handlePrivacySelect = (newPrivacy) => {
    //     setPrivacy(newPrivacy);
    //     setIsPrivacyOpen(false);
    //     if (newPrivacy !== "private") {
    //         setSelectedFollowers([]);
    //     }
    // };

    // const toggleFollower = (followerId) => {
    //     setSelectedFollowers((prev) =>
    //         prev.includes(followerId)
    //             ? prev.filter((id) => id !== followerId)
    //             : [...prev, followerId]
    //     );
    // };

    // const handleSubmit = async (e) => {
    //     e.preventDefault();
    //     setError("");

    //     // Validation
    //     if (!content.trim()) {
    //         setError("Post content is required");
    //         return;
    //     }

    //     if (content.length < MIN_CHARS) {
    //         setError(`Post must be at least ${MIN_CHARS} character`);
    //         return;
    //     }

    //     if (content.length > MAX_CHARS) {
    //         setError(`Post must be at most ${MAX_CHARS} characters`);
    //         return;
    //     }

    //     if (privacy === "private" && selectedFollowers.length === 0) {
    //         setError("Please select at least one follower for private posts");
    //         return;
    //     }

    //     try {
    //         // Prepare form data
    //         const formData = new FormData();
    //         formData.append("content", content.trim());
    //         formData.append("privacy", privacy);

    //         if (imageFile) {
    //             formData.append("image", imageFile);
    //         }

    //         if (privacy === "private") {
    //             formData.append("allowedUsers", JSON.stringify(selectedFollowers));
    //         }

    //         // Call server action
    //         const { createPost } = await import("@/services/posts/posts");
    //         const result = await createPost(formData);

    //         if (!result.success) {
    //             setError(result.error || "Failed to create post");
    //             return;
    //         }

    //         // Create optimistic post for UI
    //         const newPost = {
    //             ID: result.post.ID,
    //             BasicUserInfo: {
    //                 UserID: currentUser.ID,
    //                 Username: currentUser.Username,
    //                 Avatar: currentUser.Avatar,
    //             },
    //             Content: content.trim(),
    //             PostImage: imagePreview,
    //             CreatedAt: "Just now",
    //             NumOfComments: 0,
    //             NumOfHearts: 0,
    //             IsHearted: false,
    //             Visibility: privacy,
    //             LastComment: null,
    //         };

    //         // Notify parent component
    //         if (onPostCreated) {
    //             onPostCreated(newPost);
    //         }

    //         // Reset form
    //         setContent("");
    //         setPrivacy("public");
    //         setSelectedFollowers([]);
    //         handleRemoveImage();
    //     } catch (err) {
    //         console.error("Failed to create post:", err);
    //         setError("Failed to create post. Please try again.");
    //     }
    // };

    // const handleCancel = () => {
    //     setContent("");
    //     setPrivacy("public");
    //     setSelectedFollowers([]);
    //     handleRemoveImage();
    //     setError("");
    // };

    const charCount = content.length;
    const isOverLimit = charCount > MAX_CHARS;
    const isValid = content.trim().length >= MIN_CHARS && !isOverLimit;

    return (
        <div className="bg-background rounded-2xl p-3">
            <form>
                {/* Textarea with character counter */}
                <div className="relative">
                    <textarea
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        placeholder="What's on your mind?"
                        rows={3}
                        className="w-full bg-(--muted)/5 border border-(--border) rounded-xl px-2 py-3 pr-20 text-foreground placeholder:text-(--muted)/60 hover:border-foreground focus:outline-none focus:border-(--accent) focus:ring-2 focus:ring-(--accent)/10 transition-all resize-none"
                    />
                    {/* Character counter - bottom right */}
                    <div className="absolute bottom-3 right-3 text-xs">
                        <span
                            className={`font-medium ${isOverLimit
                                ? "text-red-500"
                                : charCount > MAX_CHARS * 0.9
                                    ? "text-orange-500"
                                    : "text-(--muted)/60"
                                }`}
                        >
                            {charCount > 0 && `${charCount}/${MAX_CHARS}`}
                        </span>
                    </div>
                </div>

                {/* Image Preview */}
                {/* {imagePreview && (
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
                            <X size={16} />
                        </button>
                    </div>
                )} */}

                {/* Error Message */}
                {error && (
                    <div className="text-red-500 text-sm bg-red-50 border border-red-200 rounded-lg px-4 py-2.5 animate-fade-in">
                        {error}
                    </div>
                )}

                {/* Follower Multi-Select for Private */}
                {privacy === "private" && (
                    <div className="border border-(--border) rounded-xl p-4 space-y-2 bg-(--muted)/5">
                        <p className="text-xs font-medium text-(--muted)">
                            Select followers who can see this post:
                        </p>
                        <div className="space-y-1.5 max-h-32 overflow-y-auto">
                            {mockFollowers.map((follower) => (
                                <label
                                    key={follower.id}
                                    className="flex items-center gap-2 cursor-pointer hover:bg-(--muted)/10 rounded-lg px-2 py-1.5 transition-colors"
                                >
                                    <input
                                        type="checkbox"
                                        checked={selectedFollowers.includes(follower.id)}
                                        onChange={() => toggleFollower(follower.id)}
                                        className="rounded border-gray-300"
                                    />
                                    <span className="text-sm">
                                        {follower.name} (@{follower.username})
                                    </span>
                                </label>
                            ))}
                        </div>
                    </div>
                )}

                {/* Bottom Controls Row */}
                <div className="flex items-center justify-between pt-2">
                    {/* Left side: Privacy and Image Upload */}
                    <div className="flex items-center gap-2">
                        {/* Privacy Selector */}
                        {/* Privacy Dropdown */}
                        <div className="relative" ref={dropdownRef}>
                            <Tooltip content={isPrivacyOpen ? "" : "Select privacy"}>
                                <button
                                    type="button"
                                    onClick={() => setIsPrivacyOpen(!isPrivacyOpen)}
                                    className="flex items-center gap-1.5 bg-(--muted)/5 border border-(--border) rounded-full px-3 py-1.5 text-sm text-foreground hover:border-foreground focus:border-(--accent) transition-colors cursor-pointer"
                                >
                                    <span className="capitalize">{privacy}</span>
                                    <ChevronDown size={14} className={`transition-transform duration-200 ${isPrivacyOpen ? "rotate-180" : ""}`} />
                                </button>
                            </Tooltip>
                            {isPrivacyOpen && (
                                <div className="absolute top-full left-0 mt-1 w-32 bg-background border border-(--border) rounded-xl z-50 animate-fade-in">

                                    <div className="flex flex-col p-1">
                                        <Tooltip content="Visible to everyone">
                                            <button
                                                type="button"
                                                onClick={() => handlePrivacySelect("public")}
                                                className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "public" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"
                                                    }`}
                                            >
                                                Public
                                            </button>
                                        </Tooltip>
                                        <Tooltip content="Visible to your friends only">
                                            <button
                                                type="button"
                                                onClick={() => handlePrivacySelect("friends")}
                                                className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "friends" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"
                                                    }`}
                                            >
                                                Friends
                                            </button>
                                        </Tooltip>
                                        <Tooltip content="Choose friends">
                                            <button
                                                type="button"
                                                onClick={() => handlePrivacySelect("private")}
                                                className={`w-full text-left px-3 py-1.5 text-sm rounded-lg transition-colors ${privacy === "private" ? "bg-(--muted)/10 font-medium" : "hover:bg-(--muted)/5 cursor-pointer"
                                                    }`}
                                            >
                                                Private
                                            </button>
                                        </Tooltip>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Image Upload Button */}
                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/jpeg,image/png,image/gif"
                            onChange={handleImageSelect}
                            className="hidden"
                        />
                        <Tooltip content="Upload image">
                            <button
                                type="button"
                                onClick={() => fileInputRef.current?.click()}
                                className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors"
                            >
                                <ImageIcon size={18} />
                                <span>Image</span>
                            </button>
                        </Tooltip>
                    </div>

                    {/* Right side: Submit and Cancel Buttons */}
                    <div className="flex items-center gap-2">
                        {(content || imageFile) && (
                            <>
                                <button
                                    type="button"
                                    onClick={handleCancel}
                                    className="px-4 py-1.5 text-sm text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={!isValid}
                                    className="px-5 py-1.5 text-sm font-medium bg-(--accent) text-white hover:bg-(--accent-hover) rounded-full disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                >
                                    Post
                                </button>
                            </>
                        )}
                    </div>
                </div>
            </form>
        </div>
    );
}
