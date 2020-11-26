package cmd

import (
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"golang-training/tracing/services/formatter"
	"net"
	"strconv"

	"go.uber.org/zap"
)

func ServiceFormatter() {
	host := net.JoinHostPort("0.0.0.0", strconv.Itoa(formatterPort))

	// create log factory
	zapLogger := logger.With(zap.String("service", "formatter"))
	logger := log.NewFactory(zapLogger)
	// server
	server := formatter.NewServer(
		host,
		tracing.Init("formatter", metricsFactory, logger),
		logger,
	)
	logError(zapLogger, server.Run())
}
