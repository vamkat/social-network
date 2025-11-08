"use client";

import { useState } from "react";
import Link from "next/link";
import { Mail, Lock, User, Calendar, Image as ImageIcon, AtSign } from "lucide-react";

const baseInputClasses =
  "w-full rounded-2xl border border-(--color-border) bg-white/80 px-4 py-3 text-sm text-(--color-text) placeholder:text-(--color-text-muted) focus:outline-none focus:ring-2 focus:ring-(--color-accent) focus:border-transparent transition-shadow shadow-sm focus:shadow-md";

const FieldLabel = ({ label, optional = false }) => (
  <span className="flex items-center justify-between text-sm font-medium text-(--color-text)">
    {label}
    {optional && <span className="text-xs text-(--color-text-muted)">Optional</span>}
  </span>
);

const InputField = ({ icon: Icon, optional = false, className = "", label, ...props }) => (
  <label className="flex flex-col gap-2">
    <FieldLabel label={label} optional={optional} />
    <div className="relative">
      {Icon && (
        <Icon size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-(--color-text-muted)" />
      )}
      <input
        {...props}
        required={props.required ?? !optional}
        className={`${baseInputClasses} ${Icon ? "pl-11" : "pl-4"} ${className}`}
      />
    </div>
  </label>
);

const TextareaField = ({ optional = false, label, className = "", ...props }) => (
  <label className="flex flex-col gap-2">
    <FieldLabel label={label} optional={optional} />
    <textarea
      {...props}
      required={props.required ?? !optional}
      className={`${baseInputClasses} min-h-[120px] resize-none ${className}`}
    />
  </label>
);

const RegistrationForm = () => {
  const [avatarName, setAvatarName] = useState("");

  const handleSubmit = (event) => {
    event.preventDefault();
  };

  const handleAvatarChange = (event) => {
    const file = event.target.files?.[0];
    setAvatarName(file ? file.name : "");
  };

  return (
    <div className="relative w-full overflow-hidden rounded-[36px] border border-(--color-border) bg-white/95 shadow-[0_20px_60px_-35px_rgba(31,27,22,0.9)]">
      <div className="absolute inset-x-10 top-0 h-px bg-linear-to-r from-transparent via-(--color-border-accent) to-transparent" />
      <div className="relative space-y-8 p-8 md:p-10">
        <div className="space-y-3">
          <p className="text-xs uppercase tracking-[0.25em] text-(--color-border-accent)">Account details</p>
          <h2 className="text-3xl font-extrabold text-(--color-text)">Just a few basics</h2>
          <p className="text-sm text-(--color-text-muted)">
            We'll use this to build your SocialSphere identity. Nothing is shared until you hit publish.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <InputField label="First Name" name="firstName" placeholder="Jamie" icon={User} />
            <InputField label="Last Name" name="lastName" placeholder="Doe" icon={User} />
          </div>

          <InputField label="Email" name="email" type="email" placeholder="jamie@samplesphere.com" icon={Mail} />

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <InputField
              label="Password"
              name="password"
              type="password"
              placeholder="Create a secure password"
              icon={Lock}
            />
            <InputField
              label="Confirm Password"
              name="confirmPassword"
              type="password"
              placeholder="Re-enter your password"
              icon={Lock}
            />
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <InputField label="Date of Birth" name="dateOfBirth" type="date" icon={Calendar} />
            <InputField
              label="Nickname"
              name="nickname"
              placeholder="What should people call you?"
              optional
              icon={AtSign}
            />
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label className="flex flex-col gap-2">
              <FieldLabel label="Avatar / Image" optional />
              <div className="relative rounded-2xl border border-dashed border-(--color-border) bg-(--color-surface)/70 px-4 py-4 text-sm shadow-inner">
                <input
                  type="file"
                  accept="image/png,image/jpeg,image/gif"
                  name="avatar"
                  onChange={handleAvatarChange}
                  className="absolute inset-0 h-full w-full cursor-pointer opacity-0"
                  aria-label="Upload avatar"
                />
                <div className="flex items-start gap-3">
                  <ImageIcon className="mt-1 text-(--color-text-muted)" size={20} />
                  <div>
                    <p className="font-medium text-(--color-text)">{avatarName || "Upload avatar"}</p>
                    <p className="text-xs text-(--color-text-muted)">PNG, JPG, or GIF — max 5MB</p>
                  </div>
                </div>
              </div>
            </label>

            <div className="md:col-span-2">
              <TextareaField
                label="About Me"
                name="about"
                placeholder="Share a short intro, your passions, or what you're looking for."
                optional
                className="min-h-[150px]"
              />
            </div>
          </div>

          <div className="rounded-2xl bg-(--color-surface)/70 px-4 py-3 text-xs text-(--color-text-muted)">
            By continuing you agree to our {" "}
            <Link href="/terms" className="font-medium text-(--color-accent)">
              community guidelines
            </Link>{" "}
            and {" "}
            <Link href="/privacy" className="font-medium text-(--color-accent)">
              privacy promise
            </Link>
            . You can update these preferences anytime.
          </div>

          <button
            type="submit"
            className="w-full rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white transition hover:bg-(--color-accent-hover)"
          >
            Create my profile
          </button>

          <div className="text-center text-sm text-(--color-text-muted)">
            Already part of SocialSphere?{" "}
            <Link href="/login" className="font-semibold text-(--color-accent)">
              Log in
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};

export default RegistrationForm;
