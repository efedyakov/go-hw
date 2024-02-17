package logger

import (
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func New(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.Level = lvl

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{zap: logger}, nil
}

func (l Logger) Debug(msg string) {
	l.zap.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.zap.Info(msg)
}

func (l Logger) Warn(msg string) {
	l.zap.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.zap.Error(msg)
}

func (l Logger) Fatal(msg string) {
	l.zap.Fatal(msg)
	os.Exit(1)
}
