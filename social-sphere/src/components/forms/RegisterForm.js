"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, Upload, X } from "lucide-react";
import Link from "next/link";
import { register } from "@/actions/auth/auth";

export default function RegisterForm() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [avatarPreview, setAvatarPreview] = useState(null);
    const [avatarName, setAvatarName] = useState("");

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);

        // Client-side validation for password match
        const password = formData.get("password");
        const confirmPassword = formData.get("confirmPassword");

        if (password !== confirmPassword) {
            setError("Passwords do not match");
            setIsLoading(false);
            return;
        }

        // Append base64 avatar if present
        if (avatarPreview) {
            formData.set("avatar", avatarPreview);
        }

        try {
            const result = await register(formData);

            if (result.success) {
                console.log("Registration successful");
                router.push("/feed/public");
            } else {
                setError(result.error || "Registration failed");
                setIsLoading(false);
            }
        } catch (err) {
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    function handleAvatarChange(event) {
        const file = event.target.files[0];
        if (file) {
            setAvatarName(file.name);
            const reader = new FileReader();
            reader.onloadend = () => {
                setAvatarPreview(reader.result);
            };
            reader.readAsDataURL(file);
        }
    }

    function removeAvatar() {
        setAvatarPreview(null);
        setAvatarName("");
        // Reset file input if needed (requires ref)
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-lg space-y-6">
            {/* Avatar Upload */}
            <div className="flex flex-col items-center mb-6">
                <div className="avatar-container">
                    {avatarPreview ? (
                        <>
                            <img src={avatarPreview} alt="Avatar preview" className="avatar-image" />
                            <button
                                type="button"
                                onClick={removeAvatar}
                                className="avatar-remove-btn"
                            >
                                <X size={14} />
                            </button>
                        </>
                    ) : (
                        <Upload className="text-(--muted) w-8 h-8" />
                    )}
                    <input
                        type="file"
                        name="avatar"
                        accept="image/png, image/jpeg"
                        onChange={handleAvatarChange}
                        className="avatar-input"
                    />
                </div>
                <span className="avatar-label">
                    {avatarName || "Upload Avatar (Optional)"}
                </span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                    <label htmlFor="firstName" className="form-label">First Name</label>
                    <input id="firstName" name="firstName" type="text" required className="form-input" placeholder="Jane" />
                </div>
                <div className="form-group">
                    <label htmlFor="lastName" className="form-label">Last Name</label>
                    <input id="lastName" name="lastName" type="text" required className="form-input" placeholder="Doe" />
                </div>
            </div>

            <div className="form-group">
                <label htmlFor="email" className="form-label">Email</label>
                <input id="email" name="email" type="email" required className="form-input" placeholder="jane@example.com" />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                    <label htmlFor="password" className="form-label">Password</label>
                    <div className="relative">
                        <input
                            id="password"
                            name="password"
                            type={showPassword ? "text" : "password"}
                            required
                            className="form-input pr-10"
                            placeholder="Min. 8 chars"
                            minLength={8}
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="form-toggle-btn"
                        >
                            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>
                </div>
                <div className="form-group">
                    <label htmlFor="confirmPassword" className="form-label">Confirm Password</label>
                    <div className="relative">
                        <input
                            id="confirmPassword"
                            name="confirmPassword"
                            type={showConfirmPassword ? "text" : "password"}
                            required
                            className="form-input pr-10"
                            placeholder="Confirm password"
                            minLength={8}
                        />
                        <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="form-toggle-btn"
                        >
                            {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                    <label htmlFor="dob" className="form-label">Date of Birth</label>
                    <input id="dob" name="dateOfBirth" type="date" required className="form-input" />
                </div>
                <div className="form-group">
                    <label htmlFor="nickname" className="form-label">Nickname (Optional)</label>
                    <input id="nickname" name="nickname" type="text" className="form-input" placeholder="@janed" />
                </div>
            </div>

            <div className="form-group">
                <label htmlFor="aboutMe" className="form-label">About Me (Optional)</label>
                <textarea
                    id="aboutMe"
                    name="aboutMe"
                    rows={3}
                    maxLength={500}
                    className="form-input resize-none"
                    placeholder="Tell us a bit about yourself..."
                />
            </div>

            {error && (
                <div className="text-red-500 text-sm animate-fade-in text-center">
                    {error}
                </div>
            )}

            <button
                type="submit"
                disabled={isLoading}
                className="w-full btn btn-primary mt-4 disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {isLoading ? "Creating Account..." : "Create Account"}
            </button>
        </form>
    );
}
