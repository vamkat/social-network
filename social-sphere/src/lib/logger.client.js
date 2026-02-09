const PREFIX = "SOC";

const LEVELS = { DEBUG: 0, INFO: 1, WARN: 2, ERROR: 3 };

function getMinLevel() {
    const env = (
        (typeof process !== "undefined" && process.env?.NEXT_PUBLIC_LOG_LEVEL) ||
        "INFO"
    ).toUpperCase();
    return LEVELS[env] ?? LEVELS.INFO;
}

function shouldLog(levelName) {
    return LEVELS[levelName] >= getMinLevel();
}

// Same template resolution as server logger
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
            const n = msg.charCodeAt(i + 1) - 48;
            if (n >= 1 && n <= max) {
                resolved += `${pairs[n - 1].key}=${pairs[n - 1].value}`;
                used.add(n - 1);
                i++;
                continue;
            }
        }
        resolved += msg[i];
    }

    const remaining = [];
    for (let i = 0; i < pairs.length; i++) {
        if (!used.has(i)) {
            remaining.push(`${pairs[i].key}:${pairs[i].value}`);
        }
    }

    return { resolved, remaining };
}

const CONSOLE_MAP = {
    DEBUG: console.debug,
    INFO: console.log,
    WARN: console.warn,
    ERROR: console.error,
};

function log(levelName, msg, args) {
    if (!shouldLog(levelName)) return;

    const { resolved, remaining } = resolveTemplate(msg, args);
    let line = `[${PREFIX}]: ${levelName} ${resolved}`;
    if (remaining.length > 0) {
        line += ` args: ${remaining.join(" ")}`;
    }

    CONSOLE_MAP[levelName](line);
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
