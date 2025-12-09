"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { signIn } from "next-auth/react";
import { Eye, EyeOff } from "lucide-react";
import { useFormValidation } from "@/hooks/useFormValidation";

export default function LoginForm() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);

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
            console.log("Step 1: Calling backend to set cookie...");
            
            // Step 1: Call backend directly to get and set the cookie
            const loginResult = await loginClient({ identifier, password });

            if (!loginResult.success) {
                setError(loginResult.error || "Invalid credentials");
                setIsLoading(false);
                return;
            }

            console.log("Step 2: Backend login successful, now signing into NextAuth...");
            
            // Step 2: Now sign into NextAuth (cookie is already set)
            const signInResult = await signIn("credentials", {
                redirect: false,
                userId: loginResult.user.UserId || loginResult.user.user_id,
                callbackUrl: "/feed/public"
            });

            console.log("Step 3: NextAuth sign in result:", signInResult);

            if (signInResult?.error) {
                console.error("NextAuth sign in error:", signInResult.error);
                setError("Session creation failed. Please try again.");
                setIsLoading(false);
            } else if (signInResult?.ok) {
                console.log("Step 4: Success! Redirecting...");
                window.location.href = "/feed/public";
            } else {
                setError("An unexpected error occurred");
                setIsLoading(false);
            }
        } catch (err) {
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
                <label htmlFor="identifier" className="form-label pl-4">
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
                    <div className="form-error">{fieldErrors.identifier}</div>
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
                        className="absolute right-3 top-1/2 -translate-y-1/2 p-2 rounded-full text-(--muted) group-focus-within:text-(--accent) hover:text-(--accent) transition-colors"
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
                <div className="form-error-box animate-fade-in">
                    {error}
                </div>
            )}

            {/* Submit Button */}
            <button
                type="submit"
                disabled={isLoading}
                className="w-1/2 mx-auto flex self-center justify-center btn btn-primary"
            >
                {isLoading ? "Signing in..." : "Sign In"}
            </button>
        </form>
    );
}