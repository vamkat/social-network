"use client";

import { useState } from "react";
import { updateUserProfile } from "@/services/user/user-actions";
import { Camera } from "lucide-react";

export default function ProfileForm({ user }) {
    const [isLoading, setIsLoading] = useState(false);
    const [message, setMessage] = useState(null);

    async function handleSubmit(formData) {
        setIsLoading(true);
        setMessage(null);

        try {
            const result = await updateUserProfile(formData);
            if (result.success) {
                setMessage({ type: "success", text: result.message });
            } else {
                setMessage({ type: "error", text: result.message });
            }
        } catch (error) {
            setMessage({ type: "error", text: "Something went wrong" });
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form action={handleSubmit} className="space-y-5 animate-in fade-in slide-in-from-bottom-4 duration-500">
            {/* Avatar Section */}
            <div className="flex flex-col items-center gap-4">
                <div className="relative group cursor-pointer">
                    <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-(--background) shadow-xl">
                        <img
                            src={user?.Avatar || "/placeholder.jpg"}
                            alt="Profile"
                            className="w-full h-full object-cover"
                        />
                    </div>
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
                        defaultValue={user?.Username}
                        className="form-input"
                        placeholder="@username"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">Date of Birth</label>
                    <input
                        type="date"
                        name="dateOfBirth"
                        defaultValue={user?.DateOfBirth}
                        className="form-input"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">First Name</label>
                    <input
                        type="text"
                        name="firstName"
                        defaultValue={user?.firstName}
                        className="form-input"
                        placeholder="First Name"
                    />
                </div>
                <div className="form-group">
                    <label className="form-label">Last Name</label>
                    <input
                        type="text"
                        name="lastName"
                        defaultValue={user?.lastName}
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

            {message && (
                <div className={`p-4 rounded-xl text-sm ${message.type === 'success' ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                    {message.text}
                </div>
            )}

            <div className="flex justify-end">
                <button
                    type="submit"
                    disabled={isLoading}
                    className="btn btn-primary px-8"
                >
                    {isLoading ? "Saving..." : "Save Changes"}
                </button>
            </div>
        </form>
    );
}
