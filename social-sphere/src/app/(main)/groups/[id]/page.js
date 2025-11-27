import Link from "next/link";
import { notFound } from "next/navigation";
import { getGroupById, getGroupMembers, getGroupPosts } from "@/actions/groups/group-actions";
import PostFeed from "@/components/features/feed/post-feed";

export default async function GroupDetailPage({ params }) {
    const { id } = params;

    const [group, posts, members] = await Promise.all([
        getGroupById(id),
        getGroupPosts(id),
        getGroupMembers()
    ]);

    if (!group) {
        return notFound();
    }

    return (
        <div className="space-y-8">
            <div className="flex flex-col gap-4 rounded-2xl border border-(--muted)/10 bg-(--background) p-6 shadow-sm">
                <Link href="/groups" className="text-sm text-(--muted) hover:text-(--foreground) transition-colors w-fit">
                    ← Back to groups
                </Link>

                <div className="flex flex-col gap-3">
                    <div className="flex items-center justify-between gap-3">
                        <div>
                            <h1 className="text-3xl font-bold">{group.Title}</h1>
                            <p className="text-(--muted) mt-1">{group.Description}</p>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-(--muted)/10 text-xs font-semibold text-(--muted)">
                                {group.MembersNum} members
                            </span>
                            {group.PendingAcceptance && (
                                <span className="inline-flex items-center px-3 py-1 rounded-full bg-amber-500/10 text-amber-700 text-xs font-semibold">
                                    Request sent
                                </span>
                            )}
                            {group.IsOwner && (
                                <span className="inline-flex items-center px-3 py-1 rounded-full bg-emerald-500/10 text-emerald-700 text-xs font-semibold">
                                    You manage this group
                                </span>
                            )}
                        </div>
                    </div>

                    <div className="flex flex-wrap items-center gap-3 text-sm text-(--muted)">
                        <span>Owner ID: {group.OwnerID}</span>
                        <span>Created: {group.CreatedAt}</span>
                        <span>Members preview: {members.slice(0, 3).map((m) => m.Username).join(", ")}{members.length > 3 ? "…" : ""}</span>
                    </div>
                </div>
            </div>

            <div>
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-xl font-semibold">Group Posts</h2>
                    <span className="text-sm text-(--muted)">Mock pagination · infinite scroll</span>
                </div>
                <PostFeed posts={posts} pageSize={4} />
            </div>
        </div>
    );
}
