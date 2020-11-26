package root

import (
	"context"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/services/formatter"
	"golang-training/tracing/services/publisher"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type service struct {
	tracer    opentracing.Tracer
	logger    log.Factory
	formatter formatter.Client
	publisher publisher.Client
}

func newService(tracer opentracing.Tracer, logger log.Factory, config Config) *service {
	return &service{
		tracer:    tracer,
		logger:    logger,
		formatter: formatter.NewClient(config.formatterHost, tracer, logger),
		publisher: publisher.NewClient(config.publisherHost, tracer, logger),
	}
}

func (s *service) Get(ctx context.Context, helloTo, greeting string, num int32) ([]string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetBaggageItem("greeting", greeting)
	}
	helloStr, err := s.formatter.Format(ctx, helloTo)
	if err != nil {
		return nil, err
	}
	helloArr, err := s.publisher.Echo(ctx, helloStr, num)
	if err != nil {
		return nil, err
	}
	s.logger.For(ctx).Info("Get format & echo success", zap.String("helloStr", helloStr), zap.String("greeting", greeting), zap.Int32("num", num))
	return helloArr, nil
}
