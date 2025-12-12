"use client";

import { useState } from "react";
import { updateUserEmail, updateUserPassword } from "@/services/user/user-actions";
import { Eye, EyeOff } from "lucide-react";

export default function SecurityForm({ user }) {
    const [loadingEmail, setLoadingEmail] = useState(false);
    const [loadingPassword, setLoadingPassword] = useState(false);
    const [emailMessage, setEmailMessage] = useState(null);
    const [passwordMessage, setPasswordMessage] = useState(null);

    const [showCurrentPassword, setShowCurrentPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    async function handleEmailSubmit(formData) {
        setLoadingEmail(true);
        setEmailMessage(null);
        const email = formData.get("email");

        try {
            const result = await updateUserEmail(email);
            if (result.success) {
                setEmailMessage({ type: "success", text: result.message });
            } else {
                setEmailMessage({ type: "error", text: result.message });
            }
        } catch (error) {
            setEmailMessage({ type: "error", text: "Failed to update email" });
        } finally {
            setLoadingEmail(false);
        }
    }

    async function handlePasswordSubmit(formData) {
        setLoadingPassword(true);
        setPasswordMessage(null);

        const currentPassword = formData.get("currentPassword");
        const newPassword = formData.get("newPassword");
        const confirmPassword = formData.get("confirmPassword");

        if (newPassword !== confirmPassword) {
            setPasswordMessage({ type: "error", text: "New passwords do not match" });
            setLoadingPassword(false);
            return;
        }

        try {
            const result = await updateUserPassword(currentPassword, newPassword);
            if (result.success) {
                setPasswordMessage({ type: "success", text: result.message });
                // Optional: Reset form
            } else {
                setPasswordMessage({ type: "error", text: result.message });
            }
        } catch (error) {
            setPasswordMessage({ type: "error", text: "Failed to update password" });
        } finally {
            setLoadingPassword(false);
        }
    }

    return (
        <div className="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
            {/* Email Section */}
            <div className="space-y-6">
                <div>
                    <h3 className="text-lg font-semibold">Email Address</h3>
                    <p className="text-sm text-(--muted)">Update the email address associated with your account.</p>
                </div>

                <form action={handleEmailSubmit} className="space-y-4 max-w-md">
                    <div className="form-group">
                        <label className="form-label">Email</label>
                        <input
                            type="email"
                            name="email"
                            defaultValue={user?.email || "user@example.com"} // Mock default
                            className="form-input"
                            placeholder="your@email.com"
                            required
                        />
                    </div>

                    {emailMessage && (
                        <div className={`p-3 rounded-xl text-sm ${emailMessage.type === 'success' ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                            {emailMessage.text}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={loadingEmail}
                        className="btn btn-primary px-6"
                    >
                        {loadingEmail ? "Updating..." : "Update Email"}
                    </button>
                </form>
            </div>

            <div className="h-px bg-(--muted)/10" />

            {/* Password Section */}
            <div className="space-y-6">
                <div>
                    <h3 className="text-lg font-semibold">Change Password</h3>
                    <p className="text-sm text-(--muted)">Ensure your account is secure by using a strong password.</p>
                </div>

                <form action={handlePasswordSubmit} className="space-y-4 max-w-md">
                    <div className="form-group">
                        <label className="form-label">Current Password</label>
                        <div className="relative">
                            <input
                                type={showCurrentPassword ? "text" : "password"}
                                name="currentPassword"
                                className="form-input pr-10"
                                placeholder="••••••••"
                                required
                            />
                            <button
                                type="button"
                                onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                                className="form-toggle-btn"
                            >
                                {showCurrentPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                            </button>
                        </div>
                    </div>
                    <div className="form-group">
                        <label className="form-label">New Password</label>
                        <div className="relative">
                            <input
                                type={showNewPassword ? "text" : "password"}
                                name="newPassword"
                                className="form-input pr-10"
                                placeholder="••••••••"
                                required
                            />
                            <button
                                type="button"
                                onClick={() => setShowNewPassword(!showNewPassword)}
                                className="form-toggle-btn"
                            >
                                {showNewPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                            </button>
                        </div>
                    </div>
                    <div className="form-group">
                        <label className="form-label">Confirm New Password</label>
                        <div className="relative">
                            <input
                                type={showConfirmPassword ? "text" : "password"}
                                name="confirmPassword"
                                className="form-input pr-10"
                                placeholder="••••••••"
                                required
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

                    {passwordMessage && (
                        <div className={`p-3 rounded-xl text-sm ${passwordMessage.type === 'success' ? 'bg-green-500/10 text-green-600' : 'bg-red-500/10 text-red-600'}`}>
                            {passwordMessage.text}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={loadingPassword}
                        className="btn btn-primary px-6"
                    >
                        {loadingPassword ? "Updating..." : "Update Password"}
                    </button>
                </form>
            </div>
        </div>
    );
}
