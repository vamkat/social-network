"use client";

import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { useFormValidation } from "@/hooks/useFormValidation";
import { login } from "@/actions/auth/login";
import { useStore } from "@/store/store";
import { useRouter } from "next/router";

export default function LoginForm() {

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const loadUserProfile = useStore((state) => state.loadUserProfile);

    // Real-time validation hook
    const { errors: fieldErrors, validateField } = useFormValidation();

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);
        const identifier = formData.get("identifier");
        const password = formData.get("password");

        try {
            // call API to login
            const resp = await login({ identifier, password });

            // chceck err
            if (!resp.success || resp.error) {
                setError(resp.error || "Invalid credentials");
                setIsLoading(false);
                return;
            }

            // get user id from response, get user profile and store in localStorage
            const user = await loadUserProfile(resp.user_id);

            // check err
            if (!user.success) {
                setError("Login successful but failed to load profile");
                setIsLoading(false);
                return;
            }

            // all good
            window.location.href = "/feed/public";

        } catch (error) {
            console.error("Login exception:", err);
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    // Real-time validation handlers
    function handleFieldValidation(name, value) {
        switch (name) {
            case "identifier":
                validateField("identifier", value, (val) => {
                    if (!val.trim()) return "Email or Username is required.";
                    return null;
                });
                break;

            case "password":
                validateField("password", value, (val) => {
                    if (!val) return "Password is required.";
                    return null;
                });
                break;
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full space-y-6">
            {/* Email/Username Field */}
            <div>
                <label htmlFor="identifier" className="form-label pl-4 text-(--accent)">
                    Email or Username
                </label>
                <input
                    id="identifier"
                    name="identifier"
                    type="text"
                    required
                    className="form-input"
                    placeholder="Enter your email or username"
                    onChange={(e) => handleFieldValidation("identifier", e.target.value)}
                    disabled={isLoading}
                />
                {fieldErrors.identifier && (
                    <div className="form-error pl-4">{fieldErrors.identifier}</div>
                )}
            </div>

            {/* Password Field */}
            <div>
                <label htmlFor="password" className="form-label pl-4">
                    Password
                </label>
                <div className="relative group">
                    <input
                        id="password"
                        name="password"
                        type={showPassword ? "text" : "password"}
                        required
                        className="form-input pr-12"
                        placeholder="Enter your password"
                        onChange={(e) => handleFieldValidation("password", e.target.value)}
                        disabled={isLoading}
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="form-toggle-btn p-3 hover:text-(--accent)"
                        disabled={isLoading}
                    >
                        {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                    </button>
                </div>
                {fieldErrors.password && (
                    <div className="form-error">{fieldErrors.password}</div>
                )}
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
                className="w-1/2 mx-auto flex justify-center btn btn-primary pt-12"
            >
                {isLoading ? "Signing in..." : "Sign In"}
            </button>
        </form>
    );
}