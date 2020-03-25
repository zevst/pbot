package protobuf_span

import (
	ot "github.com/opentracing/opentracing-go"
	jc "github.com/uber/jaeger-client-go"
	"protobuf-span/pb/span/v1"
)

//ProtoBufSpanFormat is an OpenTracing carrier format constant
const Format = "pb-span-format"

type propagator struct {
	baggagePrefix string
}

// Option sets an option on Propagator
type Option func(p *propagator)

// BaggagePrefix sets baggage prefix on Propagator
func BaggagePrefix(prefix string) Option {
	return func(p *propagator) {
		p.baggagePrefix = prefix
	}
}

func NewPropagator(opts ...Option) *propagator {
	p := new(propagator)
	for _, opt := range opts {
		opt(p)
	}
	return p
}

var emptyContext = jc.SpanContext{}

func (p propagator) Inject(sCtx jc.SpanContext, aCarrier interface{}) (err error) {
	carrier, ok := aCarrier.(*[]byte)
	if !ok {
		return ot.ErrInvalidCarrier
	}
	s := span.Span{
		TraceID:  span.Span_TraceID{Low: sCtx.TraceID().Low, High: sCtx.TraceID().High},
		SpanID:   uint64(sCtx.SpanID()),
		ParentID: uint64(sCtx.ParentID()),
		Flags:    uint32(sCtx.Flags()),
		Baggage:  make(map[string]string),
	}
	sCtx.ForeachBaggageItem(func(k, v string) bool {
		s.Baggage[p.baggagePrefix+k] = v
		return true
	})
	*carrier, err = s.Marshal()
	return err
}

func (p propagator) Extract(aCarrier interface{}) (jc.SpanContext, error) {
	b, ok := aCarrier.(*[]byte)
	if !ok {
		return emptyContext, ot.ErrInvalidCarrier
	}
	var carrier span.Span
	if err := carrier.Unmarshal(*b); err != nil {
		return emptyContext, err
	}
	traceId := jc.TraceID{High: carrier.TraceID.High, Low: carrier.TraceID.Low}
	if !traceId.IsValid() {
		return emptyContext, ot.ErrSpanContextNotFound
	}
	return jc.NewSpanContext(
		traceId,
		jc.SpanID(carrier.SpanID),
		jc.SpanID(carrier.ParentID),
		carrier.Flags != 0,
		carrier.Baggage,
	), nil
}
