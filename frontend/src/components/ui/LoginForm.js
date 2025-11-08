"use client";

import { useState } from "react";
import Link from "next/link";
import { Mail, Lock, Eye, EyeOff } from "lucide-react";

const baseInputClasses =
  "w-full rounded-2xl border border-(--color-border) bg-white/80 px-4 py-3 text-sm text-(--color-text) placeholder:text-(--color-text-muted) focus:outline-none focus:ring-2 focus:ring-(--color-accent) focus:border-transparent transition-shadow shadow-sm focus:shadow-md";

const LoginForm = () => {
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
  };

  return (
    <div className="relative w-full overflow-hidden rounded-4xl border border-(--color-border) bg-white/95 shadow-[0_18px_55px_-35px_rgba(31,27,22,0.85)]">
      <div className="absolute inset-x-8 top-0 h-px bg-linear-to-r from-transparent via-(--color-border-accent) to-transparent" />
      <div className="relative space-y-8 p-8">
        <div className="space-y-2 text-left">
          <p className="text-xs uppercase tracking-[0.25em] text-(--color-border-accent) text-center">Sign in</p>
          <h2 className="text-3xl font-extrabold text-(--color-text) text-center">Pick up where you left off</h2>
          
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <label className="flex flex-col gap-2 text-sm font-medium text-(--color-text)">
            Email
            <div className="relative">
              <Mail size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-(--color-text-muted)" />
              <input
                type="email"
                name="email"
                placeholder="you@socialsphere.io"
                required
                className={`${baseInputClasses} pl-11`}
              />
            </div>
          </label>

          <label className="flex flex-col gap-2 text-sm font-medium text-(--color-text)">
            Password
            <div className="relative">
              <Lock size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-(--color-text-muted)" />
              <input
                type={showPassword ? "text" : "password"}
                name="password"
                placeholder="••••••••"
                required
                className={`${baseInputClasses} pl-11 pr-12`}
              />
              <button
                type="button"
                onClick={() => setShowPassword((prev) => !prev)}
                className="absolute right-4 top-1/2 -translate-y-1/2 text-(--color-text-muted) transition-colors hover:text-(--color-text)"
                aria-label={showPassword ? "Hide password" : "Show password"}
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
          </label>

          <button
            type="submit"
            className="w-full rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white transition hover:bg-(--color-accent-hover)"
          >
            Continue to SocialSphere
          </button>

          <div className="text-center text-sm text-(--color-text-muted)">
            New to SocialSphere?{" "}
            <Link href="/register" className="font-semibold text-(--color-accent)">
              Create account
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};

export default LoginForm;
