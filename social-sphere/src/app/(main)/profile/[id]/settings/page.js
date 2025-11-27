"use client";

import { useState } from "react";
import { User, Shield, Lock } from "lucide-react";
import ProfileForm from "@/components/features/settings/profile-form";
import SecurityForm from "@/components/features/settings/security-form";
import PrivacyForm from "@/components/features/settings/privacy-form";
import { getUserByID } from "@/mock-data/users";
import { redirect } from "next/navigation";
import { use } from "react";

export default function SettingsPage({ params }) {
    const { id } = use(params);
    const [activeTab, setActiveTab] = useState("profile");

    // Mock user data - in real app this would come from auth context or server
    const user = getUserByID("1"); // Using mock user ID 1

    if (user.ID !== id) {
        // redirect to profile page
        redirect(`/profile/${id}`);
    }

    const tabs = [
        { id: "profile", label: "Edit Profile", icon: User },
        { id: "security", label: "Security", icon: Shield },
        { id: "privacy", label: "Privacy", icon: Lock },
    ];

    return (
        <div className="max-w-4xl mx-auto px-6 py-12">
            <div className="text-center mb-10">
                <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
                <p className="text-(--muted) mt-2">Manage your account settings and preferences.</p>
            </div>

            <div className="flex justify-center mb-8">
                <div className="inline-flex items-center p-1 rounded-xl bg-(--muted)/10">
                    {tabs.map((tab) => {
                        const Icon = tab.icon;
                        const isActive = activeTab === tab.id;
                        return (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`flex items-center gap-2 px-6 py-2.5 rounded-lg text-sm font-medium transition-all ${isActive
                                    ? "bg-(--background) text-(--foreground) shadow-sm"
                                    : "text-(--muted) hover:text-(--foreground)"
                                    }`}
                            >
                                <Icon className="w-4 h-4" />
                                {tab.label}
                            </button>
                        );
                    })}
                </div>
            </div>

            <div className="max-w-xl mx-auto">
                <div className="bg-(--background) rounded-2xl border border-(--muted)/10 p-8 shadow-sm">
                    {activeTab === "profile" && <ProfileForm user={user} />}
                    {activeTab === "security" && <SecurityForm user={user} />}
                    {activeTab === "privacy" && <PrivacyForm user={user} />}
                </div>
            </div>
        </div>
    );
}
