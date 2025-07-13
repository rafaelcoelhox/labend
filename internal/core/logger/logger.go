package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

type zapLogger struct {
	zap *zap.Logger
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func New() (Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	zapLog, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{zap: zapLog}, nil
}

func NewDevelopment() (Logger, error) {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &zapLogger{zap: zapLog}, nil
}
