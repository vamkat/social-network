"use server";

export async function login(formData) {
    const identifier = formData.get("identifier");
    const password = formData.get("password");

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // Mock validation
    if (!identifier || !password) {
        return { success: false, error: "identifier and password are required" };
    }

    // Mock authentication logic
    // In a real app, this would make a server-to-server API call to the Golang backend
    if (identifier && password) {
        console.log("Server Action: Login successful", { identifier });
        return { success: true };
    } else {
        return { success: false, error: "Invalid credentials" };
    }
}

export async function register(formData) {
    const email = formData.get("email");
    const password = formData.get("password");
    const confirmPassword = formData.get("confirmPassword");
    const firstName = formData.get("firstName");
    const lastName = formData.get("lastName");
    const dateOfBirth = formData.get("dateOfBirth");
    const nickname = formData.get("nickname");
    const aboutMe = formData.get("aboutMe");
    const avatar = formData.get("avatar"); // Base64 string

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // Mock validation
    if (!email || !password || !firstName || !lastName || !dateOfBirth) {
        return { success: false, error: "All required fields must be filled" };
    }

    if (password !== confirmPassword) {
        return { success: false, error: "Passwords do not match" };
    }

    if (password.length < 8) {
        return { success: false, error: "Password must be at least 8 characters" };
    }

    // Mock registration logic
    console.log("Server Action: Registration successful", {
        email,
        firstName,
        lastName,
        nickname,
        aboutMe,
        hasAvatar: !!avatar
    });
    return { success: true };
}
