package cmd

import (
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/services/formatter"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(formatterCmd)
}

var formatterCmd = &cobra.Command{
	Use:   "formatter",
	Short: "Start Formmater Service",
	Long:  `Start Formmater Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ServiceFormatter()
	},
}

func ServiceFormatter() error {
	// create log factory
	zapLogger := logger.With(zap.String("service", configs.Formatter.ServiceName))
	logger := log.NewFactory(zapLogger)
	// server
	server := formatter.NewServer(
		configs.Formatter,
		metricsFactory,
		logger,
	)
	return logError(zapLogger, server.Run())
}
