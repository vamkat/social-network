// Barrel file — re-exports the correct logger for the current environment.
// Server-side ("use server" files, server actions, instrumentation) get OTEL + stdout.
// Client-side (browser components) get console-only output.
//
// In practice, server actions should import "@/lib/logger.server" directly
// since they always run server-side. This barrel is for shared code
// that may run in either environment.

export { debug, info, warn, error } from "./logger.client.js";
