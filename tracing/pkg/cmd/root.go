package cmd

import (
	"context"
	"errors"
	"fmt"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"golang-training/tracing/services/root"
	"net"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/expvar"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootPort      = 8000
	formatterPort = 8084
	publisherPort = 8085

	logger         *zap.Logger
	metricsFactory metrics.Factory
)

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running cmd", zap.Error(err))
	}
	return err
}

func init() {
	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	metricsFactory = expvar.NewFactory(10) // 10 buckets for histograms
	logger.Info("Using expvar as metrics backend")
}

func Root(args []string) {
	if len(args) != 3 {
		logError(logger, errors.New("ERROR: Expecting 2 argument"))
	}
	helloTo := args[1]
	greeting := args[2]

	// init tracer
	serviceName := "hello-world"
	zapLogger := logger.With(zap.String("service", serviceName))
	slogger := log.NewFactory(zapLogger)
	tracer := tracing.Init(serviceName, metricsFactory, slogger)
	// need to set with StartSpanFromContext
	opentracing.SetGlobalTracer(tracer)

	// start root-span
	// operation name: say-hello
	ctx := context.Background()
	span := tracer.StartSpan("say-hello")
	// set tag
	span.SetTag("hello-to", helloTo)
	defer span.Finish()

	// baggage
	span.SetBaggageItem("greeting", greeting)
	fmt.Println(span)

	// attach root-span to context & pass ctx to child services
	// ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	// run
	// server := root.NewServer(
	// 	ctx,
	// 	tracer,
	// 	slogger,
	// 	net.JoinHostPort("http://localhost", strconv.Itoa(formatterPort)),
	// 	net.JoinHostPort("http://localhost", strconv.Itoa(publisherPort)),
	// )
	// server.Run(ctx, helloTo)
}

func RootWeb() {
	host := net.JoinHostPort("0.0.0.0", strconv.Itoa(rootPort))

	// create log factory
	zapLogger := logger.With(zap.String("service", "root"))
	logger := log.NewFactory(zapLogger)
	// server
	server := root.NewServer(
		host,
		tracing.Init("root", metricsFactory, logger),
		logger,
		net.JoinHostPort("localhost", strconv.Itoa(formatterPort)),
		net.JoinHostPort("localhost", strconv.Itoa(publisherPort)),
	)
	logError(zapLogger, server.Run())
}
