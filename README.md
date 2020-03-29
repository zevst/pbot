# Protobuf Span for Jaeger

###### Example. Add injector/extractor format
```go
package main

import (
	ot "github.com/opentracing/opentracing-go"
	jcc "github.com/uber/jaeger-client-go/config"
	"github.com/zevst/pbot"
)

func main() {
	cfg := jcc.Configuration{
		// ...
	}
	propagator := pbot.NewPropagator() // You can add BaggagePrefix as an option.

	// Create Jaeger tracer from Configuration
	// be sure to handle the error
	tracer, closer, _ := cfg.NewTracer(pbot.Injector(propagator), pbot.Extractor(propagator))
	ot.SetGlobalTracer(tracer)
	defer func() {
		// be sure to handle the error
		_ = closer.Close()
	}()
	// continue main()
}
```

###### Example. Inject
```go
// do not forget to handle errors.
span := opentracing.SpanFromContext(ctx)
if span == nil {
    span = opentracing.GlobalTracer().StartSpan("Operation", ext.SpanKindProducer)
}
ctx = opentracing.ContextWithSpan(ctx, span)
wrapper := pbot.NewWrapper()
_ = opentracing.GlobalTracer().Inject(span.Context(), pbot.Format, wrapper)
_ = wrapper.AddPayload(req.GetData())
value, _ := wrapper.Marshal()
// send value
```

###### Example. Extract
```go
// do not forget to handle errors.
wrapper := pbot.NewWrapper()
_ = wrapper.Unmarshal(msg.Data)
sc, _ := tracer.Extract(pbot.Format, wrapper)
opts := []opentracing.StartSpanOption{
    ext.SpanKindConsumer,
    opentracing.ChildOf(sc),
}
span := tracer.StartSpan("Operation", opts...)
defer span.Finish()
ctx = opentracing.ContextWithSpan(ctx, span)
// read wrapper.GetPayload()
```