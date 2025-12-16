"use client";

import { useActionState } from "react";
import { updateProfileAction } from "@/actions/profile/settings";
import { Camera, Loader2 } from "lucide-react";

const initialState = {
    success: false,
    message: null,
};

export default function ProfileForm({ user }) {
    const [state, formAction, isPending] = useActionState(updateProfileAction, initialState);

    return (
        <form action={formAction} className="space-y-5 animate-in fade-in slide-in-from-bottom-4 duration-500">
            {/* Avatar Section */}
            <div className="flex flex-col items-center gap-4">
                <div className="relative group cursor-pointer">
                    <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-background shadow-xl">
                        {/* <img
                            src={user?.Avatar || "/placeholder.jpg"}
                            alt="Profile"
                            className="w-full h-full object-cover"
                        /> */}
                    </div>
                    {/* Note: Avatar upload logic not fully implemented in Server Action yet as it requires separate upload flow */}
                    <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <Camera className="w-8 h-8 text-white" />
                    </div>
                    <input type="file" name="avatar" className="hidden" accept="image/*" />
                </div>
                <p className="text-sm text-(--muted)">Click to change avatar</p>
            </div>

            <div className="grid gap-6 md:grid-cols-2">
                <div className="form-group">
                    <label className="form-label">Username</label>
                    <input
                        type="text"
                        name="username"
                        defaultValue={user?.username}
                        className="form-input"
                        placeholder="@username"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">Date of Birth</label>
                    <input
                        type="date"
                        name="dateOfBirth"
                        defaultValue={user?.date_of_birth ? user.date_of_birth.split('T')[0] : ''}
                        className="form-input"
                        placeholder="YYYY-MM-DD"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">First Name</label>
                    <input
                        type="text"
                        name="firstName"
                        defaultValue={user?.first_name}
                        className="form-input"
                        placeholder="First Name"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">Last Name</label>
                    <input
                        type="text"
                        name="lastName"
                        defaultValue={user?.last_name}
                        className="form-input"
                        placeholder="Last Name"
                    />
                </div>
            </div>

            <div className="form-group">
                <label className="form-label">About Me</label>
                <textarea
                    name="about"
                    defaultValue={user?.about}
                    rows={4}
                    className="form-input resize-none"
                    placeholder="Tell us about yourself..."
                />
            </div>

            {state.message && (
                <div className={`p-4 rounded-xl text-sm ${state.success ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                    {state.message}
                </div>
            )}

            <div className="flex justify-end">
                <button
                    type="submit"
                    disabled={isPending}
                    className="btn btn-primary px-8 flex items-center gap-2"
                >
                    {isPending && <Loader2 className="w-4 h-4 animate-spin" />}
                    {isPending ? "Saving..." : "Save Changes"}
                </button>
            </div>
        </form>
    );
}
