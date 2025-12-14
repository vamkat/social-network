"use client";

import { useState } from "react";
import { Eye, EyeOff, Upload, X } from "lucide-react";
import { useFormValidation } from "@/hooks/useFormValidation";
import { register } from "@/services/auth/register";
import { useStore } from "@/store/store";

export default function RegisterForm() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [avatarPreview, setAvatarPreview] = useState(null);
    const [avatarName, setAvatarName] = useState(null);
    const [avatarFile, setAvatarFile] = useState(null);
    const loadUserProfile = useStore((state) => state.loadUserProfile);

    // Real-time validation state
    const { errors: fieldErrors, validateField } = useFormValidation();
    const [aboutCount, setAboutCount] = useState(0);

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);

        // Append avatar file if present
        if (avatarFile) {
            formData.set("avatar", avatarFile);
        }

        try {
            // call API to register
            const resp = await register(formData);

            // check err
            if (!resp.success || resp.error) {
                setError(resp.error || "Registration failed")
                setIsLoading(false);
                return;
            }

            // get user id from response, get user profile and store in localStorage
            const user = await loadUserProfile(resp.user_id);

            // check err
            if (!user.success) {
                setError("Registration successful but failed to load profile");
                setIsLoading(false);
                return;
            }

            // all good
            window.location.href = "/feed/public";

        } catch (error) {
            console.error("Registration exception:", err);
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    function handleAvatarChange(event) {
        const file = event.target.files[0];
        if (file) {
            setAvatarName(file.name);
            setAvatarFile(file);
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
        setAvatarFile(null);
    }

    // Real-time validation handlers
    function handleFieldValidation(name, value) {
        switch (name) {
            case "first_name":
                validateField("first_name", value, (val) => {
                    if (!val.trim()) return "First name is required.";
                    if (val.trim().length < 2) return "First name must be at least 2 characters.";
                    return null;
                });
                break;

            case "last_name":
                validateField("last_name", value, (val) => {
                    if (!val.trim()) return "Last name is required.";
                    if (val.trim().length < 2) return "Last name must be at least 2 characters.";
                    return null;
                });
                break;

            case "email":
                validateField("email", value, (val) => {
                    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                    if (!emailPattern.test(val.trim())) return "Please enter a valid email address.";
                    return null;
                });
                break;

            case "password":
                validateField("password", value, (val) => {
                    const strongPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$/;
                    if (val.length < 8) return "Password must be at least 8 characters.";
                    if (!strongPattern.test(val)) return "Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol.";
                    return null;
                });
                break;

            case "confirmPassword":
                validateField("confirmPassword", value, (val) => {
                    const passwordField = document.querySelector('input[name="password"]');
                    if (passwordField && val !== passwordField.value) return "Passwords do not match.";
                    return null;
                });
                break;

            case "date_of_birth":
                validateField("date_of_birth", value, (val) => {
                    if (!val) return "Date of birth is required.";
                    const today = new Date();
                    const birthDate = new Date(val);
                    let age = today.getFullYear() - birthDate.getFullYear();
                    const monthDiff = today.getMonth() - birthDate.getMonth();
                    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
                        age--;
                    }
                    if (age < 13 || age > 111) return "You must be between 13 and 111 years old.";
                    return null;
                });
                break;

            case "username":
                validateField("username", value, (val) => {
                    if (val.trim()) {
                        const usernamePattern = /^[A-Za-z0-9_.-]+$/;
                        if (val.trim().length < 4) return "Username must be at least 4 characters.";
                        if (!usernamePattern.test(val.trim())) return "Username can only use letters, numbers, dots, underscores, or dashes.";
                    }
                    return null;
                });
                break;

            case "about":
                setAboutCount(value.length);
                validateField("about", value, (val) => {
                    if (val.length > 400) return "About me must be at most 400 characters.";
                    return null;
                });
                break;
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
                {/* LEFT COLUMN - Account Info */}
                <div className="space-y-6">
                    <h3 className="text-lg font-semibold text-foreground mb-4">Account Information</h3>

                    {/* Name Fields */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label htmlFor="first_name" className="form-label pl-4">First Name <span className="text-red-500">*</span></label>
                            <input
                                id="first_name"
                                name="first_name"
                                type="text"
                                value="hello"
                                required
                                className="form-input"
                                placeholder="Jane"
                                onChange={(e) => handleFieldValidation("first_name", e.target.value)}
                            />
                            {fieldErrors.first_name && (
                                <div className="form-error">{fieldErrors.first_name}</div>
                            )}
                        </div>
                        <div>
                            <label htmlFor="last_name" className="form-label pl-4">Last Name <span className="text-red-500">*</span></label>
                            <input
                                id="last_name"
                                name="last_name"
                                type="text"
                                value="world"
                                required
                                className="form-input"
                                placeholder="Doe"
                                onChange={(e) => handleFieldValidation("last_name", e.target.value)}
                            />
                            {fieldErrors.last_name && (
                                <div className="form-error">{fieldErrors.last_name}</div>
                            )}
                        </div>
                    </div>

                    {/* Email */}
                    <div>
                        <label htmlFor="email" className="form-label pl-4">Email <span className="text-red-500">*</span></label>
                        <input
                            id="email"
                            name="email"
                            type="email"
                            value="hello@world.com"
                            required
                            className="form-input"
                            placeholder="jane@example.com"
                            onChange={(e) => handleFieldValidation("email", e.target.value)}
                        />
                        {fieldErrors.email && (
                            <div className="form-error">{fieldErrors.email}</div>
                        )}
                    </div>

                    {/* Password */}
                    <div>
                        <label htmlFor="password" className="form-label pl-4">Password <span className="text-red-500">*</span></label>
                        <div className="relative">
                            <input
                                id="password"
                                name="password"
                                type={showPassword ? "text" : "password"}
                                value="Hello12!"
                                required
                                className="form-input pr-12"
                                placeholder="HelloWorld123!"
                                minLength={8}
                                onChange={(e) => handleFieldValidation("password", e.target.value)}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="form-toggle-btn p-2"
                            >
                                {showPassword ? <EyeOff size={20} className="rounded-full" /> : <Eye size={20} className="rounded-full" />}
                            </button>
                        </div>
                        {fieldErrors.password && (
                            <div className="form-error">{fieldErrors.password}</div>
                        )}
                    </div>

                    {/* Confirm Password */}
                    <div>
                        <label htmlFor="confirmPassword" className="form-label pl-4">Confirm Password <span className="text-red-500">*</span></label>
                        <div className="relative">
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type={showConfirmPassword ? "text" : "password"}
                                value="Hello12!"
                                required
                                className="form-input pr-12"
                                placeholder="Confirm password"
                                minLength={8}
                                onChange={(e) => handleFieldValidation("confirmPassword", e.target.value)}
                            />
                            <button
                                type="button"
                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                className="form-toggle-btn p-2"
                            >
                                {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                            </button>
                        </div>
                        {fieldErrors.confirmPassword && (
                            <div className="form-error">{fieldErrors.confirmPassword}</div>
                        )}
                    </div>

                    {/* Date of Birth */}
                    <div>
                        <label htmlFor="date_of_birth" className="form-label pl-4">Date of Birth <span className="text-red-500">*</span></label>
                        <input
                            id="date_of_birth"
                            name="date_of_birth"
                            type="date"
                            value="2000-01-01"
                            required
                            className="form-input focus:outline-none"
                            onChange={(e) => handleFieldValidation("date_of_birth", e.target.value)}
                        />
                        {fieldErrors.date_of_birth && (
                            <div className="form-error">{fieldErrors.date_of_birth}</div>
                        )}
                    </div>
                </div>

                {/* RIGHT COLUMN - Profile Info */}
                <div className="space-y-6">
                    <h3 className="text-lg font-semibold text-foreground mb-4">Profile Details</h3>

                    {/* Avatar Upload */}
                    <div className="flex flex-col items-center gap-3">
                        <div className="relative">
                            <div className="avatar-container">
                                {avatarPreview ? (
                                    <img src={avatarPreview} alt="Avatar preview" className="avatar-image" />
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

                            {/* X button - only show when avatar is uploaded */}
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

                        <span className="text-sm text-muted">
                            {avatarName || "Upload Avatar (Optional)"}
                        </span>
                    </div>

                    {/* Username */}
                    <div>
                        <label htmlFor="username" className="form-label pl-4">Username (Optional)</label>
                        <input
                            id="username"
                            name="username"
                            type="text"
                            className="form-input"
                            placeholder="@janed"
                            onChange={(e) => handleFieldValidation("username", e.target.value)}
                        />
                        {fieldErrors.username && (
                            <div className="form-error">{fieldErrors.username}</div>
                        )}
                    </div>

                    {/* About Me */}
                    <div>
                        <div className="flex items-center justify-between mb-2">
                            <label htmlFor="about" className="form-label pl-4 mb-0">About Me (Optional)</label>
                            <span className="text-xs text-muted">
                                {aboutCount}/400
                            </span>
                        </div>
                        <textarea
                            id="about"
                            name="about"
                            rows={5}
                            maxLength={400}
                            className="form-input resize-none"
                            placeholder="Tell us a bit about yourself..."
                            onChange={(e) => handleFieldValidation("about", e.target.value)}
                        />
                        {fieldErrors.about && (
                            <div className="form-error">{fieldErrors.about}</div>
                        )}
                    </div>
                </div>
            </div>

            {/* Error Message */}
            {error && (
                <div className="form-error animate-fade-in mt-6 text-center pt-5">
                    {error}
                </div>
            )}

            {/* Submit Button */}
            <button
                type="submit"
                disabled={isLoading}
                className="w-1/3 mx-auto flex justify-center items-center btn btn-primary mt-3"
            >
                {isLoading ? "Creating Account..." : "Create Account"}
            </button>
        </form>
    );
}
