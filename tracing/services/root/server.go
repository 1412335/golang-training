package root

import (
	"encoding/json"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"net/http"
	"strconv"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Server struct {
	host       string
	tracer     opentracing.Tracer
	httpClient *tracing.HTTPClient
	logger     log.Factory
	service    *service
}

type Config struct {
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
		service: newService(tracer, logger, Config{
			formatterHost: formatterHost,
			publisherHost: publisherHost,
		}),
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	helloTo := r.URL.Query().Get("helloTo")
	if helloTo == "" {
		http.Error(w, "missing helloTo", http.StatusBadRequest)
		return
	}

	greeting := r.URL.Query().Get("greeting")

	num := int32(1)
	numStr := r.URL.Query().Get("num")
	if numStr != "" {
		if i, err := strconv.Atoi(numStr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.For(ctx).Error("parse num to int failed", zap.Error(err))
			return
		} else {
			num = int32(i)
		}
	}

	// do request to formatter service
	resp, err := s.service.Get(ctx, helloTo, greeting, num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.For(ctx).Error("HTTP request failed", zap.Error(err))
		return
	}

	// client logs span
	s.logger.For(ctx).Info("response", zap.Strings("helloStr", resp))

	s.writeResponse(resp, w, r)
}

func (s *Server) writeResponse(response interface{}, w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.For(r.Context()).Error("cannot marshal response", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
