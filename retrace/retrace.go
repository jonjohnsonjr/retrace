package retrace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/jonjohnsonjr/leto"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter interface {
	ExportSpans(ctx context.Context, ss []tracesdk.ReadOnlySpan) error
}

func Retrace(ctx context.Context, exp Exporter, in io.Reader) error {
	dec := json.NewDecoder(in)
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
