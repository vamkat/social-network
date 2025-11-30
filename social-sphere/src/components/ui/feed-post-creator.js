"use client";

import Link from "next/link";

export default function FeedPostCTA({
    title = "Share something new",
    subtitle = "Post an update for your friends.",
    actionLabel = "+ Post",
    href,
    onClick,
    disabled = false,
    children,
}) {
    const Wrapper = href ? Link : "button";
    const wrapperProps = href
        ? { href }
        : {
            type: "button",
            onClick,
        };

    return (
        <div className="rounded-2xl border border-(--muted)/10 bg-(--muted)/5 px-4 py-3">
            <div className="flex items-center gap-3">
                <div className="flex-1 min-w-0">
                    <div className="text-sm font-semibold text-(--foreground) truncate">
                        {title}
                    </div>
                    <div className="text-xs text-(--muted) truncate">
                        {subtitle}
                    </div>
                </div>
                <Wrapper
                    {...wrapperProps}
                    className={`btn btn-primary px-4 py-2 text-xs ${disabled ? "opacity-60 cursor-not-allowed" : ""}`}
                    aria-disabled={disabled}
                >
                    {actionLabel}
                </Wrapper>
            </div>

            {children && (
                <div className="mt-3 pt-3 border-t border-(--muted)/10">
                    {children}
                </div>
            )}
        </div>
    );
}
