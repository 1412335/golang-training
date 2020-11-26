package root

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

type Server struct {
	ctx           context.Context
	host          string
	tracer        opentracing.Tracer
	httpClient    *tracing.HTTPClient
	logger        log.Factory
	formatterHost string
	publisherHost string
}

func NewServer(host string, tracer opentracing.Tracer, logger log.Factory, formatterHost, publisherHost string) *Server {
	return &Server{
		host:   host,
		tracer: tracer,
		httpClient: &tracing.HTTPClient{
			Client: &http.Client{Transport: &nethttp.Transport{}},
			Tracer: tracer,
		},
		logger:        logger,
		formatterHost: formatterHost,
		publisherHost: publisherHost,
	}
}

func (s *Server) Run() error {
	mux := s.createServerMux()
	return http.ListenAndServe(s.host, mux)
}

func (s *Server) createServerMux() http.Handler {
	mux := tracing.NewTracerServerMux(s.tracer)
	mux.Handle("/format", http.HandlerFunc(s.format))
	return mux
}

func (s *Server) format(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// span := rootSpan.Tracer().StartSpan(
	// 	"formatString",
	// 	opentracing.ChildOf(rootSpan.Context()),
	// )
	// span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	// defer span.Finish()

	// get child span
	// span := opentracing.SpanFromContext(req.Context())
	s.logger.For(ctx).Info("Internal request received", zap.String("method", "format"))

	// do request to formatter service
	v := url.Values{}
	v.Set("helloTo", "namnn")
	// url := s.formatterHost + "/format?" + v.Encode()
	url := "http://localhost:8084/format?" + v.Encode()
	resp, err := s.httpClient.Do(ctx, url)
	if err != nil {
		s.logger.For(ctx).Error("HTTP request failed", zap.Error(err))
		return
	}

	helloStr := string(resp)
	// client logs span
	s.logger.For(ctx).Info("Internal response", zap.String("helloStr", helloStr))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(helloStr))
}

func (s *Server) printHello(helloStr string) {
	// span := rootSpan.Tracer().StartSpan(
	// 	"printHello",
	// 	opentracing.ChildOf(rootSpan.Context()),
	// )

	// span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	// defer span.Finish()

	// get root span
	// span := opentracing.SpanFromContext(req.Context())
	// span := opentracing.SpanFromContext(ctx)
	s.logger.For(s.ctx).Info("Internal request received", zap.String("method", "print"))

	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://localhost:8085/publish?" + v.Encode()
	// url := s.publisherHost + "/publish?" + v.Encode()
	_, err := s.httpClient.Do(s.ctx, url)
	if err != nil {
		// set span tag error=true
		s.logger.For(s.ctx).Error("HTTP request failed", zap.Error(err))
	}
}
