"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff } from "lucide-react";

import { loginClient } from "@/actions/auth/login-client";

export default function LoginForm() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    // Real-time validation state
    const [fieldErrors, setFieldErrors] = useState({});

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);

        try {
            const result = await loginClient(formData);

            if (result.success) {
                console.log("Login successful");
                // Refresh to trigger AuthProvider to fetch user data
                router.refresh();
                router.push("/feed/public");
            } else {
                setError(result.error || "Invalid credentials");
                setIsLoading(false);
            }
        } catch (err) {
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    // Real-time validation handlers
    function validateField(name, value) {
        const errors = { ...fieldErrors };

        switch (name) {
            case "identifier":
                if (!value.trim()) {
                    errors.identifier = "Email or Username is required.";
                } else {
                    delete errors.identifier;
                }
                break;

            case "password":
                if (!value) {
                    errors.password = "Password is required.";
                } else {
                    delete errors.password;
                }
                break;
        }

        setFieldErrors(errors);
    }

    return (
        <form onSubmit={handleSubmit} className="w-full max-w-sm space-y-6">
            <div className="form-group">
                <label htmlFor="identifier" className="form-label">
                    Email or Username
                </label>
                <input
                    id="identifier"
                    name="identifier"
                    type="text"
                    required
                    className="form-input"
                    placeholder="Email/Username"
                    onChange={(e) => validateField("identifier", e.target.value)}
                />
                {fieldErrors.identifier && (
                    <div className="text-red-500 text-xs mt-1">{fieldErrors.identifier}</div>
                )}
            </div>

            <div className="form-group">
                <div className="flex items-center justify-between">
                    <label htmlFor="password" className="form-label">
                        Password
                    </label>
                    <Link
                        href="/forgot-password"
                        className="form-link"
                    >
                        Forgot?
                    </Link>
                </div>
                <div className="relative">
                    <input
                        id="password"
                        name="password"
                        type={showPassword ? "text" : "password"}
                        required
                        className="form-input pr-10"
                        placeholder="••••••••"
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

            {error && (
                <div className="text-red-500 text-sm animate-fade-in">
                    {error}
                </div>
            )}

            <button
                type="submit"
                disabled={isLoading}
                className="w-full btn btn-primary mt-8 disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {isLoading ? "Signing in..." : "Sign In"}
            </button>
        </form>
    );
}
