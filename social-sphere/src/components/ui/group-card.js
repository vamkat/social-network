"use client";

import Link from "next/link";

export default function GroupCard({ group }) {
    return (
        <div className="group flex flex-col relative border border-(--muted)/10 rounded-xl overflow-hidden hover:bg-(--muted)/5 transition-all duration-300">
            {/* Image Section */}
            <div className="aspect-video w-full bg-(--muted)/10 relative overflow-hidden">
                {group.Image ? (
                    <img
                        src={group.Image}
                        alt={group.Title[0]}
                        className="w-full h-full object-cover"
                    />
                ) : (
                    <div className="absolute inset-0 flex items-center justify-center text-(--muted)/20 text-4xl font-bold bg-linear-to-br from-(--muted)/5 to-(--muted)/20">
                        {group.Title.charAt(0)}
                    </div>
                )}
            </div>

            {/* Content Section */}
            <div className="p-5 flex flex-col flex-1">
                <div className="mb-3">
                    <Link href={`/groups/${group.ID}`}>
                        <h3 className="font-semibold text-lg text-(--foreground) hover:underline decoration-2 underline-offset-2 mb-1">
                            {group.Title}
                        </h3>
                    </Link>
                    <p className="text-(--foreground)/90 leading-relaxed text-[15px] line-clamp-2 h-11">
                        {group.Description}
                    </p>
                </div>

                <div className="flex items-center justify-between mt-auto pt-4 border-t border-(--muted)/10">
                    <div className="flex items-center gap-2 text-xs text-(--muted) font-medium">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                            <circle cx="9" cy="7" r="4" />
                            <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
                            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                        </svg>
                        {group.MembersNum} members
                    </div>

                    <Link
                        href={`/groups/${group.ID}`}
                        className="text-xs font-semibold bg-(--foreground) text-(--background) px-4 py-2 rounded-full hover:opacity-90 transition-opacity"
                    >
                        {group.IsMember ? 'View' : 'Join'}
                    </Link>
                </div>
            </div>
        </div>
    );
}
