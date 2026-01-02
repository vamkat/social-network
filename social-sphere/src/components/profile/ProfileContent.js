"use client";

import { ProfileHeader } from "@/components/profile/ProfileHeader";
import CreatePost from "@/components/ui/CreatePost";
import PostCard from "@/components/ui/PostCard";
import Container from "@/components/layout/Container";

export default function ProfileContent({ result, posts }) {
    // Handle error state
    if (!result.success) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen gap-4 px-4">
                <div className="text-red-500 text-lg font-medium text-center">
                    {result.error || "Failed to load profile"}
                </div>
                <button
                    onClick={() => window.location.reload()}
                    className="px-4 py-2 bg-(--accent) text-white rounded-lg hover:opacity-90 transition-opacity"
                >
                    Try Again
                </button>
            </div>
        );
    }

    // Handle no user
    if (!result.user) {
        return (
            <div className="flex items-center justify-center min-h-screen px-4">
                <div className="text-(--muted) text-lg">User not found</div>
            </div>
        );
    }

    // Render profile
    return (
        <div className="w-full">
            <ProfileHeader user={result.user} />

            {result.user.own_profile ? (
                <div>
                    <Container className="pt-6 md:pt-10">
                        <CreatePost />
                    </Container>
                    <div className="mt-8 mb-6">
                        <h1 className="text-center feed-title px-4">My Feed</h1>
                        <p className="text-center feed-subtitle px-4">What's happening in my sphere?</p>
                    </div>
                    <div className="section-divider mb-6" />
                </div>
            ) : (
                <div>
                    <div className="mt-8 mb-6">
                        <h1 className="text-center feed-title px-4">{result.user.username}'s Feed</h1>
                        <p className="text-center feed-subtitle px-4">What's happening in {result.user.username}'s sphere?</p>
                    </div>
                    <div className="section-divider mb-6" />
                </div>
            )}

            <Container className="pt-6 pb-12">
                {posts?.length > 0 ? (
                    <div className="flex flex-col">
                        {posts.map((post, index) => (
                            <PostCard key={`${post.post_id}-${index}`} post={post} />
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                        <p className="text-muted text-center max-w-md px-4">
                            Nothing.
                        </p>
                    </div>
                )}
            </Container>
        </div>
    );
}