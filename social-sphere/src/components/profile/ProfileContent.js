"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import CreatePost from "@/components/ui/CreatePost";
import PostCard from "@/components/ui/PostCard";
import Container from "@/components/layout/Container";
import { Lock } from "lucide-react";
import { getUserPosts } from "@/actions/posts/get-user-posts";

export default function ProfileContent({ result, posts: initialPosts }) {
    const [posts, setPosts] = useState(initialPosts || []);
    const [offset, setOffset] = useState(10); // Start after the initial 10 posts
    // Only hasMore if we got a full batch of 10 posts
    const [hasMore, setHasMore] = useState((initialPosts || []).length >= 10);
    const [loading, setLoading] = useState(false);
    const observerTarget = useRef(null);
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

    // Check if viewer can see the profile content
    const canViewProfile = result.user.own_profile || result.user.public || result.user.viewer_is_following;

    const loadMorePosts = useCallback(async () => {
        if (loading || !hasMore || !canViewProfile) return;

        setLoading(true);
        try {
            const newPosts = await getUserPosts({ creatorId: result.user.user_id, limit: 5, offset });

            if (newPosts && newPosts.length > 0) {
                setPosts((prevPosts) => [...prevPosts, ...newPosts]);
                setOffset((prevOffset) => prevOffset + 5);

                // If we got fewer than 5 posts, we've reached the end
                if (newPosts.length < 5) {
                    setHasMore(false);
                }
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Failed to load more posts:", error);
        } finally {
            setLoading(false);
        }
    }, [offset, loading, hasMore, canViewProfile, result.user.user_id]);

    useEffect(() => {
        if (!canViewProfile) return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !loading) {
                    loadMorePosts();
                }
            },
            { threshold: 0.1 }
        );

        if (observerTarget.current) {
            observer.observe(observerTarget.current);
        }

        return () => {
            if (observerTarget.current) {
                observer.unobserve(observerTarget.current);
            }
        };
    }, [loadMorePosts, hasMore, loading, canViewProfile]);

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
            ) : canViewProfile ? (
                <div>
                    <div className="mt-8 mb-6">
                        <h1 className="text-center feed-title px-4">{result.user.username}'s Feed</h1>
                        <p className="text-center feed-subtitle px-4">What's happening in {result.user.username}'s sphere?</p>
                    </div>
                    <div className="section-divider mb-6" />
                </div>
            ) : null}

            <Container className="pt-6 pb-12">
                {canViewProfile ? (
                    posts?.length > 0 ? (
                        <div className="flex flex-col">
                            {posts.map((post, index) => (
                                <PostCard key={`${post.post_id}-${index}`} post={post} />
                            ))}

                            {/* Loading indicator */}
                            {hasMore && (
                                <div ref={observerTarget} className="flex justify-center py-8">
                                    {loading && (
                                        <div className="text-sm text-(--muted)">Loading more posts...</div>
                                    )}
                                </div>
                            )}

                            {/* End of feed message */}
                            {!hasMore && posts.length > 0 && (
                                <div className="text-center py-8 text-sm text-(--muted)">
                                    You've reached the end of the feed
                                </div>
                            )}
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                            <p className="text-muted text-center max-w-md px-4">
                                Nothing.
                            </p>
                        </div>
                    )
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                        <div className="w-16 h-16 rounded-full bg-(--muted)/10 flex items-center justify-center mb-4">
                            <Lock className="w-8 h-8 text-(--muted)" />
                        </div>
                        <h3 className="text-lg font-semibold text-foreground mb-2">
                            This profile is private
                        </h3>
                        <p className="text-(--muted) text-center max-w-md px-4">
                            Follow @{result.user.username} to see their posts and profile details.
                        </p>
                    </div>
                )}
            </Container>
        </div>
    );
}