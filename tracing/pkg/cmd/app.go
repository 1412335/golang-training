package cmd

import (
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/services/root"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start App Service",
	Long:  `Start App Service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return AppService()
	},
}

func AppService() error {
	// create log factory
	zapLogger := logger.With(zap.String("service", configs.ServiceName))
	logger := log.NewFactory(zapLogger)
	// server
	server := root.NewServer(
		configs,
		metricsFactory,
		logger,
	)
	return logError(zapLogger, server.Run())
}
