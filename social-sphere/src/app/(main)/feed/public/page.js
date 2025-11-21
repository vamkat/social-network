export default function PublicFeedPage() {
    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold">Public Feed</h1>
            <p className="text-(--muted)">This is a placeholder for the public feed.</p>

            {/* Placeholder Content to test scrolling */}
            {[...Array(5)].map((_, i) => (
                <div key={i} className="p-6 rounded-2xl bg-(--muted)/5 border border-(--muted)/10 space-y-4">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-(--muted)/20" />
                        <div>
                            <div className="font-medium">User Name</div>
                            <div className="text-sm text-(--muted)">2 hours ago</div>
                        </div>
                    </div>
                    <div className="h-32 rounded-xl bg-(--muted)/10" />
                    <div className="flex gap-4 text-(--muted)">
                        <div className="w-8 h-8 rounded-full bg-(--muted)/10" />
                        <div className="w-8 h-8 rounded-full bg-(--muted)/10" />
                    </div>
                </div>
            ))}
        </div>
    );
}
