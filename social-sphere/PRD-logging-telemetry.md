# PRD: Add Logging & Telemetry to Social-Sphere

## 1. Overview

Social-sphere (the Next.js frontend server) currently has zero observability. It uses 217 raw `console.log/error/warn` calls scattered across 90 files, with no structured logging, no centralized collection, and no trace correlation with the backend.

The backend already has a production-ready OpenTelemetry pipeline: structured logs and traces are exported via OTLP gRPC to Grafana Alloy, which forwards them to VictoriaLogs and VictoriaTraces, visualized in Grafana. This PRD describes how to plug social-sphere into that same pipeline so the entire application has unified observability.

---

## 2. Goals

1. **Structured server-side logging** from social-sphere flows into the same Grafana/VictoriaLogs instance as the backend
2. **Same log format** as the backend — `HH:mm:ss.SSS [SOC]: LEVEL message key=value` on stdout, OTLP export to Alloy
3. **Same message template syntax** — `@1`, `@2` placeholders with key-value argument pairs
4. **Trace context propagation** — outgoing fetch calls from `serverApiRequest` inject `traceparent`/`tracestate` headers so the backend gateway can continue the trace, enabling end-to-end distributed tracing
5. **Minimal initial scope** — instrument the `serverApiRequest` chokepoint first, not every component
6. **Client-side log formatting** — consistent console output in the browser with log-level filtering (no OTEL export from browser)

---

## 3. Non-Goals (Out of Scope for V1)

- Browser-side telemetry export (e.g. sending logs from the browser to Alloy)
- Metrics collection (backend has this as a TODO too — `NewMeterProvider` is commented out)
- Replacing all 217 console calls in one pass (will be done incrementally after the foundation is in place)
- Custom Grafana dashboards (logs will be queryable in the existing VictoriaLogs datasource immediately)

---

## 4. Current State

### 4.1 Backend Telemetry Stack (reference implementation)

| Component | Details |
|-----------|---------|
| **Telemetry package** | `backend/shared/go/telemetry/` — `telemetry.go`, `logger.go`, `otel.go`, `tracer.go` |
| **Public API** | `tele.Debug()`, `tele.Info()`, `tele.Warn()`, `tele.Error()`, `tele.Fatal()`, `tele.Trace()` |
| **Logger bridge** | `otelslog` — bridges Go's `slog` to the OTEL LoggerProvider |
| **Exporters** | `otlptracegrpc` + `otlploggrpc` → Alloy at `alloy:4317` (gRPC, insecure) |
| **Propagator** | `TraceContext` + `Baggage` |
| **Resource attributes** | `service.name`, `service.namespace=social-network`, `deployment.environment=dev` |
| **Log attributes** | `customArgs` (key-value pairs), `context` (user-id, request-id, ip), `callers` (function+line), `prefix` (3-letter code) |
| **Stdout format** | `HH:mm:ss.SSS [PREFIX]: LEVEL message args` |
| **Message templates** | `@1`, `@2`...`@9` replaced with `key=value` from positional arg pairs |
| **Service prefixes** | `USR`, `API`, `MED`, `POS`, `CHA`, `NOT`, `LIV` |
| **Context keys** | `user-id`, `request-id`, `ip` (defined in `backend/shared/go/ct/ctxKey.go`) |
| **Init pattern** | Each service calls `tele.InitTelemetry(ctx, "name", "PREFIX", "alloy:4317", keys, debug, simplePrint)` at startup |

### 4.2 Infrastructure (already running)

| Service | Address | Role |
|---------|---------|------|
| **Grafana Alloy** | `alloy:4317` (gRPC), `alloy:4318` (HTTP) | OTLP receiver → batch → export |
| **VictoriaLogs** | `victoria-logs:9428` | Log storage, queryable from Grafana |
| **VictoriaTraces** | `victoria-traces:10428` | Trace storage (Jaeger-compatible) |
| **Grafana** | `localhost:3001` | Visualization with VictoriaLogs + VictoriaTraces datasources |

### 4.3 Social-Sphere Current State

