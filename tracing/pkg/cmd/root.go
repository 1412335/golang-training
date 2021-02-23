package cmd

import (
	"golang-training/tracing/pkg/config"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/expvar"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Used for flags.
	cfgFile        string
	metricsBackend string

	configs        *config.ServiceConfig
	logger         *zap.Logger
	metricsFactory metrics.Factory

	formatterHost string
	publisherHost string

	rootCmd = &cobra.Command{
		Use:   "tracing",
		Short: "Tracing Example With Jaeger",
		Long:  `Tracing Example With Jaeger`,
	}
)

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running cmd", zap.Error(err))
	}
	return err
}

func initConfig() {
	configs = &config.ServiceConfig{}
	if err := config.LoadConfig(cfgFile, configs); err != nil {
		logger.Fatal("Load config failed", zap.Error(err))
	}
	logger.Info("Load config success", zap.String("config_file", viper.ConfigFileUsed()))

	if configs.Metrics == "expvar" {
		metricsFactory = expvar.NewFactory(10) // 10 buckets for histograms
		logger.Info("Using expvar as metrics backend")
	} else {
		metricsFactory = jprom.New().Namespace(metrics.NSOptions{Name: "tracing", Tags: nil})
		logger.Info("Using prometheus as metrics backend")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yml)")
	rootCmd.PersistentFlags().StringVarP(&metricsBackend, "metrics", "m", "prometheus", "metrics backend expvar|prometheus (default: prometheus)")

	viper.BindPFlag("metrics", rootCmd.PersistentFlags().Lookup("metrics"))

	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Execute cmd failed", zap.Error(err))
		os.Exit(-1)
	}
}
