package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/jonjohnsonjr/retrace/retrace"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	exp, err := texporter.New()
	if err != nil {
		return err
	}
	defer exp.Shutdown(ctx)

	return retrace.Retrace(ctx, exp, os.Stdin)
}
