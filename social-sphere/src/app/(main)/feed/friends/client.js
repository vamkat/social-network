"use client";

import { useRef } from "react";
import { fetchFeedPosts } from "@/services/posts/posts";
import FeedList from "@/components/feed/feed-list";
import CreatePost from "@/components/ui/create-post";

export default function FriendsFeedClient({ initialPosts }) {
    const addPostRef = useRef(null);

    const handlePostCreated = (newPost) => {
        if (addPostRef.current) {
            addPostRef.current(newPost);
        }
    };

    return (
        <div className="w-full py-8">
            <div className="max-w-7xl mx-auto px-6">
                <div className="flex gap-6">
                    {/* Left Sidebar - Reserved for future */}
                    <aside className="hidden xl:block w-48 shrink-0">
                        {/* Future: Navigation, shortcuts, etc */}
                    </aside>

                    {/* Main Feed */}
                    <main className="flex-1 max-w-2xl mx-auto">
                        <CreatePost onPostCreated={handlePostCreated} />
                        <div className="section-divider mb-6" />
                        <div className="mt-8 mb-6">
                            <h1 className="feed-title text-center">Friends Feed</h1>
                            <p className="text-center feed-subtitle">Updates from your friends</p>
                        </div>

                        <div className="section-divider mb-6" />

                        <FeedList
                            initialPosts={initialPosts}
                            fetchPosts={fetchFeedPosts}
                            onPostCreated={(handler) => { addPostRef.current = handler; }}
                        />
                    </main>

                    {/* Right Sidebar - Reserved for widgets */}
                    <aside className="hidden lg:block w-80 shrink-0">
                        {/* Future: Recommended users, trending, etc */}
                    </aside>
                </div>
            </div>
        </div>
    );
}
