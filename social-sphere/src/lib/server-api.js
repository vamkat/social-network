"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { propagation, context } from "@opentelemetry/api";
import * as logger from "@/lib/logger.server";

const API_BASE = process.env.GATEWAY

export async function serverApiRequest(endpoint, options = {}) {
    const method = options.method || "GET";
    const start = performance.now();

    logger.info("outgoing request @1 @2", "method", method, "url", endpoint);

    try {
        const cookieStore = await cookies();
        const jwt = cookieStore.get("jwt")?.value;

        const headers = { ...(options.headers || {}) };
        if (jwt) headers["Cookie"] = `jwt=${jwt}`;

        // Inject W3C trace context (traceparent/tracestate) for distributed tracing
        propagation.inject(context.active(), headers);

        const res = await fetch(`${API_BASE}${endpoint}`, {
            ...options,
            headers,
            cache: "no-store"
        });


        if (options.forwardCookies) {
            // Handle multiple Set-Cookie headers
            const setCookieHeaders = res.headers.getSetCookie ? res.headers.getSetCookie() : [];

            // Fallback for environments where getSetCookie might not be available
            if (setCookieHeaders.length === 0) {
                const header = res.headers.get('Set-Cookie');
                if (header) setCookieHeaders.push(header);
            }

            if (setCookieHeaders.length > 0) {
                setCookieHeaders.forEach(cookieStr => {

                    const parts = cookieStr.split(';');
                    const [nameValue, ...optionsParts] = parts;
                    const [name, ...valueParts] = nameValue.split('=');
                    const value = valueParts.join('=');

                    if (name && value !== undefined) {
                        const cookieOptions = {
                            secure: false,
                            httpOnly: true,
                            path: '/',
                            sameSite: 'lax',
                        };

                        optionsParts.forEach(part => {
                            const [optKey, optVal] = part.trim().split('=');
                            const keyLower = optKey.toLowerCase();
                            if (keyLower === 'path') cookieOptions.path = optVal;
                            if (keyLower === 'httponly') cookieOptions.httpOnly = true;
                            if (keyLower === 'secure') cookieOptions.secure = true;
                            if (keyLower === 'samesite') cookieOptions.sameSite = optVal.toLowerCase();
                            if (keyLower === 'max-age') cookieOptions.maxAge = parseInt(optVal);
                            if (keyLower === 'expires') cookieOptions.expires = new Date(optVal);
                        });

                        cookieStore.set(name.trim(), value, cookieOptions);
                    }
                });
            }
        }

        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            const duration = Math.round(performance.now() - start);
            const errMsg = err.error || err.message || "Unknown error";

            logger.error("request failed @1 @2 @3 @4 @5",
                "method", method, "url", endpoint, "status", res.status, "error", errMsg, "duration_ms", duration);

            if (res.status === 403) {
                return {ok: false, status: res.status, message: err.error || err.message || "Forbidden"}
            }
            if (res.status === 400) {
                return {ok: false, status: res.status, message: err.error || err.message || "Bad request"}
            }
            if (res.status === 401) {
                cookieStore.delete("jwt");
                redirect("/login");
            }

            return {ok: false, status: res.status, message: errMsg}
        }

        const duration = Math.round(performance.now() - start);

        // Handle empty response bodies (like delete endpoints)
        const text = await res.text();
        if (!text || text.trim() === '') {
            logger.info("request succeeded @1 @2 @3 @4",
                "method", method, "url", endpoint, "status", res.status, "duration_ms", duration);
            return {ok: true, data: null};
        } else {
            logger.info("request succeeded @1 @2 @3 @4",
                "method", method, "url", endpoint, "status", res.status, "duration_ms", duration);
            return {ok: true, data: JSON.parse(text)};
        }

    } catch (e) {
        if (e?.digest?.startsWith("NEXT_REDIRECT")) throw e;
        const duration = Math.round(performance.now() - start);
        logger.error("request exception @1 @2 @3 @4",
            "method", method, "url", endpoint, "error", e?.message || "unknown", "duration_ms", duration);
        return {ok: false, message: "Network error"}
    }
}
