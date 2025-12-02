export default function ProfileStats({ stats }) {
    return (
        <div className="flex items-center gap-8 px-1">
            <button className="flex items-center gap-2 hover:opacity-80 transition-opacity group">
                <span className="text-lg font-bold text-(--foreground) group-hover:text-blue-500 transition-colors">
                    {stats.followers.toLocaleString()}
                </span>
                <span className="text-sm text-(--muted)">Followers</span>
            </button>

            <button className="flex items-center gap-2 hover:opacity-80 transition-opacity group">
                <span className="text-lg font-bold text-(--foreground) group-hover:text-blue-500 transition-colors">
                    {stats.following.toLocaleString()}
                </span>
                <span className="text-sm text-(--muted)">Following</span>
            </button>

            <button className="flex items-center gap-2 hover:opacity-80 transition-opacity group">
                <span className="text-lg font-bold text-(--foreground) group-hover:text-blue-500 transition-colors">
                    {stats.groups.toLocaleString()}
                </span>
                <span className="text-sm text-(--muted)">Groups</span>
            </button>
        </div>
    );
}
