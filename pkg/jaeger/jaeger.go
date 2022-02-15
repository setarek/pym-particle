package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	cnf "github.com/setarek/pym-particle-microservice/config"
	"io"
)

func InitJaeger(cnf *cnf.Config) (opentracing.Tracer, io.Closer, error) {
	return config.Configuration{
		ServiceName: cnf.GetString("service_name"),
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}.NewTracer(
		config.Metrics(metrics.NullFactory),
	)
}
