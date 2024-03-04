# retrace

`retrace` takes a [`stdouttrace`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/stdout/stdouttrace)-generated JSON file on stdin and re-exports the traces via [`oltptracehttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp).

This allows you to capture traces offline or in environments where it's difficult to use otel in the intended way.

For example, if you have a file produced by `stdouttrace` at `/tmp/trace.json` and want to view it in [`otel-desktop-viewer`](https://github.com/CtrlSpice/otel-desktop-viewer), you can do this:

```console
$ otel-desktop-viewer &

$ export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
$ export OTEL_TRACES_EXPORTER="otlp"
$ export OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"

$ retrace < /tmp/trace.json
```

## Install

```console
go install github.com/jonjohnsonjr/retrace@latest
```
