package context

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func init() {
	keyRequestID = "id"
	if value := os.Getenv("CONTEXT_KEY_REQUEST_ID"); len(value) > 0 {
		keyRequestID = value
	}
}

type Context interface {
	context.Context

	WithDeadline(d time.Time)
	CopyWithDeadline(d time.Time) Context

	WithTimeout(timeout time.Duration)
	CopyWithTimeout(timeout time.Duration) Context
	Cancel()

	Copy() Context

	Value
}

var (
	keyRequestID string
)

type local struct {
	base       context.Context
	cancelFunc context.CancelFunc
}

func (l local) Deadline() (deadline time.Time, ok bool) {
	return l.base.Deadline()
}

func (l local) Done() <-chan struct{} {
	return l.base.Done()
}

func (l local) Err() error {
	return l.base.Err()
}

func (l *local) Cancel() {
	if l.cancelFunc != nil {
		l.cancelFunc()
	}
}

func (l local) Copy() Context {
	return &l
}

func (l *local) isEmptyID() bool {
	_, ok := l.id()
	return !ok
}

var cancelFunc = func() {}

func Empty() Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
	}

	ctx.WithValue(keyRequestID, uuid.New())

	return ctx
}

func New(option interface{}) Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
	}

	switch baseCtx := option.(type) {
	case Context:
		ctx.withValue(keyRequestID, baseCtx.ID())
	case context.Context:
		ctx.base = baseCtx
	}

	if ctx.isEmptyID() {
		ctx.withValue(keyRequestID, uuid.New().String())
	}
	return ctx
}

func NewWithSpanContext(option interface{}, spanContext opentracing.SpanContext) Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
	}

	switch baseCtx := option.(type) {
	case Context:
		ctx.withValue(keyRequestID, baseCtx.ID())
	case context.Context:
		ctx.base = baseCtx
	}

	ctx.SetTraceID(spanContext.(jaeger.SpanContext).TraceID().String())

	if ctx.isEmptyID() {
		ctx.withValue(keyRequestID, uuid.New().String())
	}
	return ctx
}
