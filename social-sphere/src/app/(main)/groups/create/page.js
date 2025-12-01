import GroupForm from "@/components/forms/GroupForm";

export const metadata = {
    title: "Create Group | Social Sphere",
};

export default function CreateGroupPage() {
    return (
        <div className="max-w-3xl mx-auto space-y-6">
            <div className="space-y-2">
                <h1 className="text-3xl font-bold tracking-tight">Create a group</h1>
                <p className="text-(--muted)">Set up your community with a name, description, and visibility.</p>
            </div>

            <GroupForm />
        </div>
    );
}
