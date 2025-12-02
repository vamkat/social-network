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

        const firstName = formData.get("firstName")?.trim() || "";
        if (!firstName) {
            setError("First name is required.");
            setIsLoading(false);
            return;
        } else {
            if (firstName.length < 2) {
                setError("First name must be at least 2 characters.");
                setIsLoading(false);
                return;
            }
        }
        
        // Client-side validation for last name
        const lastName = formData.get("lastName")?.trim() || "";
        if (!lastName) {
            setError("Last name is required.");
            setIsLoading(false);
            return;
        } else {
            if (lastName.length < 2) {
                setError("Last name must be at least 2 characters.");
                setIsLoading(false);
                return;
            }
        }

        // Client-side validation for email
        const email = formData.get("email");
        const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailPattern.test(email)) {
            setError("Please enter a valid email address.");
            setIsLoading(false);
            return;
        }

        // Client-side validation for password match
        const password = formData.get("password");
        const confirmPassword = formData.get("confirmPassword");
        if (!password || !confirmPassword) {
            setError("Please enter both password and confirm password.");
            setIsLoading(false);
            return;
        }
        if (password.length < 8 ) {
            setError("Password must be at least 8 characters.");
            setIsLoading(false);
            return;
        }

        
        const strongPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$/;
        if (!strongPattern.test(password)) {
        setError("Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol.");
        setIsLoading(false);
        return;
        }
        // Check if password and confirm password match
        if (password !== confirmPassword) {
            setError("Passwords do not match");
            setIsLoading(false);
            return;
        }

        const dateOfBirth = formData.get("dateOfBirth")?.trim() || "";
        if (!dateOfBirth) {
            setError("Date of birth is required.");
            setIsLoading(false);
            return;
        } else {
            const age = calculateAge(dateOfBirth);
            if (age < 13 || age > 111) {
                setError("You must be between 13 and 111 years old.");
                setIsLoading(false);
                return;
            }
        }

        // Client-side validation for username
        const username = formData.get("nickname")?.trim() || "";
        if (username) {
            if (username.length < 4) {
                setError("Username must be at least 4 characters.");
                setIsLoading(false);
                return;
            }
            const safePattern = /^[A-Za-z0-9_.-]+$/; // basic “safe” set; adjust as needed
            if (!safePattern.test(username)) {
                setError("Username can only use letters, numbers, dots, underscores, or dashes.");
                setIsLoading(false);
                return;
            }
        }

        // Client-side validation for first name
        

        const aboutMe = formData.get("aboutMe")?.trim() || "";
        if (aboutMe) {
            if (aboutMe.length > 800) {
                setError("About me must be at most 800 characters.");
                setIsLoading(false);
                return;
            }
            const safePattern = /^[A-Za-z0-9_.-]+$/; // basic “safe” set; adjust as needed
            if (!safePattern.test(aboutMe)) {
                setError("This section can only use letters, numbers, dots, underscores, or dashes.");
                setIsLoading(false);
                return;
            }
        }

        

        // Append base64 avatar if present
        if (avatarPreview) {
            const allowedDataUrl = /^data:image\/(jpeg|png|gif);base64,[A-Za-z0-9+/]+=*$/i;
            if (!allowedDataUrl.test(avatarPreview)) {
                setError("Avatar must be base64 JPEG, PNG, or GIF.");
                setIsLoading(false);
                return;
            }
            formData.set("avatar", avatarPreview);handleAvatarChange
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
            const allowed = ["image/jpeg", "image/png", "image/gif"];
            if (!allowed.includes(file.type)) {
                setError("Avatar must be JPEG, PNG, or GIF.");
                return;
            }
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

    function calculateAge(dateOfBirth) {
        const today = new Date();
        const birthDate = new Date(dateOfBirth);
        let age = today.getFullYear() - birthDate.getFullYear();
        const month = today.getMonth() - birthDate.getMonth();
        if (month < 0 || (month === 0 && today.getDate() < birthDate.getDate())) {
            age--;
        }
        return age;
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
