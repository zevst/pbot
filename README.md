# Protobuf Span for Jaeger

###### Example
```go
package main

import (
	ot "github.com/opentracing/opentracing-go"
	jcc "github.com/uber/jaeger-client-go/config"
	"github.com/zevst/protobuf_span"
)

func main() {
	cfg := jcc.Configuration{
		// ...
	}
	propagator := protobuf_span.NewPropagator() // You can add BaggagePrefix as an option.

	// Create Jaeger tracer from Configuration
	// be sure to handle the error
	tracer, closer, _ := cfg.NewTracer(
		jcc.Injector(protobuf_span.Format, propagator),
		jcc.Extractor(protobuf_span.Format, propagator),
	)
	ot.SetGlobalTracer(tracer)
	defer func() {
		// be sure to handle the error
		_ = closer.Close()
	}()
	// continue main()
}
```