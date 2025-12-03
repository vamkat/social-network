"use server";

import { validateRegistrationForm, validateLoginForm } from "@/utils/validation";

export async function login(formData) {

    // server side validation 
    const validation = validateLoginForm(formData);
    if (!validation.valid) {
        return { success: false, error: validation.error };
    }

    // extract fields for backend request
    const identifier = formData.get("identifier")?.trim();
    const password = formData.get("password");

    // prepare data for backend
    const backendFormData = new FormData();
    backendFormData.append("identifier", identifier);
    backendFormData.append("password", password);

    // login endpoint
    const loginEndpoint = process.env.LOGIN || "/login";
    const apiBase = process.env.API_BASE || "http://localhost:8081";

    // request login 
    const response = await fetch(`${apiBase}${loginEndpoint}`, {
        method: "POST",
        body: backendFormData,
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error("Login failed:", errorData);
        return { success: false, error: errorData.error || "Login failed. Please try again." };
    }

    return { success: true };
}

export async function register(formData) {
    // Extract avatar file
    const avatar = formData.get("avatar");

    // Shared validation
    const validation = validateRegistrationForm(formData, avatar);
    if (!validation.valid) {
        return { success: false, error: validation.error };
    }

    // Extract fields for backend request
    const firstName = formData.get("firstName")?.trim();
    const lastName = formData.get("lastName")?.trim();
    const email = formData.get("email")?.trim();
    const password = formData.get("password");
    const dateOfBirth = formData.get("dateOfBirth")?.trim();
    const nickname = formData.get("nickname")?.trim();
    const aboutMe = formData.get("aboutMe")?.trim();

    // Prepare data for backend
    const backendFormData = new FormData();
    backendFormData.append("first_name", firstName);
    backendFormData.append("last_name", lastName);
    backendFormData.append("email", email);
    backendFormData.append("password", password);
    backendFormData.append("date_of_birth", dateOfBirth);
    backendFormData.append("public", "true"); // Default to public

    if (nickname) backendFormData.append("username", nickname);
    if (aboutMe) backendFormData.append("about", aboutMe);
    if (avatar && avatar.size > 0) backendFormData.append("avatar", avatar);

    try {
        const apiBase = process.env.API_BASE || "http://localhost:8081";
        const registerEndpoint = process.env.REGISTER || "/register";

        const response = await fetch(`${apiBase}${registerEndpoint}`, {
            method: "POST",
            body: backendFormData,
            // Next.js/Node fetch automatically sets Content-Type to multipart/form-data with boundary
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error("Registration failed:", errorData);
            return { success: false, error: errorData.error || "Registration failed. Please try again." };
        }

        return { success: true };
    } catch (error) {
        console.error("Registration error:", error);
        return { success: false, error: "Network error. Please try again later." };
    }
}

