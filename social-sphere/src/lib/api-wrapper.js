/**
 * Wrapper for making safe API calls with standardized error handling.
 * Includes cookies by default.
 * @param {string} url - The API endpoint URL
 * @param {Object} options - Fetch options (method, headers, body, etc.)
 * @returns {Promise<{success: boolean, data?: any, error?: string}>}
 */
export async function safeApiCall(url, options = {}) {
    console.log(url)
    try {
        const response = await fetch(url, {
            credentials: "include", // Default to include cookies
            ...options,
        });
       
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error(`API Error (${url}):`, errorData);
            return {
                success: false,
                error: errorData.error || "Request failed. Please try again."
            };
        }

        const data = await response.json();
        return { success: true, data };
    } catch (error) {
        console.error(`Network Error (${url}):`, error);
        return {
            success: false,
            error: "Network error. Please try again later."
        };
    }
}
