import { NextResponse } from 'next/server';
import { validateRegistrationForm } from '@/lib/validation';

export async function POST(request) {
    try {
        const formData = await request.formData();
        console.log("form data below");
        console.log(formData);

        const avatarFile = formData.get('avatar');
        const validation = validateRegistrationForm(formData, avatarFile);

        if (!validation.valid) {
            return NextResponse.json(
                { error: validation.error },
                { status: 400 }
            );
        }

        // register with a public profile
        formData.append('public', 'true');

        const apiBase = process.env.API_BASE || "http://localhost:8081";
        const registerEndpoint = process.env.REGISTER || "/register";

        const backendResponse = await fetch(`${apiBase}${registerEndpoint}`, {
            method: "POST",
            body: formData,
        });

        const responseData = await backendResponse.json().catch(() => null);
        const setCookieHeader = backendResponse.headers.get('set-cookie');

        const response = NextResponse.json(
            responseData || { error: "Registration failed" },
            { status: backendResponse.status }
        );

        if (setCookieHeader) {
            const cookieParts = setCookieHeader.split(';').map(part => part.trim());
            const [nameValue] = cookieParts;
            const [name, value] = nameValue.split('=');

            const attributes = {};
            cookieParts.slice(1).forEach(part => {
                const [key, val] = part.split('=');
                attributes[key.toLowerCase()] = val || true;
            });

            const cookieOptions = {
                path: attributes.path || '/',
                httpOnly: attributes.httponly === true,
                secure: attributes.secure === true,
                domain: 'localhost',
            };

            if (attributes.expires) {
                cookieOptions.expires = new Date(attributes.expires);
            }

            response.cookies.set(name, value, cookieOptions);
        }

        return response;
    } catch (error) {
        console.error("Register API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}
