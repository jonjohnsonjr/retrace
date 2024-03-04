package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jonjohnsonjr/retrace/retrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

	return retrace.Retrace(ctx, exp, os.Stdin)
}
