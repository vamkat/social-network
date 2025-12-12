"use server";

import { getMockPosts, GetPostsByUserId } from "@/mock-data/posts";

export async function fetchPublicPosts(offset = 0, limit = 5) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts(offset, limit);
}

export async function fetchFeedPosts(offset = 0, limit = 5) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    return getMockPosts(offset, limit);
}

export async function fetchPostById(postId) {
    await new Promise((resolve) => setTimeout(resolve, 200));
    const post = getMockPosts(0, 1000).find((p) => String(p.ID) === String(postId));
    return post ?? null;
}

import { unstable_cache } from "next/cache";

export async function fetchUserPosts(userID, offset = 0, limit = 5) {
    const getCachedPosts = unstable_cache(
        async () => {
            console.log(`[CACHE MISS] Fetching posts for user ${userID}, offset ${offset}, limit ${limit}`);
            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 100));
            return GetPostsByUserId(userID, offset, limit);
        },
        [`user-posts-${userID}-${offset}-${limit}`],
        { revalidate: 60 }
    );

    return getCachedPosts();
}

/**
 * Create a new post
 * @param {FormData} formData - Form data containing post content, image, privacy, and allowedUsers
 * @returns {Promise<{success: boolean, post?: object, error?: string}>}
 */
export async function createPost(formData) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 500));

    try {
        // Extract form data
        const content = formData.get("content")?.trim();
        const privacy = formData.get("privacy") || "public";
        const imageFile = formData.get("image");
        const allowedUsersJson = formData.get("allowedUsers");

        // Validate content
        if (!content) {
            return { success: false, error: "Post content is required." };
        }

        if (content.length < 1 || content.length > 5000) {
            return { success: false, error: "Post content must be between 1 and 5000 characters." };
        }

        // Validate privacy
        const validPrivacyOptions = ["public", "friends", "private"];
        if (!validPrivacyOptions.includes(privacy)) {
            return { success: false, error: "Invalid privacy option." };
        }

        // Parse allowed users for private posts
        let allowedUsers = [];
        if (privacy === "private") {
            try {
                allowedUsers = allowedUsersJson ? JSON.parse(allowedUsersJson) : [];
                if (allowedUsers.length === 0) {
                    return { success: false, error: "Private posts must have at least one selected follower." };
                }
            } catch (e) {
                return { success: false, error: "Invalid allowed users data." };
            }
        }

        // Validate image if provided
        if (imageFile && imageFile.size > 0) {
            const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
            if (!allowedTypes.includes(imageFile.type)) {
                return { success: false, error: "Image must be JPEG, PNG, or GIF." };
            }

            const maxSize = 20 * 1024 * 1024; // 20MB
            if (imageFile.size > maxSize) {
                return { success: false, error: "Image must be less than 20MB." };
            }
        }

        // In a real implementation, this would:
        // 1. Upload image to storage service
        // 2. Send post data to backend API
        // 3. Return the created post from backend

        // For now, return a mock success response
        const newPost = {
            ID: `post-${Date.now()}`,
            Content: content,
            Visibility: privacy,
            // Image URL would come from backend after upload
            PostImage: imageFile && imageFile.size > 0 ? "/placeholder-image.jpg" : null,
            AllowedUsers: privacy === "private" ? allowedUsers : null,
        };

        return { success: true, post: newPost };
    } catch (error) {
        console.error("Error creating post:", error);
        return { success: false, error: "Failed to create post. Please try again." };
    }
}