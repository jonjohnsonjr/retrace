package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/jonjohnsonjr/leto"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// TODO: Probably default to http and have grpc + others be optional.
	client := otlptracehttp.NewClient()

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return err
	}
	defer exp.Shutdown(ctx)

	dec := json.NewDecoder(os.Stdin)

	spans := make([]tracesdk.ReadOnlySpan, 0, 100)

	var span leto.SpanStub
	for {
		if err := dec.Decode(&span); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("decoding span: %w", err)
		}

		spans = append(spans, span.Snapshot())

		if len(spans) == 100 {
			// Flush every 100 spans.
			if err := exp.ExportSpans(ctx, spans); err != nil {
				return fmt.Errorf("exporting spans: %w", err)
			}

			// Reuse the underlying storage to avoid allocs.
			spans = spans[:0]
		}
	}

	if err := exp.ExportSpans(ctx, spans); err != nil {
		return fmt.Errorf("exporting spans: %w", err)
	}

	return nil
}
