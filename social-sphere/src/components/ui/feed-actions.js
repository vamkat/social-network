"use client";

import { useState } from "react";
import PostForm from "@/components/forms/PostForm";
import FeedPostCTA from "./feed-post-creator";

// Renders the +Post CTA and expands it inline with the composer.
export default function FeedActions({ ctaProps = {}, onPostCreated, postFormProps = {} }) {
    const [showComposer, setShowComposer] = useState(false);

    return (
        <FeedPostCTA
            {...ctaProps}
            href={undefined} // force button behavior; no navigation
            onClick={() => setShowComposer((v) => !v)}
        >
            {showComposer && (
                <PostForm
                    {...postFormProps}
                    embed
                    onPostCreated={(post) => {
                        if (onPostCreated) {
                            onPostCreated(post);
                        }
                        setShowComposer(false);
                    }}
                    onCancel={() => setShowComposer(false)}
                />
            )}
        </FeedPostCTA>
    );
}
