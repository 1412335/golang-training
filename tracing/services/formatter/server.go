package formatter

import (
	"fmt"
	"golang-training/tracing/pkg/config"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"net"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Server struct {
	host           string
	tracer         opentracing.Tracer
	metricsFactory metrics.Factory
	logger         log.Factory
}

func NewServer(configs *config.Formatter, metricsFactory metrics.Factory, logger log.Factory) *Server {
	host := net.JoinHostPort("0.0.0.0", strconv.Itoa(configs.Port))
	// create tracer
	tracer := tracing.Init(configs.ServiceName, metricsFactory, logger)
	return &Server{
		host:           host,
		tracer:         tracer,
		metricsFactory: metricsFactory,
		logger:         logger,
	}
}

func (s *Server) Run() error {
	mux := s.createServerMux()
	s.logger.Bg().Info("Starting server", zap.String("host", s.host))
	return http.ListenAndServe(s.host, mux)
}

func (s *Server) createServerMux() http.Handler {
	mux := tracing.NewTracerServerMux(s.tracer)
	mux.Handle("/format", http.HandlerFunc(s.format))
	return mux
}

func (s *Server) format(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	// get span from request context
	span := opentracing.SpanFromContext(ctx)
	// defer span.Finish() //=> duplicate span

	// get baggage
	greeting := span.BaggageItem("greeting")
	if greeting == "" {
		greeting = "hello"
	}

	helloTo := r.FormValue("helloTo")
	helloStr := fmt.Sprintf("%s, %s!", greeting, helloTo)
	// logs
	s.logger.For(ctx).Info("HTTP response", zap.String("format", helloStr))
	// response
	w.Write([]byte(helloStr))
}
