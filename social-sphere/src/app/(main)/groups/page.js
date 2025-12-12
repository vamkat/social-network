import { getAllGroups, getMyGroups } from "@/services/groups/group-actions";
import Link from "next/link";
import GroupsSections from "@/components/features/groups/groups-sections";

export default async function GroupsPage() {
    const allGroups = await getAllGroups();
    const myGroups = await getMyGroups();
    const myGroupIds = new Set(myGroups.map((g) => g.ID));
    const availableGroups = allGroups.filter((group) => !group.IsMember && !myGroupIds.has(group.ID));

    return (
        <div className="w-full py-10">
            <div className="max-w-6xl mx-auto px-6 space-y-8">
                {/* Header */}
                <div className="flex items-center justify-between gap-4 flex-wrap">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">Groups</h1>
                        <p className="text-muted mt-1">Discover and join communities</p>
                    </div>
                    <Link href="/groups/create" className="btn btn-primary">
                        Create Group
                    </Link>
                </div>

                {/* Search Bar */}
                <div className="relative max-w-2xl">
                    <input
                        type="text"
                        placeholder="Search groups..."
                        className="w-full bg-(--muted)/5 border border-(--muted)/20 rounded-full px-6 py-3 pl-12 text-(--foreground) placeholder:text-(--muted) focus:outline-none focus:border-(--foreground) transition-colors"
                    />
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        className="absolute left-4 top-1/2 -translate-y-1/2 text-(--muted)"
                    >
                        <circle cx="11" cy="11" r="8" />
                        <path d="m21 21-4.3-4.3" />
                    </svg>
                </div>

                <GroupsSections myGroups={myGroups} availableGroups={availableGroups} />
            </div>
        </div>
    );
}

