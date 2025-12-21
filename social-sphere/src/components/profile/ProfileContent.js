"use client";

import { LogoutButton } from "@/components/LogoutButton";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import CreatePost from "@/components/ui/CreatePost";
import PostCard from "@/components/ui/PostCard";

export default function ProfileContent({ result, posts }) {
    // Handle error state
    if (!result.success) {
        return (
            <div className="flex flex-col items-center justify-center min-h-screen gap-4">
                <div className="text-red-500 text-lg font-medium">
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
            <div className="flex items-center justify-center min-h-screen">
                <div className="text-(--muted) text-lg">User not found</div>
            </div>
        );
    }

    console.log(result.user);

    // Render profile
    return (
        <div >
            <div>
                <ProfileHeader user={result.user} />
                {result.user.own_profile ? (
                    <div>
                        <div className="pt-3 flex flex-col px-70">
                            <CreatePost />
                        </div>
                        <div className="mt-8 mb-6">
                            <h1 className="text-center feed-title">My Feed</h1>
                            <p className="text-center feed-subtitle">What's happening in my sphere?</p>
                        </div>
                        <div className="section-divider mb-6" />
                    </div>
                ) : (
                    <div>
                        <div className="mt-8 mb-6">
                            <h1 className="text-center feed-title">{result.user.username}'s Feed</h1>
                            <p className="text-center feed-subtitle">What's happening in {result.user.username}'s sphere?</p>
                        </div>
                        <div className="section-divider mb-6" />
                    </div>
                )}
                <div className="pt-6 flex flex-col px-70">
                {posts?.length > 0 ? (
                    posts.map((post, index) => {
                        if (posts.length === index + 1) {
                            return (
                                <div key={`${post.ID}-${index}`}>
                                    <PostCard post={post} className="mb-6"/>
                                </div>
                            );
                        } else {
                            return <PostCard key={`${post.ID}-${index}`} post={post} />;
                        }
                    })
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                        <p className="text-muted text-center max-w-md">
                            Nothing.
                        </p>
                    </div>
                )}
                </div>
                <LogoutButton />
            </div>
        </div>
    );
}