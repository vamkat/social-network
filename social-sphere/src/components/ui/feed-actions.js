"use client";

import PostForm from "@/components/forms/PostForm";

// Always render the composer (no toggle CTA).
export default function FeedActions({ ctaProps = {}, onPostCreated, postFormProps = {} }) {
    const { title, subtitle } = ctaProps;

    return (
        <div className="space-y-2">
            {(title || subtitle) && (
                <div className="px-1">
                    {title && <div className="text-sm font-semibold text-(--foreground)">{title}</div>}
                    {subtitle && <div className="text-xs text-(--muted)">{subtitle}</div>}
                </div>
            )}
            <PostForm
                {...postFormProps}
                embed
                onPostCreated={(post) => {
                    if (onPostCreated) {
                        onPostCreated(post);
                    }
                }}
            />
        </div>
    );
}
