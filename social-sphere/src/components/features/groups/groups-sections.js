"use client";

import { useMemo, useState } from "react";
import GroupCard from "@/components/ui/group-card";

const PAGE_SIZE = 6;

export default function GroupsSections({ myGroups, availableGroups }) {
    const hasMyGroups = myGroups.length > 0;
    const hasAvailableGroups = availableGroups.length > 0;

    const [myLimit, setMyLimit] = useState(Math.min(PAGE_SIZE, myGroups.length));
    const [availableLimit, setAvailableLimit] = useState(Math.min(PAGE_SIZE, availableGroups.length));
    const [showAvailable, setShowAvailable] = useState(!hasMyGroups);

    const visibleMyGroups = useMemo(
        () => myGroups.slice(0, myLimit),
        [myGroups, myLimit]
    );

    const visibleAvailableGroups = useMemo(
        () => availableGroups.slice(0, availableLimit),
        [availableGroups, availableLimit]
    );

    const canLoadMoreMy = myLimit < myGroups.length;
    const canLoadMoreAvailable = availableLimit < availableGroups.length;
    const canShowLessMy = myLimit > PAGE_SIZE;
    const canShowLessAvailable = availableLimit > PAGE_SIZE;

    const handleLoadMoreMy = () => {
        setMyLimit((prev) => Math.min(prev + PAGE_SIZE, myGroups.length));
    };

    const handleShowLessMy = () => {
        setMyLimit(Math.min(PAGE_SIZE, myGroups.length));
    };

    const handleLoadMoreAvailable = () => {
        setAvailableLimit((prev) => Math.min(prev + PAGE_SIZE, availableGroups.length));
    };

    const handleShowLessAvailable = () => {
        setAvailableLimit(Math.min(PAGE_SIZE, availableGroups.length));
    };

    return (
        <div className="space-y-10">
            {hasMyGroups && (
                <section className="bg-(--muted)/5 border border-(--muted)/10 rounded-2xl px-5 py-6 shadow-sm">
                    <div className="flex items-center justify-between flex-wrap gap-3 mb-5">
                        <h2 className="text-xl font-semibold">My groups</h2>
                        <p className="text-sm text-(--muted)">{myGroups.length} total</p>
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                        {visibleMyGroups.map((group) => (
                            <GroupCard key={group.ID} group={group} />
                        ))}
                    </div>

                    {(canLoadMoreMy || canShowLessMy) && (
                        <div className="flex justify-center mt-6 gap-3">
                            {canLoadMoreMy && (
                                <button
                                    type="button"
                                    onClick={handleLoadMoreMy}
                                    className="px-5 py-2 rounded-full border border-(--muted)/30 text-sm font-medium hover:border-(--foreground) transition-colors"
                                >
                                    Load more
                                </button>
                            )}
                            {canShowLessMy && (
                                <button
                                    type="button"
                                    onClick={handleShowLessMy}
                                    className="px-5 py-2 rounded-full border border-(--muted)/30 text-sm font-medium hover:border-(--foreground) transition-colors"
                                >
                                    Show less
                                </button>
                            )}
                        </div>
                    )}
                </section>
            )}

            {hasAvailableGroups && !showAvailable && hasMyGroups && (
                <div className="flex justify-center">
                    <button
                        type="button"
                        onClick={() => setShowAvailable(true)}
                        className="px-5 py-2 rounded-full border border-(--muted)/30 text-sm font-semibold bg-(--background) hover:border-(--foreground) transition-colors"
                    >
                        Show groups to join
                    </button>
                </div>
            )}

            {hasAvailableGroups && showAvailable && (
                <section className="bg-(--muted)/5 border border-(--muted)/10 rounded-2xl px-5 py-6 shadow-sm">
                    <div className="flex items-center justify-between flex-wrap gap-3 mb-5">
                        <h2 className="text-xl font-semibold">Groups you can join</h2>
                        <p className="text-sm text-(--muted)">{availableGroups.length} available</p>
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                        {visibleAvailableGroups.map((group) => (
                            <GroupCard key={group.ID} group={group} />
                        ))}
                    </div>

                    {(canLoadMoreAvailable || canShowLessAvailable) && (
                        <div className="flex justify-center mt-6 gap-3">
                            {canLoadMoreAvailable && (
                                <button
                                    type="button"
                                    onClick={handleLoadMoreAvailable}
                                    className="px-5 py-2 rounded-full border border-(--muted)/30 text-sm font-medium hover:border-(--foreground) transition-colors"
                                >
                                    Load more
                                </button>
                            )}
                            {canShowLessAvailable && (
                                <button
                                    type="button"
                                    onClick={handleShowLessAvailable}
                                    className="px-5 py-2 rounded-full border border-(--muted)/30 text-sm font-medium hover:border-(--foreground) transition-colors"
                                >
                                    Show less
                                </button>
                            )}
                        </div>
                    )}
                </section>
            )}

            {!hasMyGroups && !hasAvailableGroups && (
                <div className="text-center text-(--muted) text-sm">
                    No groups to display yet.
                </div>
            )}
        </div>
    );
}