| Aspect | Status |
|--------|--------|
| **Framework** | Next.js 16.0.10, React 19.2.1, JavaScript (not TypeScript) |
| **Output mode** | Standalone (`output: 'standalone'` in `next.config.mjs`) |
| **Docker** | `node:20-slim`, runs `node server.js` |
| **Logging library** | None |
| **Instrumentation file** | None |
| **Middleware** | None |
| **API routes** | None — all server-side logic is in server actions (`src/actions/`) |
| **Server actions** | 62 files across 10 directories (auth, chat, events, groups, notifs, posts, profile, requests, search, users) |
| **Central fetch utility** | `src/lib/server-api.js` — `serverApiRequest(endpoint, options)` — every server action calls this |
| **Console usage** | 217 statements across 90 files |
| **LOG_LEVEL env vars** | Defined in `.env.*.example` but completely unused in code |
| **OTEL env vars** | Not present in `docker-compose.yml` for social-sphere |

---

## 5. Architecture

### 5.1 How Telemetry Flows

```
Browser (user clicks)
    ↓
Next.js Server Action (e.g. createPost)
    ↓
serverApiRequest("POST", "/posts")     ← LOG here: outgoing request
    ↓                                   ← INJECT traceparent header
    ↓ fetch()
    ↓
API Gateway (backend)                   ← picks up traceparent, continues trace
    ↓
Posts Service (backend)                 ← backend already logs its side
    ↓
Response returns
    ↓
serverApiRequest receives response      ← LOG here: response status + duration
    ↓
Server Action returns to browser
```

**Log pipeline:**
```
social-sphere logger
    ↓ OTLP gRPC
Grafana Alloy (alloy:4317)             ← same collector as backend
    ↓ batch processor
VictoriaLogs (victoria-logs:9428)       ← same storage as backend
    ↓
Grafana (localhost:3001)                ← query with service.name=social-sphere
```

### 5.2 What Gets Logged

Instrument `serverApiRequest` in `src/lib/server-api.js` — the single chokepoint. Every server action goes through it.

| Event | Level | Attributes | Example Output |
|-------|-------|------------|----------------|
| Outgoing request | `info` | method, url | `14:32:05.123 [SOC]: INFO outgoing request method=POST url=/posts` |
| Successful response | `info` | method, url, status, duration_ms | `14:32:05.207 [SOC]: INFO request succeeded method=POST url=/posts status=201 duration_ms=84` |
| Failed response (4xx/5xx) | `error` | method, url, status, error, duration_ms | `14:32:05.207 [SOC]: ERROR request failed method=POST url=/posts status=500 error=internal duration_ms=84` |
| Fetch exception (network error) | `error` | method, url, error | `14:32:05.207 [SOC]: ERROR request exception method=POST url=/posts error=ECONNREFUSED` |

That's it for V1. Four log lines covering the most operationally important data.

### 5.3 What Does NOT Get Logged

- Button clicks, UI interactions, re-renders (use browser devtools)
- State changes in Zustand stores
- Individual component lifecycle events
- Request/response bodies (too large, potential PII)

---

## 6. Implementation Plan

### 6.1 New Files to Create

#### 6.1.1 `src/instrumentation.js`

Next.js has a built-in hook: if a file named `instrumentation.js` exists at the project root (or `src/`), Next.js calls its `register()` export once when the server starts. This is where we initialize the OTEL SDK.

