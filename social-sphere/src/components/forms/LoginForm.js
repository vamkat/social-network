"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff } from "lucide-react";

import { login } from "@/actions/auth/auth";

export default function LoginForm() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);

        try {
            const result = await login(formData);

            if (result.success) {
                console.log("Login successful");
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
                    placeholder="Email/Nickname"
                />
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
