package log

import "go.uber.org/zap"

type Factory struct {
	logger *zap.Logger
}

func NewFactory(logger *zap.Logger) Factory {
	return Factory{logger: logger}
}

func (f *Factory) Bg() Logger {
	return logger{logger: f.logger}
}
