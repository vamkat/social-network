"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, Upload, X } from "lucide-react";
import { register } from "@/actions/auth/auth";
import { validateRegistrationForm } from "@/utils/validation";

export default function RegisterForm() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [avatarPreview, setAvatarPreview] = useState(null);
    const [avatarName, setAvatarName] = useState("");
    const [avatarFile, setAvatarFile] = useState(null);

    // Real-time validation state
    const [fieldErrors, setFieldErrors] = useState({});
    const [aboutMeCount, setAboutMeCount] = useState(0);

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);

        // Client-side validation using validation utilities
        const validation = validateRegistrationForm(formData, avatarFile);
        if (!validation.valid) {
            setError(validation.error);
            setIsLoading(false);
            return;
        }

        // Append avatar file if present
        if (avatarFile) {
            formData.set("avatar", avatarFile);
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
            setAvatarFile(file); // Store the file object
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
        // Reset file input if needed (requires ref)
    }

    // Real-time validation handlers
    function validateField(name, value) {
        const errors = { ...fieldErrors };

        switch (name) {
            case "firstName":
                if (!value.trim()) {
                    errors.firstName = "First name is required.";
                } else if (value.trim().length < 2) {
                    errors.firstName = "First name must be at least 2 characters.";
                } else {
                    delete errors.firstName;
                }
                break;

            case "lastName":
                if (!value.trim()) {
                    errors.lastName = "Last name is required.";
                } else if (value.trim().length < 2) {
                    errors.lastName = "Last name must be at least 2 characters.";
                } else {
                    delete errors.lastName;
                }
                break;

            case "email":
                const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                if (!emailPattern.test(value.trim())) {
                    errors.email = "Please enter a valid email address.";
                } else {
                    delete errors.email;
                }
                break;

            case "password":
                const strongPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$/;
                if (value.length < 8) {
                    errors.password = "Password must be at least 8 characters.";
                } else if (!strongPattern.test(value)) {
                    errors.password = "Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol.";
                } else {
                    delete errors.password;
                }
                break;

            case "confirmPassword":
                const passwordField = document.querySelector('input[name="password"]');
                if (passwordField && value !== passwordField.value) {
                    errors.confirmPassword = "Passwords do not match.";
                } else {
                    delete errors.confirmPassword;
                }
                break;

            case "dateOfBirth":
                if (!value) {
                    errors.dateOfBirth = "Date of birth is required.";
                } else {
                    const today = new Date();
                    const birthDate = new Date(value);
                    let age = today.getFullYear() - birthDate.getFullYear();
                    const monthDiff = today.getMonth() - birthDate.getMonth();
                    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
                        age--;
                    }
                    if (age < 13 || age > 111) {
                        errors.dateOfBirth = "You must be between 13 and 111 years old.";
                    } else {
                        delete errors.dateOfBirth;
                    }
                }
                break;

            case "nickname":
                if (value.trim()) {
                    const usernamePattern = /^[A-Za-z0-9_.-]+$/;
                    if (value.trim().length < 4) {
                        errors.nickname = "Username must be at least 4 characters.";
                    } else if (!usernamePattern.test(value.trim())) {
                        errors.nickname = "Username can only use letters, numbers, dots, underscores, or dashes.";
                    } else {
                        delete errors.nickname;
                    }
                } else {
                    delete errors.nickname;
                }
                break;

            case "aboutMe":
                setAboutMeCount(value.length);
                if (value.length > 400) {
                    errors.aboutMe = "About me must be at most 400 characters.";
                } else {
                    delete errors.aboutMe;
                }
                break;
        }

        setFieldErrors(errors);
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
                    <input
                        id="firstName"
                        name="firstName"
                        type="text"
                        required
                        className="form-input"
                        placeholder="Jane"
                        onChange={(e) => validateField("firstName", e.target.value)}
                    />
                    {fieldErrors.firstName && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.firstName}</div>
                    )}
                </div>
                <div className="form-group">
                    <label htmlFor="lastName" className="form-label">Last Name</label>
                    <input
                        id="lastName"
                        name="lastName"
                        type="text"
                        required
                        className="form-input"
                        placeholder="Doe"
                        onChange={(e) => validateField("lastName", e.target.value)}
                    />
                    {fieldErrors.lastName && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.lastName}</div>
                    )}
                </div>
            </div>

            <div className="form-group">
                <label htmlFor="email" className="form-label">Email</label>
                <input
                    id="email"
                    name="email"
                    type="email"
                    required
                    className="form-input"
                    placeholder="jane@example.com"
                    onChange={(e) => validateField("email", e.target.value)}
                />
                {fieldErrors.email && (
                    <div className="text-red-500 text-xs mt-1">{fieldErrors.email}</div>
                )}
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
                            placeholder="HelloWorld123!"
                            minLength={8}
                            onChange={(e) => validateField("password", e.target.value)}
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="form-toggle-btn"
                        >
                            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>
                    {fieldErrors.password && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.password}</div>
                    )}
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
                            onChange={(e) => validateField("confirmPassword", e.target.value)}
                        />
                        <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="form-toggle-btn"
                        >
                            {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                        </button>
                    </div>
                    {fieldErrors.confirmPassword && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.confirmPassword}</div>
                    )}
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                    <label htmlFor="dob" className="form-label">Date of Birth</label>
                    <input
                        id="dob"
                        name="dateOfBirth"
                        type="date"
                        required
                        className="form-input"
                        onChange={(e) => validateField("dateOfBirth", e.target.value)}
                    />
                    {fieldErrors.dateOfBirth && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.dateOfBirth}</div>
                    )}
                </div>
                <div className="form-group">
                    <label htmlFor="nickname" className="form-label">Nickname (Optional)</label>
                    <input
                        id="nickname"
                        name="nickname"
                        type="text"
                        className="form-input"
                        placeholder="@janed"
                        onChange={(e) => validateField("nickname", e.target.value)}
                    />
                    {fieldErrors.nickname && (
                        <div className="text-red-500 text-xs mt-1">{fieldErrors.nickname}</div>
                    )}
                </div>
            </div>

            <div className="form-group">
                <div className="flex items-center justify-between mb-1">
                    <label htmlFor="aboutMe" className="form-label">About Me (Optional)</label>
                    <span className="text-xs text-gray-500">
                        {aboutMeCount}/400
                    </span>
                </div>
                <textarea
                    id="aboutMe"
                    name="aboutMe"
                    rows={3}
                    maxLength={400}
                    className="form-input resize-none"
                    placeholder="Tell us a bit about yourself..."
                    onChange={(e) => validateField("aboutMe", e.target.value)}
                />
                {fieldErrors.aboutMe && (
                    <div className="text-red-500 text-xs mt-1">{fieldErrors.aboutMe}</div>
                )}
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
