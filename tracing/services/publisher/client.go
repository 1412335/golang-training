package publisher

import (
	"context"
	"golang-training/tracing/pkg/log"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

// interface export formatter client
type Client interface {
	Echo(context.Context, string, int32) ([]string, error)
}

type ClientImpl struct {
	host   string
	tracer opentracing.Tracer
	logger log.Factory
	client PublisherServiceClient
}

func NewClient(host string, tracer opentracing.Tracer, logger log.Factory) Client {
	conn, err := grpc.Dial(host,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		logger.Bg().Fatal("cannot create grpc connection", zap.Error(err))
		return nil
	}
	return &ClientImpl{
		host:   host,
		tracer: tracer,
		logger: logger,
		client: NewPublisherServiceClient(conn),
	}
}

func (c *ClientImpl) Echo(ctx context.Context, helloStr string, num int32) ([]string, error) {
	c.logger.For(ctx).Info("Echo request", zap.String("helloStr", helloStr), zap.Int32("num", num))
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	echoResp, err := c.client.Echo(ctx, &EchoRequest{HelloStr: helloStr, Num: num})
	if err != nil {
		return nil, err
	}
	return c.formatProtoResp(echoResp), nil
}

func (c *ClientImpl) formatProtoResp(echoResp *EchoResponse) []string {
	helloArr := make([]string, len(echoResp.HelloStr))
	for i, helloStr := range echoResp.HelloStr {
		helloArr[i] = helloStr
	}
	return helloArr
}
