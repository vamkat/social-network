import { NodeSDK } from "@opentelemetry/sdk-node";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { OTLPLogExporter } from "@opentelemetry/exporter-logs-otlp-grpc";
import { SimpleLogRecordProcessor } from "@opentelemetry/sdk-logs";
import { Resource } from "@opentelemetry/resources";
import {
    W3CTraceContextPropagator,
    W3CBaggagePropagator,
    CompositePropagator,
} from "@opentelemetry/core";
import { credentials } from "@grpc/grpc-js";

const endpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT || "localhost:4317";
const insecure = process.env.OTEL_EXPORTER_OTLP_INSECURE === "true";
const environment = process.env.DEPLOYMENT_ENVIRONMENT || "dev";

const grpcCredentials = insecure
    ? credentials.createInsecure()
    : credentials.createSsl();

const resource = new Resource({
    "service.name": "social-sphere",
    "service.namespace": "social-network",
    "deployment.environment": environment,
});

const traceExporter = new OTLPTraceExporter({
    url: `http://${endpoint}`,
    credentials: grpcCredentials,
});

const logExporter = new OTLPLogExporter({
    url: `http://${endpoint}`,
    credentials: grpcCredentials,
});

const sdk = new NodeSDK({
    resource,
    traceExporter,
    logRecordProcessors: [new SimpleLogRecordProcessor(logExporter)],
    textMapPropagator: new CompositePropagator({
        propagators: [
            new W3CTraceContextPropagator(),
            new W3CBaggagePropagator(),
        ],
    }),
});

sdk.start();

const time = new Date().toLocaleTimeString("en-GB", {
    hour12: false,
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    fractionalSecondDigits: 3,
});
process.stdout.write(
    `${time} [SOC]: INFO OpenTelemetry initialized endpoint=${endpoint} insecure=${insecure}\n`
);

process.on("SIGTERM", () => {
    sdk.shutdown().then(
        () => process.exit(0),
        () => process.exit(1)
    );
});
