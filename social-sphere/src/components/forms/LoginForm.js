"use client";

import { useState, useEffect } from "react";
import { Eye, EyeOff } from "lucide-react";
import { login } from "@/actions/auth/login";
import { useStore } from "@/store/store";
import LoadingThreeDotsJumping from '@/components/ui/LoadingDots';

export default function LoginForm() {

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const setUser = useStore((state) => state.setUser);
    const clearUser = useStore((state) => state.clearUser);

    // Clear any stale user data when login page loads
    useEffect(() => {
        clearUser();
    }, [clearUser]);

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true);
        setError("");

        const formData = new FormData(event.currentTarget);
        const email = formData.get("email");
        const password = formData.get("password");

        try {
            // call API to login
            const resp = await login({ email, password });
            // check err
            if (!resp.success || resp.error) {
                setError(resp.error || "Invalid credentials");
                setIsLoading(false);
                return;
            }

            // Store user data directly from login response
            setUser({
                id: resp.user_id,
                username: resp.username,
                avatar_url: resp.avatar_url || ""
            });

            // all good
            window.location.href = "/feed/public";

        } catch (error) {
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full space-y-6">
            {/* Email/Username Field */}
            <div>
                <label htmlFor="email" className="form-label pl-4 text-(--accent)">
                    Email
                </label>
                <input
                    id="email"
                    name="email"
                    type="email"
                    required
                    className="form-input"
                    placeholder="Enter your email"
                    disabled={isLoading}
                />
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
                className="w-1/2 mx-auto flex justify-center items-center btn btn-primary mt-12"
            >
                {isLoading ? <LoadingThreeDotsJumping /> : "Sign In"}
            </button>
        </form>
    );
}
