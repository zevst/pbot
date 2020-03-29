package pbot

import (
	ot "github.com/opentracing/opentracing-go"
	jc "github.com/uber/jaeger-client-go"
	jcc "github.com/uber/jaeger-client-go/config"
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

func Injector(p *propagator) jcc.Option {
	return jcc.Injector(Format, p)
}

func Extractor(p *propagator) jcc.Option {
	return jcc.Extractor(Format, p)
}

func (p propagator) Inject(sCtx jc.SpanContext, aCarrier interface{}) error {
	c, ok := aCarrier.(SpanProtobufSetters)
	if !ok {
		return ot.ErrInvalidCarrier
	}
	c.SetTraceID(sCtx.TraceID())
	c.SetSpanID(sCtx.SpanID())
	c.SetParentID(sCtx.ParentID())
	c.SetFlags(sCtx.Flags())
	sCtx.ForeachBaggageItem(func(k, v string) bool { return c.SetBaggage(p.baggagePrefix+k, v) })
	return nil
}

func (p propagator) Extract(aCarrier interface{}) (jc.SpanContext, error) {
	c, ok := aCarrier.(SpanProtobufGetters)
	if !ok {
		return emptyContext, ot.ErrInvalidCarrier
	}
	if !c.GetTraceID().IsValid() {
		return emptyContext, ot.ErrSpanContextNotFound
	}
	return jc.NewSpanContext(c.GetTraceID(), c.GetSpanID(), c.GetParentID(), c.GetFlags(), c.GetBaggage()), nil
}
