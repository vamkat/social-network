import { logs, SeverityNumber } from "@opentelemetry/api-logs";

const PREFIX = "SOC";

const LEVELS = { DEBUG: 0, INFO: 1, WARN: 2, ERROR: 3 };

const SEVERITY_MAP = {
    DEBUG: SeverityNumber.DEBUG,
    INFO: SeverityNumber.INFO,
    WARN: SeverityNumber.WARN,
    ERROR: SeverityNumber.ERROR,
};

function getMinLevel() {
    const env = (process.env.LOG_LEVEL || "INFO").toUpperCase();
    return LEVELS[env] ?? LEVELS.INFO;
}

function shouldLog(levelName) {
    if (levelName === "DEBUG" && process.env.ENABLE_DEBUG_LOGS !== "true") {
        return false;
    }
    return LEVELS[levelName] >= getMinLevel();
}

// Resolve @1..@9 templates and build key-value pairs from args
// Mirrors backend logger.go:91-118
function resolveTemplate(msg, args) {
    const pairs = [];
    for (let i = 0; i < args.length - 1; i += 2) {
        pairs.push({ key: String(args[i]), value: String(args[i + 1]) });
    }

    const used = new Set();
    const max = Math.min(9, pairs.length);

    let resolved = "";
    for (let i = 0; i < msg.length; i++) {
        if (msg[i] === "@" && i + 1 < msg.length) {
            const n = msg.charCodeAt(i + 1) - 48; // '0' = 48
            if (n >= 1 && n <= max) {
                resolved += `${pairs[n - 1].key}=${pairs[n - 1].value}`;
                used.add(n - 1);
                i++; // skip the digit
                continue;
            }
        }
        resolved += msg[i];
    }

    // Collect remaining (not referenced by @N) args
    const remaining = [];
    for (let i = 0; i < pairs.length; i++) {
        if (!used.has(i)) {
            remaining.push(`${pairs[i].key}:${pairs[i].value}`);
        }
    }

    return { resolved, pairs, remaining };
}

// Get caller info from stack trace (JS equivalent of runtime.Callers in Go)
function getCallers() {
    const err = new Error();
    const lines = (err.stack || "").split("\n");
    // Skip: Error, getCallers, log, debug/info/warn/error → start at index 4
    const callerLines = lines.slice(4, 7);
    return callerLines
        .map((line) => {
            const match = line.match(/at\s+(.+?)\s+\((.+):(\d+):\d+\)/);
            if (match) return `by ${match[1]} at ${match[3]}`;
            const match2 = line.match(/at\s+(.+):(\d+):\d+/);
            if (match2) return `by ${match2[1]} at ${match2[2]}`;
            return line.trim();
        })
        .join("\n");
}

// Format timestamp matching backend: HH:mm:ss.SSS
function timestamp() {
    const now = new Date();
    const h = String(now.getHours()).padStart(2, "0");
    const m = String(now.getMinutes()).padStart(2, "0");
    const s = String(now.getSeconds()).padStart(2, "0");
    const ms = String(now.getMilliseconds()).padStart(3, "0");
    return `${h}:${m}:${s}.${ms}`;
}

function log(levelName, msg, args) {
    if (!shouldLog(levelName)) return;

    const { resolved, pairs, remaining } = resolveTemplate(msg, args);
    const callerInfo = getCallers();

    // 1. Emit OTEL log record (mirrors logger.go:129-136)
    try {
        const logger = logs.getLogger("social-sphere");
        const customArgs = {};
        for (const { key, value } of pairs) {
            customArgs[key] = value;
        }

        logger.emit({
            severityNumber: SEVERITY_MAP[levelName],
            severityText: levelName,
            body: resolved,
            attributes: {
                prefix: PREFIX,
                callers: callerInfo,
                ...Object.fromEntries(
                    Object.entries(customArgs).map(([k, v]) => [`customArgs.${k}`, v])
                ),
            },
        });
    } catch {
        // OTEL not initialized yet (e.g. during startup) — continue to stdout
    }

    // 2. Print to stdout (mirrors logger.go:143-155)
    let line = `${timestamp()} [${PREFIX}]: ${levelName} ${resolved}`;
    if (remaining.length > 0) {
        line += ` args: ${remaining.join(" ")}`;
    }
    process.stdout.write(line + "\n");
}

export function debug(msg, ...args) {
    log("DEBUG", msg, args);
}

export function info(msg, ...args) {
    log("INFO", msg, args);
}

export function warn(msg, ...args) {
    log("WARN", msg, args);
}

export function error(msg, ...args) {
    log("ERROR", msg, args);
}
