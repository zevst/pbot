package pbot

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/uber/jaeger-client-go"
	"github.com/zevst/pbot/pb/span/v1"
)

type SpanProtobufPayload interface {
	proto.Marshaler
	proto.Unmarshaler

	AddPayload(any *types.Any) error
	GetPayload() *types.Any
}

type SpanProtobufSetters interface {
	SetTraceID(id jaeger.TraceID)
	SetSpanID(id jaeger.SpanID)
	SetParentID(id jaeger.SpanID)
	SetFlags(flags byte)
	SetBaggage(key, value string) bool
}

type SpanProtobufGetters interface {
	GetTraceID() jaeger.TraceID
	GetSpanID() jaeger.SpanID
	GetParentID() jaeger.SpanID
	GetFlags() bool
	GetBaggage() map[string]string
}
type pbMsg span.Span

func NewWrapper() SpanProtobufPayload {
	return new(pbMsg)
}

func (s *pbMsg) Marshal() ([]byte, error) {
	msg := (span.Span)(*s)
	defer func() { *s = pbMsg(msg) }()
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg.Marshal()
}

func (s *pbMsg) Unmarshal(b []byte) error {
	var msg span.Span
	defer func() { *s = pbMsg(msg) }()
	if err := msg.Unmarshal(b); err != nil {
		return err
	}
	return msg.Validate()
}

func (s *pbMsg) SetTraceID(id jaeger.TraceID) {
	s.TraceId = span.Span_TraceID{Low: id.Low, High: id.High}
}

func (s *pbMsg) GetTraceID() jaeger.TraceID {
	return jaeger.TraceID{High: s.TraceId.High, Low: s.TraceId.Low}
}

func (s *pbMsg) SetSpanID(id jaeger.SpanID) {
	s.SpanId = uint64(id)
}

func (s *pbMsg) GetSpanID() jaeger.SpanID {
	return jaeger.SpanID(s.SpanId)
}

func (s *pbMsg) SetParentID(id jaeger.SpanID) {
	s.ParentId = uint64(id)
}

func (s *pbMsg) GetParentID() jaeger.SpanID {
	return jaeger.SpanID(s.ParentId)
}

func (s *pbMsg) SetFlags(flags byte) {
	s.Flags = flags
}

func (s *pbMsg) GetFlags() bool {
	return s.Flags != 0
}

func (s *pbMsg) SetBaggage(key, value string) bool {
	if s.Baggage == nil {
		s.Baggage = make(map[string]string)
	}
	s.Baggage[key] = value
	return true
}

func (s *pbMsg) GetBaggage() map[string]string {
	return s.Baggage
}

func (s *pbMsg) AddPayload(any *types.Any) error {
	if s.Payload != nil {
		return errors.New("payload already exists")
	}
	s.Payload = any
	return nil
}

func (s *pbMsg) GetPayload() *types.Any {
	return s.Payload
}
