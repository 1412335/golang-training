package formatter

import (
	"context"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"net/http"
	"net/url"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// interface export formatter client
type Client interface {
	Format(context.Context, string) (string, error)
}

type ClientImpl struct {
	host       string
	tracer     opentracing.Tracer
	httpClient *tracing.HTTPClient
	logger     log.Factory
}

func NewClient(host string, tracer opentracing.Tracer, logger log.Factory) Client {
	return &ClientImpl{
		host:   host,
		tracer: tracer,
		logger: logger,
		httpClient: &tracing.HTTPClient{
			Tracer: tracer,
			Client: &http.Client{Transport: &nethttp.Transport{}},
		},
	}
}

func (c *ClientImpl) Format(ctx context.Context, helloTo string) (string, error) {
	c.logger.For(ctx).Info("Request format", zap.String("helloTo", helloTo))
	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://" + c.host + "/format?" + v.Encode()
	resp, err := c.httpClient.Do(ctx, url)
	return string(resp), err
}
