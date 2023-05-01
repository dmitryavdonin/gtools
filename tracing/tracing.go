package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func New(ctx context.Context, serviceName, host string, port int) (io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		RPCMetrics:  true,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", host, port),
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil

}
