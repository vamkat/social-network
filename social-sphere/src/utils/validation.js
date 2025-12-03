/**
 * Shared validation utilities for client-side and server-side validation
 */

// Email validation pattern
export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Password strength pattern: at least 1 lowercase, 1 uppercase, 1 number, 1 symbol
export const STRONG_PASSWORD_PATTERN = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$/;

// Username/nickname pattern: letters, numbers, dots, underscores, dashes
export const USERNAME_PATTERN = /^[A-Za-z0-9_.-]+$/;

// File validation constants
export const MAX_FILE_SIZE = 20 * 1024 * 1024; // 20MB
export const ALLOWED_FILE_TYPES = ["image/jpeg", "image/png", "image/gif"];

/**
 * Calculate age from date of birth
 * @param {string} dateOfBirth - Date string in YYYY-MM-DD format
 * @returns {number} Age in years
 */
export function calculateAge(dateOfBirth) {
    const today = new Date();
    const birthDate = new Date(dateOfBirth);
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDiff = today.getMonth() - birthDate.getMonth();

    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }

    return age;
}

/**
 * Validate email format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid
 */
export function isValidEmail(email) {
    return EMAIL_PATTERN.test(email);
}

/**
 * Validate password strength
 * @param {string} password - Password to validate
 * @returns {boolean} True if meets strength requirements
 */
export function isStrongPassword(password) {
    return password.length >= 8 && STRONG_PASSWORD_PATTERN.test(password);
}

/**
 * Validate username/nickname format
 * @param {string} username - Username to validate
 * @returns {boolean} True if valid
 */
export function isValidUsername(username) {
    return username.length >= 4 && USERNAME_PATTERN.test(username);
}

/**
 * Validate avatar file
 * @param {File} file - File object to validate
 * @returns {{valid: boolean, error: string}} Validation result
 */
export function isValidAvatarFile(file) {
    if (!file) return { valid: true, error: "" }; // Optional

    if (!ALLOWED_FILE_TYPES.includes(file.type)) {
        return { valid: false, error: "Avatar must be JPEG, PNG, or GIF." };
    }

    if (file.size > MAX_FILE_SIZE) {
        return { valid: false, error: "Avatar must be less than 5MB." };
    }

    return { valid: true, error: "" };
}

/**
 * Validate registration form data (client-side)
 * @param {FormData} formData - Form data to validate
 * @param {File|null} avatarFile - Avatar file object
 * @returns {{valid: boolean, error: string}} Validation result
 */
export function validateRegistrationForm(formData, avatarFile = null) {
    // First name validation
    const firstName = formData.get("firstName")?.trim() || "";
    if (!firstName) {
        return { valid: false, error: "First name is required." };
    }
    if (firstName.length < 2) {
        return { valid: false, error: "First name must be at least 2 characters." };
    }

    // Last name validation
    const lastName = formData.get("lastName")?.trim() || "";
    if (!lastName) {
        return { valid: false, error: "Last name is required." };
    }
    if (lastName.length < 2) {
        return { valid: false, error: "Last name must be at least 2 characters." };
    }

    // Email validation
    const email = formData.get("email")?.trim() || "";
    if (!isValidEmail(email)) {
        return { valid: false, error: "Please enter a valid email address." };
    }

    // Password validation
    const password = formData.get("password");
    const confirmPassword = formData.get("confirmPassword");
    if (!password || !confirmPassword) {
        return { valid: false, error: "Please enter both password and confirm password." };
    }
    if (password.length < 8) {
        return { valid: false, error: "Password must be at least 8 characters." };
    }
    if (!STRONG_PASSWORD_PATTERN.test(password)) {
        return { valid: false, error: "Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol." };
    }
    if (password !== confirmPassword) {
        return { valid: false, error: "Passwords do not match" };
    }

    // Date of birth validation
    const dateOfBirth = formData.get("dateOfBirth")?.trim() || "";
    if (!dateOfBirth) {
        return { valid: false, error: "Date of birth is required." };
    }
    const age = calculateAge(dateOfBirth);
    if (age < 13 || age > 111) {
        return { valid: false, error: "You must be between 13 and 111 years old." };
    }

    // Nickname validation (optional)
    const username = formData.get("nickname")?.trim() || "";
    if (username) {
        if (username.length < 4) {
            return { valid: false, error: "Username must be at least 4 characters." };
        }
        if (!USERNAME_PATTERN.test(username)) {
            return { valid: false, error: "Username can only use letters, numbers, dots, underscores, or dashes." };
        }
    }

    // About me validation (optional)
    const aboutMe = formData.get("aboutMe")?.trim() || "";
    if (aboutMe && aboutMe.length > 400) {
        return { valid: false, error: "About me must be at most 400 characters." };
    }

    // Avatar validation (optional)
    if (avatarFile) {
        const avatarValidation = isValidAvatarFile(avatarFile);
        if (!avatarValidation.valid) {
            return avatarValidation;
        }
    }

    return { valid: true, error: "" };
}

/**
 * Validate login form data (client-side)
 * @param {FormData} formData - Form data to validate
 * @returns {{valid: boolean, error: string}} Validation result
 */
export function validateLoginForm(formData) {
    // Identifier validation
    const identifier = formData.get("identifier")?.trim() || "";
    if (!identifier) {
        return { valid: false, error: "Email or username is required." };
    }

    // Password validation
    const password = formData.get("password");
    if (!password) {
        return { valid: false, error: "Password is required." };
    }

    return { valid: true, error: "" };
}
