package cmd

import (
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"golang-training/tracing/services/publisher"
	"net"
	"strconv"

	"go.uber.org/zap"
)

func ServicePublisher() {
	host := net.JoinHostPort("0.0.0.0", strconv.Itoa(publisherPort))

	// create log factory
	serviceName := "publisher"
	zapLogger := logger.With(zap.String("service", serviceName))
	logger := log.NewFactory(zapLogger)
	// server
	server := publisher.NewServer(
		host,
		tracing.Init(serviceName, metricsFactory, logger),
		logger,
	)
	logError(zapLogger, server.Run())
}
