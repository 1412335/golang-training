package cmd

import (
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/services/publisher"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(publisherCmd)
}

var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "Start Publisher Service",
	Long:  `Start Publisher Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ServicePublisher()
	},
}

func ServicePublisher() error {
	// create log factory
	zapLogger := logger.With(zap.String("service", configs.Publisher.ServiceName))
	logger := log.NewFactory(zapLogger)
	// server
	server := publisher.NewServer(
		configs.Publisher,
		metricsFactory,
		logger,
	)
	return logError(zapLogger, server.Run())
}
