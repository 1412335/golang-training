package formatter

import (
	"fmt"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Server struct {
	host   string
	tracer opentracing.Tracer
	logger log.Factory
}

func NewServer(host string, tracer opentracing.Tracer, logger log.Factory) *Server {
	return &Server{
		host:   host,
		tracer: tracer,
		logger: logger,
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
