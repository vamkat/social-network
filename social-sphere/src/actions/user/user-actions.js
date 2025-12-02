"use server";

export async function updateUserProfile(formData) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // In a real app, we would validate and send data to backend
    console.log("Updating profile:", Object.fromEntries(formData));

    return { success: true, message: "Profile updated successfully" };
}

export async function updateUserEmail(email) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000));

    console.log("Updating email to:", email);

    return { success: true, message: "Email updated successfully" };
}

export async function updateUserPassword(currentPassword, newPassword) {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000));

    console.log("Updating password");

    // Mock validation
    if (currentPassword === "wrong") {
        return { success: false, message: "Incorrect current password" };
    }

    return { success: true, message: "Password updated successfully" };
}
