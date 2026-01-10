"use client";

import { useState } from "react";
import { Camera, Loader2, X } from "lucide-react";
import { updateProfileInfo } from "@/actions/profile/update-profile";
import { validateUpload } from "@/actions/auth/validate-upload";
import { validateProfileForm, validateImage } from "@/lib/validation";
import { useStore } from "@/store/store";

export default function ProfileForm({ user }) {
    const [isLoading, setIsLoading] = useState(false);
    const [message, setMessage] = useState({ success: false, text: null });
    const [imageErr, setImageErr] = useState(null);
    const [avatarFile, setAvatarFile] = useState(null);
    const [avatarPreview, setAvatarPreview] = useState(user?.avatar_url || null);
    const setUser = useStore((state) => state.setUser);
    const currentUser = useStore((state) => state.user);

    const handleAvatarChange = async (e) => {
        const file = e.target.files?.[0];
        if (!file) return;

        // Validate image file (type, size, dimensions)
        const validation = await validateImage(file);
        if (!validation.valid) {
            setImageErr(validation.error);
            return;
        }

        setImageErr(null);
        setAvatarFile(file);
        const reader = new FileReader();
        reader.onloadend = () => {
            setAvatarPreview(reader.result);
        };
        reader.readAsDataURL(file);
    };

    const removeAvatar = () => {
        setAvatarFile(null);
        setAvatarPreview(null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setMessage({ success: false, text: null });

        const formData = new FormData(e.currentTarget);
        const profileData = {
            username: formData.get("username")?.trim(),
            first_name: formData.get("firstName")?.trim(),
            last_name: formData.get("lastName")?.trim(),
            date_of_birth: formData.get("dateOfBirth"),
            about: formData.get("about")?.trim() || "",
        };

        // Add avatar metadata if a new file is selected
        if (avatarFile) {
            profileData.avatar_name = avatarFile.name;
            profileData.avatar_size = avatarFile.size;
            profileData.avatar_type = avatarFile.type;
        }

        try {
            // Validate profile data and avatar
            const validation = await validateProfileForm(profileData, avatarFile);
            if (!validation.valid) {
                setMessage({ success: false, text: validation.error });
                setIsLoading(false);
                return;
            }

            // Step 1: Update profile (with or without avatar metadata)
            const resp = await updateProfileInfo(profileData);

            if (!resp.success && resp.error) {
                setMessage({ success: false, text: resp.error });
                setIsLoading(false);
                return;
            }

            let newAvatarUrl = null;

            // Step 2: If avatar was provided, upload and validate
            if (avatarFile && resp.FileId && resp.UploadUrl) {
                try {
                    // Upload to MinIO
                    const uploadRes = await fetch(resp.UploadUrl, {
                        method: "PUT",
                        body: avatarFile
                    });

                    if (!uploadRes.ok) {
                        const errorText = await uploadRes.text();
                        console.error(`Storage error (${uploadRes.status}): ${errorText}`);
                        setMessage({ success: false, text: "Failed to upload avatar" });
                    } else {
                        // Validate upload
                        const validateResp = await validateUpload(resp.FileId);
                        if (!validateResp.success) {
                            setMessage({ success: false, text: validateResp.error || "Failed to validate upload" });
                        } else if (validateResp.download_url) {
                            newAvatarUrl = validateResp.download_url;
                        }
                    }
                } catch (uploadError) {
                    setMessage({ success: false, text: "Avatar upload failed, but profile updated" });
                }
            }

            // Step 3: Update store with new avatar URL
            if (currentUser) {
                setUser({
                    ...currentUser,
                    avatar_url: newAvatarUrl
                });
            }

            setMessage({ success: true, text: "Profile updated successfully!" });
            setAvatarFile(null);

        } catch (error) {
            console.error("Profile update error:", error);
            setMessage({ success: false, text: error.message || "An unexpected error occurred" });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-5 animate-in fade-in slide-in-from-bottom-4 duration-500">
            {/* Avatar Section */}
            <div className="flex flex-col items-center gap-4">
                <div className="relative group cursor-pointer">
                    <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-background shadow-xl flex items-center justify-center bg-gray-100">
                        {avatarPreview ? (
                            <img
                                src={avatarPreview}
                                alt="Profile"
                                className="w-full h-full object-cover"
                            />
                        ) : (
                            <Camera className="w-8 h-8 text-gray-400" />
                        )}
                    </div>
                    <label htmlFor="avatar-input" className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                        <Camera className="w-8 h-8 text-white" />
                    </label>

                    <input
                        id="avatar-input"
                        type="file"
                        onChange={handleAvatarChange}
                        className="hidden"
                        accept="image/jpeg,image/png,image/gif,image/webp"
                    />

                    {/* X button - only show when avatar exists */}
                    {avatarPreview && (
                        <button
                            type="button"
                            onClick={removeAvatar}
                            className="absolute -top-1 -right-1 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors z-10"
                        >
                            <X size={14} />
                        </button>
                    )}
                </div>
                {!imageErr ? (
                    <p className="text-sm text-(--muted)">Click to change avatar</p>
                ) : (
                    <p className="text-sm text-red-500">{imageErr}</p>
                )}
                
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

            {message.text && (
                <div className={`p-4 rounded-xl text-sm ${message.success ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                    {message.text}
                </div>
            )}

            <div className="flex justify-end">
                <button
                    type="submit"
                    disabled={isLoading}
                    className="btn btn-primary px-8 flex items-center gap-2"
                >
                    {isLoading && <Loader2 className="w-4 h-4 animate-spin" />}
                    {isLoading ? "Saving..." : "Save Changes"}
                </button>
            </div>
        </form>
    );
}