**Responsibilities:**
- Initialize `NodeSDK` from `@opentelemetry/sdk-node`
- Configure `OTLPTraceExporter` targeting `process.env.OTEL_EXPORTER_OTLP_ENDPOINT || "alloy:4317"` (gRPC, insecure)
- Configure `OTLPLogExporter` targeting the same address
- Set resource attributes: `service.name=social-sphere`, `service.namespace=social-network`, `deployment.environment=dev`
- Set propagator: `CompositePropagator` with `W3CTraceContextPropagator` + `W3CBaggagePropagator` (matches backend's `TraceContext{}` + `Baggage{}`)
- Register `HttpInstrumentation` for auto-instrumented outgoing fetch calls
- Only runs server-side (Next.js guarantees this)

**Reference:** Backend equivalent is `SetupOTelSDK()` in `backend/shared/go/telemetry/otel.go:29-84`

#### 6.1.2 `src/lib/logger.server.js`

Server-side structured logger that mirrors the backend's `tele` package.

**Exported functions:**
```js
export function debug(message, ...args) { }
export function info(message, ...args)  { }
export function warn(message, ...args)  { }
export function error(message, ...args) { }
```

**Argument format** (matches backend `logger.go`):
```js
// key-value pairs, same as backend
info("outgoing request @1 @2", "method", "POST", "url", "/posts")
// → resolves @1 to method=POST, @2 to url=/posts
```

**Each call does two things** (same as backend's `logger.go:77-156`):

1. **Emit to OTEL** via the `LoggerProvider` set up in `instrumentation.js`:
   - Log body = resolved message string
   - Attributes:
     - `customArgs` group — the key-value pairs
     - `context` group — `user-id`, `request-id`, `ip` (if available)
     - `callers` — function name + line number (from `new Error().stack`)
     - `prefix` — `"SOC"`
   - This mirrors the backend's `slog.Log()` call with `slog.GroupAttrs("customArgs", ...)`, `slog.GroupAttrs("context", ...)`, `slog.String("callers", ...)`, `slog.String("prefix", ...)`

2. **Print to stdout** in the same format as the backend:
   ```
   HH:mm:ss.SSS [SOC]: LEVEL message key=value key=value
   ```
   This mirrors the backend's stdout formatting in `logger.go:143-155`

**Log level filtering:**
- Read from `process.env.LOG_LEVEL` (finally using the already-defined env var)
- `DEBUG` only shown when `process.env.ENABLE_DEBUG_LOGS === "true"` (same as backend's `enableDebug` flag)
- Levels: `DEBUG < INFO < WARN < ERROR` (same as Go's `slog.Level`)

**Message template resolution** (matches backend `logger.go:91-118`):
- Scan message for `@N` where N is 1-9
- Replace with `key=value` from the corresponding argument pair
- Remaining args appended as ` args: key:value key:value`

#### 6.1.3 `src/lib/logger.client.js`

Client-side lightweight logger. No OTEL export.

**Exported functions:** Same API as `logger.server.js`

**Behavior:**
- Routes to `console.debug/log/warn/error`
- Formats with same `[SOC]: LEVEL message` prefix for recognizable output
- Respects `NEXT_PUBLIC_LOG_LEVEL` env var
- Resolves `@1`, `@2` templates identically

#### 6.1.4 `src/lib/logger.js`

Barrel file that re-exports the correct logger depending on environment:

```js
if (typeof window === "undefined") {
  // server
  module.exports = require("./logger.server.js");
} else {
  // client
  module.exports = require("./logger.client.js");
}
```

Or alternatively, use separate imports where needed (`"use server"` files import server logger directly).

---

### 6.2 Files to Modify

#### 6.2.1 `src/lib/server-api.js`

This is the primary instrumentation target. Current state has 4 `console.log` calls.

**Changes:**
1. Import the server logger
2. At the start of `serverApiRequest`: record start time, log outgoing request
3. On successful response: log status + duration
4. On error response (4xx/5xx): log status + error message + duration
5. On catch (network exception): log error
6. **Trace context injection**: use OTEL propagation API to inject `traceparent` and `tracestate` into the outgoing fetch headers. This is what enables end-to-end distributed tracing with the backend.

**Before:**
```js
console.log("ERROR: ", err);
console.log("STATUS: ", res.status);
console.log("Data: ", JSON.parse(text));
```

**After:**
```js
import * as logger from "@/lib/logger.server";
import { propagation, context, trace } from "@opentelemetry/api";

export async function serverApiRequest(endpoint, options = {}) {
    const method = options.method || "GET";
    const start = performance.now();

    logger.info("outgoing request @1 @2", "method", method, "url", endpoint);

    // Inject trace context into headers for distributed tracing
    const headers = { ...(options.headers || {}) };
    propagation.inject(context.active(), headers);

    // ... existing fetch logic ...

    if (!res.ok) {
        const duration = Math.round(performance.now() - start);
        logger.error("request failed @1 @2 @3 @4",
            "method", method, "url", endpoint, "status", res.status, "duration_ms", duration);
        // ... existing error handling ...
    }

    const duration = Math.round(performance.now() - start);
    logger.info("request succeeded @1 @2 @3",
        "method", method, "url", endpoint, "status", res.status, "duration_ms", duration);
}
```

#### 6.2.2 `docker-compose.yml`

Add OTEL environment variables to the `social-sphere` service block:

```yaml
social-sphere:
  environment:
    - NODE_ENV=production
    - GATEWAY=http://api-gateway:8081
    - LIVE=ws://localhost:8082
    # New: OpenTelemetry configuration (matches backend services)
    - OTEL_EXPORTER_OTLP_INSECURE=true
    - OTEL_EXPORTER_OTLP_ENDPOINT=alloy:4317
    - OTEL_RESOURCE_ATTRIBUTES=service.name=social-sphere,service.namespace=social-network,deployment.environment=dev
    - ENABLE_DEBUG_LOGS=true
    - LOG_LEVEL=DEBUG
```

No changes needed to:
- **Alloy config** (`config.alloy`) — it already accepts OTLP from any service on the `social-network` network
- **VictoriaLogs** — it stores whatever Alloy sends
- **Grafana datasources** — VictoriaLogs datasource already exists

#### 6.2.3 `next.config.mjs`

Add the `serverExternalPackages` configuration to ensure gRPC native modules work correctly in the standalone build:

```js
const nextConfig = {
  // ... existing config ...
  serverExternalPackages: ["@grpc/grpc-js"],
};
```

#### 6.2.4 `.env.development.example` and `.env.production.example`

Add new OTEL variables:

```
# OpenTelemetry
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
ENABLE_DEBUG_LOGS=true
```

#### 6.2.5 `Dockerfile`

No changes needed. The standalone output already includes `node_modules` dependencies. The `node server.js` entrypoint works because Next.js auto-loads `instrumentation.js` during server startup.

---

### 6.3 NPM Dependencies to Add

```json
{
  "@opentelemetry/api": "^1.9.0",
  "@opentelemetry/sdk-node": "^0.57.0",
  "@opentelemetry/sdk-trace-node": "^1.30.0",
  "@opentelemetry/sdk-logs": "^0.57.0",
  "@opentelemetry/api-logs": "^0.57.0",
  "@opentelemetry/exporter-trace-otlp-grpc": "^0.57.0",
  "@opentelemetry/exporter-logs-otlp-grpc": "^0.57.0",
  "@opentelemetry/resources": "^1.30.0",
  "@opentelemetry/semantic-conventions": "^1.28.0",
  "@opentelemetry/instrumentation-http": "^0.57.0",
  "@grpc/grpc-js": "^1.12.0"
}
```

---

## 7. Log Attribute Schema

To ensure social-sphere logs are queryable alongside backend logs in Grafana, we must use the same attribute structure.

### 7.1 OTEL Log Record Attributes

| Attribute | Type | Source | Backend Equivalent |
|-----------|------|--------|--------------------|
| `prefix` | string | Hardcoded `"SOC"` | `slog.String("prefix", l.prefix)` in `logger.go:136` |
| `callers` | string | Parsed from `new Error().stack` | `slog.String("callers", callerInfo)` in `logger.go:135` |
| `customArgs.*` | group | Key-value pairs from function args | `slog.GroupAttrs("customArgs", kvPairsToAttrs(args)...)` in `logger.go:133` |
| `context.user-id` | string | Extracted from JWT cookie | `slog.GroupAttrs("context", ctxArgs...)` in `logger.go:134` |
| `context.request-id` | string | Generated UUID per request | Same |
| `context.ip` | string | From request headers | Same |

### 7.2 OTEL Resource Attributes

Set via `OTEL_RESOURCE_ATTRIBUTES` env var (same pattern as all backend services):

```
service.name=social-sphere
service.namespace=social-network
deployment.environment=dev
```

### 7.3 Trace Context Headers

Injected into outgoing `fetch()` calls by the OTEL propagation API:

| Header | Purpose | Backend Equivalent |
|--------|---------|-------------------|
| `traceparent` | W3C Trace Context trace ID + span ID | Set by `otelhttp` middleware in backend |
| `tracestate` | Vendor-specific trace data | Set by `propagation.Baggage{}` in backend |

The backend gateway already reads these headers via `otelhttp.NewHandler()` in `backend/shared/go/http-middleware/`. No backend changes needed.

---

## 8. Implementation Order

| Step | Task | Files | Effort |
|------|------|-------|--------|
| 1 | Install npm dependencies | `package.json` | Small |
| 2 | Create `instrumentation.js` | `src/instrumentation.js` | Small |
| 3 | Create server logger | `src/lib/logger.server.js` | Medium |
| 4 | Create client logger | `src/lib/logger.client.js` | Small |
| 5 | Create logger barrel | `src/lib/logger.js` | Small |
| 6 | Instrument `serverApiRequest` | `src/lib/server-api.js` | Small |
| 7 | Update docker-compose | `docker-compose.yml` | Small |
| 8 | Update next.config.mjs | `next.config.mjs` | Small |
| 9 | Update env examples | `.env.*.example` | Small |
| 10 | Verify in Grafana | — | Testing |

**After V1 is verified working:**

| Step | Task | Files | Effort |
|------|------|-------|--------|
| 11 | Replace `console.*` in server actions (62 files) | `src/actions/**/*.js` | Medium-Large |
| 12 | Replace `console.*` in client components (~55 files) | `src/components/**/*.js` | Medium-Large |
| 13 | Replace `console.*` in contexts (3 files) | `src/context/*.js` | Small |
| 14 | Update K8s deployment | `backend/k8s/` | Small |

---

## 9. Verification Criteria

### 9.1 Functional

- [ ] `docker compose up` starts social-sphere with OTEL env vars
- [ ] `instrumentation.js` `register()` runs on server startup (visible in stdout)
- [ ] Every `serverApiRequest` call produces an `info` log line on stdout in the format `HH:mm:ss.SSS [SOC]: INFO outgoing request method=X url=Y`
- [ ] Every response produces a log line with status and duration
- [ ] Error responses produce `ERROR` level logs
- [ ] Logs appear in Grafana → Explore → VictoriaLogs datasource when filtering `service.name=social-sphere`
- [ ] `LOG_LEVEL` env var controls which levels are emitted
- [ ] `ENABLE_DEBUG_LOGS=false` suppresses debug-level logs

### 9.2 Distributed Tracing

- [ ] Outgoing fetch calls from `serverApiRequest` include `traceparent` header
- [ ] In Grafana → Explore → VictoriaTraces, a trace starting from social-sphere shows child spans from the backend gateway and downstream services
- [ ] Trace IDs correlate between social-sphere logs and backend logs

### 9.3 No Regressions

- [ ] All existing server actions continue to work (auth, posts, chat, etc.)
- [ ] Docker image builds successfully with new dependencies
- [ ] Standalone output includes OTEL dependencies
- [ ] No client-side bundle size increase (OTEL packages are server-only)

---

## 10. Reference: Backend vs Social-Sphere Mapping

| Backend (Go) | Social-Sphere (JS) | Notes |
|--------------|---------------------|-------|
| `tele.InitTelemetry()` | `register()` in `instrumentation.js` | Next.js calls this automatically |
| `SetupOTelSDK()` in `otel.go` | `NodeSDK` init in `instrumentation.js` | Same exporters, same address |
| `NewLoggerProvider()` in `otel.go` | `OTLPLogExporter` + `LoggerProvider` | Same gRPC endpoint |
| `NewTracerProvider()` in `otel.go` | `OTLPTraceExporter` + `NodeTracerProvider` | Same gRPC endpoint |
| `NewPropagator()` in `otel.go` | `W3CTraceContextPropagator` + `W3CBaggagePropagator` | Same propagation standards |
| `tele.Info(ctx, msg, args...)` | `logger.info(msg, ...args)` | Same API shape |
| `@1`, `@2` template syntax | `@1`, `@2` template syntax | Identical |
| `logging.log()` in `logger.go` | `log()` in `logger.server.js` | Both emit to OTEL + stdout |
| `functionCallers()` in `logger.go` | `getCallers()` via `new Error().stack` | JS equivalent of `runtime.Callers` |
| `context2Attributes()` in `logger.go` | Extract from cookie/headers | JS doesn't have Go's context, use cookies |
| 3-letter prefix (`USR`, `API`) | `SOC` | One new prefix |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Same env var | Same value: `alloy:4317` |
| `ENABLE_DEBUG_LOGS` | Same env var | Same behavior |
